package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

type fakeStore struct {
	mu    sync.Mutex
	calls int
	// name is `any` (not string) because the handler passes a real Go nil
	// for the SQL parameter when no name was supplied, so it can be
	// distinguished from an explicit empty string.
	name  any
	email string
	err   error
}

func (s *fakeStore) Exec(_ context.Context, _ string, args ...any) (pgconn.CommandTag, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls++
	s.name = args[0]
	s.email = args[1].(string)
	return pgconn.CommandTag{}, s.err
}
func routerFor(s *fakeStore) http.Handler {
	return newRouter(s, "http://localhost:3000", newLimiter(func() time.Time { return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) }))
}
func post(r http.Handler, body string, origin string) *httptest.ResponseRecorder {
	q := httptest.NewRequest(http.MethodPost, "/api/waitlist", strings.NewReader(body))
	q.RemoteAddr = "203.0.113.10:1234"
	if origin != "" {
		q.Header.Set("Origin", origin)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, q)
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
	var v errorBody
	if err := json.Unmarshal(w.Body.Bytes(), &v); err != nil {
		t.Fatal(err)
	}
	return v.Error.Code
}

func TestWaitlistSuccessNormalizesAndPersistsOnce(t *testing.T) {
	s := &fakeStore{}
	w := post(routerFor(s), `{"name":"  Ada Lovelace ","email":" ADA@EXAMPLE.TEST "}`, "")
	if w.Code != 201 {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	requireJSON(t, w)
	if w.Body.String() != `{"status":"created"}` {
		t.Fatalf("body=%s", w.Body.String())
	}
	if s.calls != 1 || s.name != "Ada Lovelace" || s.email != "ada@example.test" {
		t.Fatalf("store %#v", s)
	}
}
func TestWaitlistAcceptsEmailOnly(t *testing.T) {
	s := &fakeStore{}
	w := post(routerFor(s), `{"email":" Ada@Example.Test "}`, "")
	if w.Code != 201 {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	requireJSON(t, w)
	if w.Body.String() != `{"status":"created"}` {
		t.Fatalf("body=%s", w.Body.String())
	}
	if s.calls != 1 || s.name != nil || s.email != "ada@example.test" {
		t.Fatalf("store %#v", s)
	}
}
func TestWaitlistWhitespaceOnlyNameStoredAsNull(t *testing.T) {
	s := &fakeStore{}
	w := post(routerFor(s), `{"name":"   ","email":"grace@example.test"}`, "")
	if w.Code != 201 {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	if s.calls != 1 || s.name != nil || s.email != "grace@example.test" {
		t.Fatalf("store %#v", s)
	}
}
func TestWaitlistValidationRejectsUnsafePayloads(t *testing.T) {
	cases := []string{`{}`, `{"name":"A"}`, `{"name":"A","email":" "}`, `{"name":"A","email":"invalid"}`, `{"name":"A","email":"a@example.test","extra":true}`, `{`, `{"name":"A","email":"a@example.test"} {}`, `{"name":"A","email":"a@example.test"}{"name":"B"}`, `[]`, `{"email":"invalid"}`, `{"email":"a@example.test","extra":true}`, `{"email":" "}`}
	cases = append(cases, `{"name":"`+strings.Repeat("a", 101)+`","email":"a@example.test"}`, `{"name":"A","email":"`+strings.Repeat("a", 246)+`@example.test"}`)
	for _, body := range cases {
		t.Run(body[:min(16, len(body))], func(t *testing.T) {
			s := &fakeStore{}
			w := post(routerFor(s), body, "")
			if w.Code != 400 {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			requireJSON(t, w)
			if errorCode(t, w) != validationCode {
				t.Fatalf("code=%s", errorCode(t, w))
			}
			if s.calls != 0 {
				t.Fatal("store was called")
			}
			for _, unsafe := range []string{"sql", "stack", "path", "postgres"} {
				if strings.Contains(strings.ToLower(w.Body.String()), unsafe) {
					t.Fatalf("unsafe body %s", w.Body.String())
				}
			}
		})
	}
}
func TestWaitlistMapsDuplicateAndInternalErrors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		err    error
		status int
		code   string
	}{{"duplicate", &pgconn.PgError{Code: "23505", Detail: "Key (email_normalized)=(a@example.test) already exists"}, 409, "WAITLIST_EMAIL_EXISTS"}, {"internal", errors.New("postgres://secret@db/remember_my_pill /srv/backend/main.go"), 500, internalCode}} {
		t.Run(tc.name, func(t *testing.T) {
			w := post(routerFor(&fakeStore{err: tc.err}), `{"name":"Ada","email":"Ada@Example.Test"}`, "")
			if w.Code != tc.status || errorCode(t, w) != tc.code {
				t.Fatalf("status=%d code=%s", w.Code, errorCode(t, w))
			}
			if strings.Contains(w.Body.String(), "secret") || strings.Contains(w.Body.String(), "postgres") || strings.Contains(w.Body.String(), "/srv") {
				t.Fatalf("leak: %s", w.Body.String())
			}
		})
	}
}
func TestRateLimitAndHealth(t *testing.T) {
	s := &fakeStore{}
	r := routerFor(s)
	for i := 0; i < 5; i++ {
		if w := post(r, `{"name":"Ada","email":"ada`+string(rune('a'+i))+`@example.test"}`, ""); w.Code != 201 {
			t.Fatalf("request %d=%d", i, w.Code)
		}
	}
	w := post(r, `{"name":"Ada","email":"later@example.test"}`, "")
	if w.Code != 429 || errorCode(t, w) != "WAITLIST_RATE_LIMITED" {
		t.Fatalf("rate response=%d %s", w.Code, w.Body.String())
	}
	q := httptest.NewRequest(http.MethodGet, "/health", nil)
	q.RemoteAddr = "203.0.113.10:999"
	h := httptest.NewRecorder()
	r.ServeHTTP(h, q)
	if h.Code != 200 || h.Body.String() != `{"status":"ok"}` {
		t.Fatalf("health %d %s", h.Code, h.Body.String())
	}
	q2 := httptest.NewRequest(http.MethodPost, "/api/waitlist", strings.NewReader(`{"name":"Ada","email":"different@example.test"}`))
	q2.RemoteAddr = "203.0.113.11:999"
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, q2)
	if w2.Code != 201 {
		t.Fatalf("separate ip=%d", w2.Code)
	}
}
func TestHealthContract(t *testing.T) {
	w := httptest.NewRecorder()
	routerFor(&fakeStore{}).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/health", nil))
	if w.Code != 200 || w.Body.String() != `{"status":"ok"}` {
		t.Fatalf("%d %s", w.Code, w.Body.String())
	}
	requireJSON(t, w)
	for _, bad := range []string{"postgres", "database", "password", "path", "stack"} {
		if strings.Contains(strings.ToLower(w.Body.String()), bad) {
			t.Fatal("health leak")
		}
	}
}
func TestCORS(t *testing.T) {
	r := routerFor(&fakeStore{})
	pre := httptest.NewRequest(http.MethodOptions, "/api/waitlist", nil)
	pre.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, pre)
	if w.Code != 204 || w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" || w.Header().Get("Access-Control-Allow-Methods") != "POST, OPTIONS" || w.Header().Get("Access-Control-Allow-Headers") != "Content-Type" {
		t.Fatalf("allowed cors: %d %#v", w.Code, w.Header())
	}
	if w.Header().Get("Access-Control-Allow-Credentials") != "" {
		t.Fatal("credentials enabled")
	}
	bad := httptest.NewRequest(http.MethodOptions, "/api/waitlist", nil)
	bad.Header.Set("Origin", "https://evil.example")
	bw := httptest.NewRecorder()
	r.ServeHTTP(bw, bad)
	if bw.Code != 403 || bw.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("bad cors %d %#v", bw.Code, bw.Header())
	}
	same := post(r, `{"name":"Ada","email":"ada@example.test"}`, "")
	if same.Code != 201 {
		t.Fatalf("same origin client %d", same.Code)
	}
}
func TestAllowedOriginValidation(t *testing.T) {
	for _, v := range []string{"", "localhost:3000", "ftp://example.test", "http://example.test/path", "https://example.test"} {
		want := v == "https://example.test"
		if got := validAllowedOrigin(v); got != want {
			t.Fatalf("%q=%v", v, got)
		}
	}
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
