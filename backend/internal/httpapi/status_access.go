package httpapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"remember_my_pill/backend/internal/access"
)

const statusAccessUnavailableCode = "STATUS_ACCESS_UNAVAILABLE"

type statusAccessLimiter struct {
	mu   sync.Mutex
	hits map[string][]time.Time
	key  []byte
	now  func() time.Time
}

func newStatusAccessLimiter(key []byte, now func() time.Time) *statusAccessLimiter {
	if now == nil {
		now = time.Now
	}
	return &statusAccessLimiter{hits: make(map[string][]time.Time), key: key, now: now}
}

func (l *statusAccessLimiter) allow(kind, value string, limit int, window time.Duration) bool {
	mac := hmac.New(sha256.New, l.key)
	_, _ = mac.Write([]byte(kind + ":" + value))
	key := hex.EncodeToString(mac.Sum(nil))
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	kept := make([]time.Time, 0, len(l.hits[key]))
	for _, hit := range l.hits[key] {
		if now.Sub(hit) < window {
			kept = append(kept, hit)
		}
	}
	if len(kept) >= limit {
		l.hits[key] = kept
		return false
	}
	l.hits[key] = append(kept, now)
	return true
}

type statusAccessRequest struct {
	Email string `json:"email"`
}

type statusAccessExchange struct {
	VerificationToken string `json:"verificationToken"`
}

