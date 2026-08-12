package httpapi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/netip"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ValidationCode = "WAITLIST_VALIDATION_ERROR"
	InternalCode   = "WAITLIST_INTERNAL_ERROR"
)

var (
	emailPattern     = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	requestIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{8,64}$`)
)

type Store interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

type Readiness interface {
	Ping(context.Context) error
}

type Options struct {
	AllowedOrigin             string
	ConsentVersion            string
	PilotLegalContentApproved bool
	TrustedProxies            []netip.Prefix
	Logger                    *slog.Logger
	Now                       func() time.Time
}

type Limiter struct {
	mu   sync.Mutex
	hits map[string][]time.Time
	now  func() time.Time
}

func NewLimiter(now func() time.Time) *Limiter {
	if now == nil {
		now = time.Now
	}
	return &Limiter{hits: map[string][]time.Time{}, now: now}
}

func (l *Limiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	kept := make([]time.Time, 0, len(l.hits[ip]))
	for _, hit := range l.hits[ip] {
		if now.Sub(hit) < time.Minute {
			kept = append(kept, hit)
		}
	}
	if len(kept) >= 5 {
		l.hits[ip] = kept
		return false
	}
	l.hits[ip] = append(kept, now)
	return true
}

type errorBody struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type signupRequest struct {
	Name                    string  `json:"name"`
	Email                   string  `json:"email"`
	Consent                 *bool   `json:"consent"`
	ConsentVersion          *string `json:"consentVersion"`
	ReferralCode            *string `json:"referralCode"`
	MarketingConsent        *bool   `json:"marketingConsent"`
	MarketingConsentVersion *string `json:"marketingConsentVersion"`
	Company                 *string `json:"company"`
}

type contextKey string

const requestIDKey contextKey = "request_id"

func NewRouter(store Store, readiness Readiness, options Options) http.Handler {
	logger := options.Logger
	if logger == nil {
		logger = slog.Default()
	}
	limiter := NewLimiter(options.Now)
	r := chi.NewRouter()
	r.Use(requestLogging(logger, options.TrustedProxies))
	r.Use(cors(options.AllowedOrigin))
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, `{"status":"ok"}`)
	})
	r.Get("/ready", func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
		defer cancel()
		if err := readiness.Ping(ctx); err != nil {
			logger.Warn("readiness failed", "request_id", requestID(req.Context()))
			writeJSON(w, http.StatusServiceUnavailable, `{"status":"unavailable"}`)
			return
		}
		writeJSON(w, http.StatusOK, `{"status":"ready"}`)
	})
	r.Post("/api/waitlist", func(w http.ResponseWriter, req *http.Request) {
		if !limiter.allow(clientIP(req, options.TrustedProxies)) {
			logger.Warn("rate limit rejected", "request_id", requestID(req.Context()))
			writeError(w, http.StatusTooManyRequests, "WAITLIST_RATE_LIMITED", "Please wait a moment before trying again.")
			return
		}

		decoder := json.NewDecoder(http.MaxBytesReader(w, req.Body, 4096))
		decoder.DisallowUnknownFields()
		var in signupRequest
		if decoder.Decode(&in) != nil || decoder.Decode(&struct{}{}) != io.EOF {
			writeError(w, http.StatusBadRequest, ValidationCode, "Enter a valid email.")
			return
		}
		if in.Company != nil && strings.TrimSpace(*in.Company) != "" {
			writeAccepted(w)
			return
		}
		in.Name = strings.TrimSpace(in.Name)
		in.Email = strings.ToLower(strings.TrimSpace(in.Email))
		if in.ReferralCode != nil {
			code := strings.TrimSpace(*in.ReferralCode)
			if len(code) > 50 {
				writeError(w, http.StatusBadRequest, ValidationCode, "Enter a valid email.")
				return
			}
		}
		if len(in.Name) > 100 || len(in.Email) == 0 || len(in.Email) > 255 || !emailPattern.MatchString(in.Email) {
			writeError(w, http.StatusBadRequest, ValidationCode, "Enter a valid email.")
			return
		}
		consentVersion, validConsent := validateConsent(in, options.ConsentVersion, options.PilotLegalContentApproved)
		if !validConsent {
			writeError(w, http.StatusBadRequest, ValidationCode, "Enter a valid email.")
			return
		}
		marketingVersion, validMarketing := validateMarketingConsent(in)
		if !validMarketing || (marketingVersion != nil && consentVersion == nil) {
			writeError(w, http.StatusBadRequest, ValidationCode, "Enter a valid email.")
			return
		}

		var name any
		if in.Name != "" {
			name = in.Name
		}
		_, err := store.Exec(req.Context(), `INSERT INTO waitlist_entries (name, email_normalized, consent_version, consented_at, marketing_consent_version, marketing_consented_at) VALUES ($1, $2, $3::varchar, CASE WHEN $3::varchar IS NULL THEN NULL ELSE CURRENT_TIMESTAMP END, $4::varchar, CASE WHEN $4::varchar IS NULL THEN NULL ELSE CURRENT_TIMESTAMP END)`, name, in.Email, consentVersion, marketingVersion)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				writeAccepted(w)
				return
			}
			logger.Error("waitlist insert failed", "request_id", requestID(req.Context()))
			writeError(w, http.StatusInternalServerError, InternalCode, "We could not join the waitlist. Please try again.")
			return
		}
		writeAccepted(w)
	})
	return r
}

func validateConsent(in signupRequest, activeVersion string, required bool) (any, bool) {
	if in.Consent == nil && in.ConsentVersion == nil {
		return nil, !required
	}
	if in.Consent == nil || !*in.Consent || in.ConsentVersion == nil {
		return nil, false
	}
	version := strings.TrimSpace(*in.ConsentVersion)
	if activeVersion == "" || version == "" || version != activeVersion {
		return nil, false
	}
	return version, true
}

func validateMarketingConsent(in signupRequest) (any, bool) {
	if in.MarketingConsent == nil {
		return nil, in.MarketingConsentVersion == nil
	}
	if *in.MarketingConsent {
		if in.MarketingConsentVersion == nil || strings.TrimSpace(*in.MarketingConsentVersion) != "marketing-consent-v1" {
			return nil, false
		}
		return "marketing-consent-v1", true
	}
	if in.MarketingConsentVersion != nil {
		return nil, false
	}
	return nil, true
}

func writeAccepted(w http.ResponseWriter) {
	writeJSON(w, http.StatusAccepted, `{"status":"accepted"}`)
}

func writeJSON(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = io.WriteString(w, body)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	var body errorBody
	body.Error.Code = code
	body.Error.Message = message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func cors(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			origin := req.Header.Get("Origin")
			if origin != "" && origin == allowedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
			}
			if req.Method == http.MethodOptions {
				if origin == "" || origin != allowedOrigin {
					w.WriteHeader(http.StatusForbidden)
					return
				}
				w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}

func requestLogging(logger *slog.Logger, trustedProxies []netip.Prefix) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			id := incomingRequestID(req, trustedProxies)
			if id == "" {
				id = generatedRequestID()
			}
			w.Header().Set("X-Request-ID", id)
			start := time.Now()
			recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(recorder, req.WithContext(context.WithValue(req.Context(), requestIDKey, id)))
			logger.Info("request completed", "request_id", id, "method", req.Method, "path", req.URL.Path, "status", recorder.status, "duration_ms", time.Since(start).Milliseconds())
		})
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (w *statusRecorder) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func requestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}

func incomingRequestID(req *http.Request, trustedProxies []netip.Prefix) string {
	if !peerIsTrusted(req.RemoteAddr, trustedProxies) {
		return ""
	}
	id := req.Header.Get("X-Request-ID")
	if !requestIDPattern.MatchString(id) {
		return ""
	}
	return id
}

func generatedRequestID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "generated-request-id"
	}
	return hex.EncodeToString(bytes)
}

func clientIP(req *http.Request, trustedProxies []netip.Prefix) string {
	if peerIsTrusted(req.RemoteAddr, trustedProxies) {
		if value := forwardedFor(req.Header.Get("Forwarded")); value != "" {
			return value
		}
		if value := strings.TrimSpace(strings.Split(req.Header.Get("X-Forwarded-For"), ",")[0]); net.ParseIP(value) != nil {
			return value
		}
	}
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err == nil {
		return host
	}
	return req.RemoteAddr
}

// forwardedFor accepts the first RFC 7239 for= value only after the direct
// connection has been identified as a configured trusted proxy.
func forwardedFor(header string) string {
	first := strings.TrimSpace(strings.Split(header, ",")[0])
	for _, part := range strings.Split(first, ";") {
		key, value, found := strings.Cut(strings.TrimSpace(part), "=")
		if !found || !strings.EqualFold(key, "for") {
			continue
		}
		value = strings.Trim(strings.TrimSpace(value), "\"")
		value = strings.Trim(value, "[]")
		if net.ParseIP(value) != nil {
			return value
		}
	}
	return ""
}

func peerIsTrusted(remoteAddr string, trustedProxies []netip.Prefix) bool {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		host = remoteAddr
	}
	addr, err := netip.ParseAddr(host)
	if err != nil {
		return false
	}
	for _, prefix := range trustedProxies {
		if prefix.Contains(addr) {
			return true
		}
	}
	return false
}
