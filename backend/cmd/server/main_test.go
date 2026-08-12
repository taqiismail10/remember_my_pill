package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"remember_my_pill/backend/internal/httpapi"
)

type fakeStore struct {
	mu               sync.Mutex
	calls            int
	name             any
	email            string
	consent          any
	marketingConsent any
	query            string
	err              error
}

func (s *fakeStore) Exec(_ context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls++
	s.name, s.email, s.consent, s.marketingConsent, s.query = args[0], args[1].(string), args[2], args[3], query
	return pgconn.CommandTag{}, s.err
}

type fakeReadiness struct{ err error }

func (r fakeReadiness) Ping(context.Context) error { return r.err }

func routerFor(store *fakeStore, readiness error, options ...httpapi.Options) http.Handler {
	option := httpapi.Options{
		AllowedOrigin:  "http://localhost:3000",
		ConsentVersion: "waitlist-consent-v1",
		Logger:         slog.New(slog.NewTextHandler(ioDiscard{}, nil)),
		Now:            func() time.Time { return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) },
	}
	if len(options) > 0 {
		option = options[0]
		if option.Logger == nil {
			option.Logger = slog.New(slog.NewTextHandler(ioDiscard{}, nil))
		}
	}
	return httpapi.NewRouter(store, fakeReadiness{err: readiness}, option)
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) { return len(p), nil }