func registerStatusAccessRoutes(r chi.Router, options Options, _ any) {
	limiter := newStatusAccessLimiter(options.StatusAccessRateLimitKey, options.Now)
	active := func() bool {
		return options.StatusAccessEnabled && options.PilotLegalContentApproved && options.StatusStore != nil && options.StatusEmailSender != nil
	}
	unavailable := func(w http.ResponseWriter) {
		writeError(w, http.StatusNotFound, statusAccessUnavailableCode, "Status access is not available.")
	}
	r.Post("/api/waitlist/status-access/request", func(w http.ResponseWriter, req *http.Request) {
		if !active() {
			unavailable(w)
			return
		}
		decoder := json.NewDecoder(http.MaxBytesReader(w, req.Body, 1024))
		decoder.DisallowUnknownFields()
		var in statusAccessRequest
		if decoder.Decode(&in) != nil || decoder.Decode(&struct{}{}) != io.EOF {
			writeAccepted(w)
			return
		}
		email := strings.ToLower(strings.TrimSpace(in.Email))
		if !emailPattern.MatchString(email) || !limiter.allow("ip", clientIP(req, options.TrustedProxies), 3, time.Hour) || !limiter.allow("email", email, 3, 24*time.Hour) {
			writeAccepted(w)
			return
		}
		eligibility, entryID, err := access.EvaluateStatusAccessEligibility(req.Context(), options.StatusStore, email)
		if err != nil || eligibility != access.StatusAccessEligible {
			writeAccepted(w)
			return
		}
		raw, err := access.GenerateToken()
		if err != nil {
			writeAccepted(w)
			return
		}
		if _, err := access.CreateVerificationToken(req.Context(), options.StatusStore, entryID, raw, now(options.Now)); err == nil {
			if verificationURL, err := access.BuildVerificationURL(options.StatusAccessBaseURL, raw); err == nil {
				// Delivery outcomes are intentionally never exposed. This boundary
				// is activated only in isolated tests until legal/template approval.
				_ = options.StatusEmailSender.SendStatusAccessEmail(req.Context(), access.StatusAccessEmail{To: email, VerificationURL: verificationURL})
			}
		}
		writeAccepted(w)
	})
	r.Post("/api/waitlist/status-access/exchange", func(w http.ResponseWriter, req *http.Request) {
		if !active() {
			unavailable(w)
			return
		}
		if req.URL.RawQuery != "" || !isJSON(req) {
			writeStatusUnauthorized(w)
			return
		}
		decoder := json.NewDecoder(http.MaxBytesReader(w, req.Body, 1024))
		decoder.DisallowUnknownFields()
		var in statusAccessExchange
		if decoder.Decode(&in) != nil || decoder.Decode(&struct{}{}) != io.EOF || strings.TrimSpace(in.VerificationToken) == "" {
			writeStatusUnauthorized(w)
			return
		}
		session, err := access.ConsumeVerificationAndCreateBrowserSession(req.Context(), options.StatusStore, strings.TrimSpace(in.VerificationToken), now(options.Now))
		if err != nil {
			writeStatusUnauthorized(w)
			return
		}
		http.SetCookie(w, access.StatusSessionCookie(session.RawID, options.StatusSessionCookieSecure, now(options.Now)))
		writeJSON(w, http.StatusOK, `{"status":"authenticated"}`)
	})
	r.Get("/api/waitlist/status", func(w http.ResponseWriter, req *http.Request) {
		if !active() {
			unavailable(w)
			return
		}
		cookie, err := req.Cookie(statusCookieName(options.StatusSessionCookieSecure))
		if err != nil || cookie.Value == "" {
			writeStatusUnauthorized(w)
			return
		}
		entryID, err := access.LookupBrowserSession(req.Context(), options.StatusStore, cookie.Value, now(options.Now))
		if err != nil {
			writeStatusUnauthorized(w)
			return
		}
		status, err := access.LookupWaitlistStatus(req.Context(), options.StatusStore, entryID)
		if err != nil {
			writeStatusUnauthorized(w)
			return
		}
		response := struct {
			Rank          int64  `json:"rank"`
			ReferralCount int64  `json:"referralCount"`
			ReferralCode  string `json:"referralCode,omitempty"`
			ReferralURL   string `json:"referralUrl,omitempty"`
		}{Rank: status.Rank, ReferralCount: status.ReferralCount, ReferralCode: status.ReferralCode}
		if status.ReferralCode != "" {
			response.ReferralURL = referralURL(options.StatusAccessBaseURL, status.ReferralCode)
		}
		writeJSONValue(w, http.StatusOK, response)
	})
	r.Post("/api/waitlist/status/logout", func(w http.ResponseWriter, req *http.Request) {
		if !active() {
			unavailable(w)
			return
		}
		if !isJSON(req) || !validLogoutCSRF(req, options.AllowedOrigin) {
			writeError(w, http.StatusForbidden, "STATUS_ACCESS_CSRF_REJECTED", "Request could not be completed.")
			return
		}
		if cookie, err := req.Cookie(statusCookieName(options.StatusSessionCookieSecure)); err == nil && cookie.Value != "" {
			_ = access.RevokeBrowserSession(req.Context(), options.StatusStore, cookie.Value, now(options.Now))
		}
		http.SetCookie(w, access.ClearStatusSessionCookie(options.StatusSessionCookieSecure))
		w.WriteHeader(http.StatusNoContent)
	})
}

func now(clock func() time.Time) time.Time {
	if clock == nil {
		return time.Now().UTC()
	}
	return clock().UTC()
}

func isJSON(req *http.Request) bool {
	contentType, _, _ := strings.Cut(req.Header.Get("Content-Type"), ";")
	return strings.EqualFold(strings.TrimSpace(contentType), "application/json")
}

func validLogoutCSRF(req *http.Request, allowedOrigin string) bool {
	if req.Header.Get("Origin") != allowedOrigin {
		return false
	}
	site := req.Header.Get("Sec-Fetch-Site")
	return site == "same-origin" || site == "same-site"
}

func statusCookieName(secure bool) string {
	if secure {
		return access.ProductionStatusCookieName
	}
	return access.DevelopmentStatusCookieName
}

func referralURL(base, code string) string {
	u, err := url.Parse(base)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return ""
	}
	u.Path = strings.TrimRight(u.Path, "/") + "/waitlist"
	u.RawQuery = url.Values{"ref": []string{code}}.Encode()
	u.Fragment = ""
	return u.String()
}

func writeStatusUnauthorized(w http.ResponseWriter) {
	writeError(w, http.StatusUnauthorized, "STATUS_ACCESS_UNAUTHORIZED", "Status access could not be verified.")
}

func writeJSONValue(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