func post(r http.Handler, body string, remoteAddr string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/api/waitlist", strings.NewReader(body))
	req.RemoteAddr = remoteAddr
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func requireJSON(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()
	if !strings.HasPrefix(w.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("content type = %q", w.Header().Get("Content-Type"))
	}
}

func errorCode(t *testing.T, w *httptest.ResponseRecorder) string {
	t.Helper()
	var body struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	return body.Error.Code
}

func requireAccepted(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()
	if w.Code != http.StatusAccepted || w.Body.String() != `{"status":"accepted"}` {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	requireJSON(t, w)
}

func TestWaitlistAcceptsEmailOnlyAndOptionalName(t *testing.T) {
	for _, tc := range []struct {
		name, body, wantName, wantEmail string
		wantNil                         bool
	}{
		{"email only", `{"email":" Ada@Example.Test "}`, "", "ada@example.test", true},
		{"name and email", `{"name":" Ada Lovelace ","email":" ADA@EXAMPLE.TEST "}`, "Ada Lovelace", "ada@example.test", false},
		{"blank name", `{"name":"  ","email":"grace@example.test"}`, "", "grace@example.test", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			store := &fakeStore{}
			requireAccepted(t, post(routerFor(store, nil), tc.body, "203.0.113.10:1234"))
			if store.calls != 1 || store.email != tc.wantEmail || (tc.wantNil && store.name != nil) || (!tc.wantNil && store.name != tc.wantName) || store.consent != nil {
				t.Fatalf("store=%#v", store)
			}
		})
	}
}

func TestWaitlistRejectsInvalidPayloads(t *testing.T) {
	cases := []string{`{}`, `{"name":"A"}`, `{"email":" "}`, `{"email":"invalid"}`, `{"email":"a@example.test","extra":true}`, `{`, `{"email":"a@example.test"} {}`, `[]`, `{"email":"a@example.test","consentVersion":"waitlist-consent-v1"}`, `{"email":"a@example.test","consent":false,"consentVersion":"waitlist-consent-v1"}`, `{"email":"a@example.test","consent":true,"consentVersion":"stale"}`, `{"email":"a@example.test","marketingConsent":true}`, `{"email":"a@example.test","marketingConsent":true,"marketingConsentVersion":"wrong"}`, `{"email":"a@example.test","marketingConsent":true,"marketingConsentVersion":"marketing-consent-v1"}`, `{"email":"a@example.test","marketingConsent":false,"marketingConsentVersion":"marketing-consent-v1"}`}
	cases = append(cases, `{"name":"`+strings.Repeat("a", 101)+`","email":"a@example.test"}`, `{"email":"`+strings.Repeat("a", 246)+`@example.test"}`, `{"email":"a@example.test","name":"`+strings.Repeat("a", 4096)+`"}`)
	for _, body := range cases {
		t.Run(body[:min(16, len(body))], func(t *testing.T) {
			store := &fakeStore{}
			w := post(routerFor(store, nil), body, "203.0.113.10:1234")
			if w.Code != http.StatusBadRequest || errorCode(t, w) != httpapi.ValidationCode || store.calls != 0 {
				t.Fatalf("status=%d body=%s calls=%d", w.Code, w.Body.String(), store.calls)
			}
		})
	}
}

func TestConsentCompatibilityAndServerTimestampQuery(t *testing.T) {
	store := &fakeStore{}
	w := post(routerFor(store, nil), `{"email":"ada@example.test","consent":true,"consentVersion":"waitlist-consent-v1"}`, "203.0.113.10:1234")
	requireAccepted(t, w)
	if store.consent != "waitlist-consent-v1" {
		t.Fatalf("consent=%#v", store.consent)
	}
	if !strings.Contains(store.query, "CURRENT_TIMESTAMP") {
		t.Fatalf("consent timestamp was not assigned by SQL: %s", store.query)
	}
}

func TestMarketingConsentIsIndependent(t *testing.T) {
	store := &fakeStore{}
	requireAccepted(t, post(routerFor(store, nil), `{"email":"ada@example.test","consent":true,"consentVersion":"waitlist-consent-v1","marketingConsent":false}`, "203.0.113.10:1234"))
	if store.marketingConsent != nil {
		t.Fatalf("marketing consent fabricated: %#v", store.marketingConsent)
	}
	store = &fakeStore{}
	requireAccepted(t, post(routerFor(store, nil), `{"email":"ada@example.test","consent":true,"consentVersion":"waitlist-consent-v1","marketingConsent":true,"marketingConsentVersion":"marketing-consent-v1"}`, "203.0.113.10:1234"))
	if store.marketingConsent != "marketing-consent-v1" || !strings.Contains(store.query, "marketing_consented_at") || !strings.Contains(store.query, "CURRENT_TIMESTAMP") {
		t.Fatalf("marketing=%#v query=%s", store.marketingConsent, store.query)
	}
}

func TestLegalActivationRequiresWaitlistConsent(t *testing.T) {
	store := &fakeStore{}
	options := httpapi.Options{AllowedOrigin: "http://localhost:3000", ConsentVersion: "waitlist-consent-v1", PilotLegalContentApproved: true, Logger: slog.New(slog.NewTextHandler(ioDiscard{}, nil))}
	w := post(routerFor(store, nil, options), `{"email":"legacy@example.test"}`, "203.0.113.10:1234")
	if w.Code != http.StatusBadRequest || store.calls != 0 {
		t.Fatalf("status=%d calls=%d", w.Code, store.calls)
	}
	requireAccepted(t, post(routerFor(store, nil, options), `{"email":"new@example.test","consent":true,"consentVersion":"waitlist-consent-v1","marketingConsent":false}`, "203.0.113.11:1234"))
}

func TestDuplicateAndInternalErrors(t *testing.T) {
	duplicate := post(routerFor(&fakeStore{err: &pgconn.PgError{Code: "23505"}}, nil), `{"email":"ada@example.test"}`, "203.0.113.10:1234")
	requireAccepted(t, duplicate)
	internal := post(routerFor(&fakeStore{err: errors.New("postgres://secret@db/internal")}, nil), `{"email":"ada@example.test"}`, "203.0.113.10:1234")
	if internal.Code != http.StatusInternalServerError || errorCode(t, internal) != httpapi.InternalCode || strings.Contains(strings.ToLower(internal.Body.String()), "secret") {
		t.Fatalf("internal=%d %s", internal.Code, internal.Body.String())
	}
}

func TestHoneypotReturnsAcceptedWithoutInsert(t *testing.T) {
	store := &fakeStore{}
	requireAccepted(t, post(routerFor(store, nil), `{"email":"bot@example.test","company":"Bot Inc"}`, "203.0.113.10:1234"))
	if store.calls != 0 {
		t.Fatalf("honeypot inserted %d rows", store.calls)
	}
}

func TestRateLimitHealthAndReadiness(t *testing.T) {
	store := &fakeStore{}
	r := routerFor(store, nil)
	for i := 0; i < 5; i++ {
		requireAccepted(t, post(r, `{"email":"person`+string(rune('a'+i))+`@example.test"}`, "203.0.113.10:1234"))
	}
	limited := post(r, `{"email":"later@example.test"}`, "203.0.113.10:1234")
	if limited.Code != http.StatusTooManyRequests || errorCode(t, limited) != "WAITLIST_RATE_LIMITED" {
		t.Fatalf("limited=%d %s", limited.Code, limited.Body.String())
	}
	for _, endpoint := range []string{"/health", "/ready"} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, endpoint, nil))
		if w.Code != http.StatusOK {
			t.Fatalf("%s=%d %s", endpoint, w.Code, w.Body.String())
		}
	}
	unavailable := routerFor(&fakeStore{}, errors.New("down"))
	w := httptest.NewRecorder()
	unavailable.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ready", nil))
	if w.Code != http.StatusServiceUnavailable || strings.Contains(strings.ToLower(w.Body.String()), "down") {
		t.Fatalf("ready=%d %s", w.Code, w.Body.String())
	}
}

func TestCORSAndRequestIDProxyPolicy(t *testing.T) {
	r := routerFor(&fakeStore{}, nil)
	preflight := httptest.NewRequest(http.MethodOptions, "/api/waitlist", nil)
	preflight.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, preflight)
	if w.Code != http.StatusNoContent || w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Fatalf("cors=%d %#v", w.Code, w.Header())
	}

	untrusted := httptest.NewRequest(http.MethodPost, "/api/waitlist", strings.NewReader(`{"email":"ada@example.test"}`))
	untrusted.RemoteAddr = "203.0.113.10:1234"
	untrusted.Header.Set("X-Forwarded-For", "198.51.100.1")
	untrusted.Header.Set("X-Request-ID", "attacker-id-123")
	uw := httptest.NewRecorder()
	r.ServeHTTP(uw, untrusted)
	if uw.Header().Get("X-Request-ID") == "attacker-id-123" || uw.Header().Get("X-Request-ID") == "" {
		t.Fatalf("untrusted request id=%q", uw.Header().Get("X-Request-ID"))
	}
	for i := 0; i < 4; i++ {
		request := httptest.NewRequest(http.MethodPost, "/api/waitlist", strings.NewReader(`{"email":"untrusted`+string(rune('a'+i))+`@example.test"}`))
		request.RemoteAddr = "203.0.113.10:1234"
		request.Header.Set("X-Forwarded-For", "198.51.100.1")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, request)
		if w.Code != http.StatusAccepted {
			t.Fatalf("untrusted forwarded header bypassed limiter: %d", w.Code)
		}
	}
	blocked := httptest.NewRequest(http.MethodPost, "/api/waitlist", strings.NewReader(`{"email":"blocked@example.test"}`))
	blocked.RemoteAddr = "203.0.113.10:1234"
	blocked.Header.Set("X-Forwarded-For", "203.0.113.11")
	blockedWriter := httptest.NewRecorder()
	r.ServeHTTP(blockedWriter, blocked)
	if blockedWriter.Code != http.StatusTooManyRequests {
		t.Fatalf("untrusted forwarded header changed client IP: %d", blockedWriter.Code)
	}

	trusted := routerFor(&fakeStore{}, nil, httpapi.Options{AllowedOrigin: "http://localhost:3000", ConsentVersion: "waitlist-consent-v1", TrustedProxies: []netip.Prefix{netip.MustParsePrefix("127.0.0.1/32")}, Logger: slog.New(slog.NewTextHandler(ioDiscard{}, nil))})
	request := httptest.NewRequest(http.MethodPost, "/api/waitlist", strings.NewReader(`{"email":"trusted@example.test"}`))
	request.RemoteAddr = "127.0.0.1:1234"
	request.Header.Set("X-Forwarded-For", "198.51.100.1")
	request.Header.Set("X-Request-ID", "trusted-id-123")
	tw := httptest.NewRecorder()
	trusted.ServeHTTP(tw, request)
	if tw.Header().Get("X-Request-ID") != "trusted-id-123" {
		t.Fatalf("trusted request id=%q", tw.Header().Get("X-Request-ID"))
	}
	forwarded := httptest.NewRequest(http.MethodPost, "/api/waitlist", strings.NewReader(`{"email":"forwarded@example.test"}`))
	forwarded.RemoteAddr = "127.0.0.1:1234"
	forwarded.Header.Set("Forwarded", "for=198.51.100.2")
	fw := httptest.NewRecorder()
	trusted.ServeHTTP(fw, forwarded)
	if fw.Code != http.StatusAccepted {
		t.Fatalf("trusted Forwarded header was not accepted: %d", fw.Code)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
