//go:build integration

package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"remember_my_pill/backend/internal/httpapi"
	"remember_my_pill/backend/migrations"
)

var integrationPool *pgxpool.Pool

func TestMain(m *testing.M) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		os.Exit(0)
	}
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		panic(err)
	}
	if config.ConnConfig.Database != "remember_my_pill_test" {
		panic("integration tests require the disposable remember_my_pill_test database")
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		panic(err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		panic(err)
	}
	integrationPool = pool
	code := m.Run()
	pool.Close()
	os.Exit(code)
}

func orderedMigrations(t *testing.T) []migrations.Migration {
	t.Helper()
	items, err := migrations.Ordered()
	if err != nil {
		t.Fatal(err)
	}
	return items
}

func dropWaitlist(t *testing.T) {
	t.Helper()
	if _, err := integrationPool.Exec(context.Background(), "DROP TABLE IF EXISTS waitlist_entries CASCADE"); err != nil {
		t.Fatal(err)
	}
}

func applyUp(t *testing.T, items []migrations.Migration) {
	t.Helper()
	for _, migration := range items {
		if _, err := integrationPool.Exec(context.Background(), migration.Up); err != nil {
			t.Fatalf("apply %s: %v", migration.Name, err)
		}
	}
}

func applyDown(t *testing.T, items []migrations.Migration) {
	t.Helper()
	for i := len(items) - 1; i >= 0; i-- {
		if _, err := integrationPool.Exec(context.Background(), items[i].Down); err != nil {
			t.Fatalf("revert %s: %v", items[i].Name, err)
		}
	}
}

func resetDatabase(t *testing.T) {
	t.Helper()
	dropWaitlist(t)
	applyUp(t, orderedMigrations(t))
}

func integrationRouter(store httpapi.Store, readiness httpapi.Readiness) http.Handler {
	return httpapi.NewRouter(store, readiness, httpapi.Options{
		AllowedOrigin:  "http://localhost:3000",
		ConsentVersion: "policy-2026-08",
		Logger:         slog.New(slog.NewTextHandler(ioDiscard{}, nil)),
	})
}

func integrationPost(r http.Handler, body, ip string) *httptest.ResponseRecorder {
	q := httptest.NewRequest(http.MethodPost, "/api/waitlist", strings.NewReader(body))
	q.RemoteAddr = ip + ":1234"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, q)
	return w
}

func TestIntegrationMigrationLifecycleAndLegacyConsent(t *testing.T) {
	items := orderedMigrations(t)
	dropWaitlist(t)
	applyUp(t, items)
	applyDown(t, items)
	applyUp(t, items)
	var exists bool
	if err := integrationPool.QueryRow(context.Background(), `SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'waitlist_entries' AND column_name = 'updated_at')`).Scan(&exists); err != nil || !exists {
		t.Fatalf("full up/down/up did not restore migration 003 schema: exists=%v err=%v", exists, err)
	}

	// Exercise the legacy path separately. Migration 002 intentionally makes
	// names nullable, so its historical down migration cannot restore NOT NULL
	// while a valid legacy NULL name exists without deleting or inventing data.
	dropWaitlist(t)
	applyUp(t, items[:2])
	if _, err := integrationPool.Exec(context.Background(), `INSERT INTO waitlist_entries (name, email_normalized) VALUES (NULL, 'legacy@example.test')`); err != nil {
		t.Fatal(err)
	}
	applyUp(t, items[2:])

	var consentVersion *string
	var consentedAt *time.Time
	var referralCode, statusToken *string
	if err := integrationPool.QueryRow(context.Background(), `SELECT consent_version, consented_at, referral_code, status_token_hash FROM waitlist_entries WHERE email_normalized = 'legacy@example.test'`).Scan(&consentVersion, &consentedAt, &referralCode, &statusToken); err != nil {
		t.Fatal(err)
	}
	if consentVersion != nil || consentedAt != nil || referralCode != nil || statusToken != nil {
		t.Fatalf("legacy record was fabricated: consent=%v consented_at=%v referral=%v token=%v", consentVersion, consentedAt, referralCode, statusToken)
	}

	applyDown(t, items[2:])
	applyUp(t, items[2:])
	if err := integrationPool.QueryRow(context.Background(), `SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'waitlist_entries' AND column_name = 'updated_at')`).Scan(&exists); err != nil || !exists {
		t.Fatalf("migration 003 down/up did not restore schema: exists=%v err=%v", exists, err)
	}
}

func TestIntegrationPersistenceDuplicateConsentAndConstraints(t *testing.T) {
	resetDatabase(t)
	r := integrationRouter(integrationPool, integrationPool)
	w := integrationPost(r, `{"name":"  Grace Hopper  ","email":" GRACE@EXAMPLE.TEST ","consent":true,"consentVersion":"policy-2026-08"}`, "203.0.113.1")
	requireAccepted(t, w)

	var id, name, email string
	var created, consentedAt, updatedAt time.Time
	var consentVersion string
	if err := integrationPool.QueryRow(context.Background(), `SELECT id::text, name, email_normalized, created_at, consent_version, consented_at, updated_at FROM waitlist_entries WHERE email_normalized = 'grace@example.test'`).Scan(&id, &name, &email, &created, &consentVersion, &consentedAt, &updatedAt); err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`).MatchString(id) || name != "Grace Hopper" || email != "grace@example.test" || created.IsZero() || consentVersion != "policy-2026-08" || consentedAt.IsZero() || updatedAt.IsZero() {
		t.Fatalf("stored id=%s name=%q email=%q created=%v consent=%q consentedAt=%v updatedAt=%v", id, name, email, created, consentVersion, consentedAt, updatedAt)
	}
	duplicate := integrationPost(r, `{"email":"GrAcE@Example.Test"}`, "203.0.113.2")
	requireAccepted(t, duplicate)
	var count int
	if err := integrationPool.QueryRow(context.Background(), `SELECT count(*) FROM waitlist_entries WHERE email_normalized = 'grace@example.test'`).Scan(&count); err != nil || count != 1 {
		t.Fatalf("count=%d err=%v", count, err)
	}
	if _, err := integrationPool.Exec(context.Background(), `INSERT INTO waitlist_entries (name, email_normalized) VALUES (NULL, 'null-name@example.test')`); err != nil {
		t.Fatalf("NULL name should be allowed: %v", err)
	}
	for _, query := range []string{
		`INSERT INTO waitlist_entries (name,email_normalized) VALUES ('valid',NULL)`,
		`INSERT INTO waitlist_entries (name,email_normalized) VALUES ('` + strings.Repeat("a", 101) + `','long-name@example.test')`,
		`INSERT INTO waitlist_entries (name,email_normalized) VALUES ('valid','` + strings.Repeat("a", 246) + `@example.test')`,
		`INSERT INTO waitlist_entries (name,email_normalized) VALUES ('again','grace@example.test')`,
	} {
		if _, err := integrationPool.Exec(context.Background(), query); err == nil {
			t.Fatalf("constraint did not reject %q", query)
		}
	}
}

func TestIntegrationEmailOnlyHoneypotAndReadiness(t *testing.T) {
	resetDatabase(t)
	r := integrationRouter(integrationPool, integrationPool)
	ready := httptest.NewRecorder()
	r.ServeHTTP(ready, httptest.NewRequest(http.MethodGet, "/ready", nil))
	if ready.Code != http.StatusOK || ready.Body.String() != `{"status":"ready"}` {
		t.Fatalf("ready=%d %s", ready.Code, ready.Body.String())
	}

	requireAccepted(t, integrationPost(r, `{"email":"compatibility@example.test"}`, "203.0.113.3"))
	var name, consentVersion *string
	var consentedAt *time.Time
	if err := integrationPool.QueryRow(context.Background(), `SELECT name, consent_version, consented_at FROM waitlist_entries WHERE email_normalized = 'compatibility@example.test'`).Scan(&name, &consentVersion, &consentedAt); err != nil {
		t.Fatal(err)
	}
	if name != nil || consentVersion != nil || consentedAt != nil {
		t.Fatalf("compatibility record should remain unconsented: name=%v consent=%v consented_at=%v", name, consentVersion, consentedAt)
	}

	requireAccepted(t, integrationPost(r, `{"email":"bot@example.test","company":"Bot Inc"}`, "203.0.113.4"))
	var bots int
	if err := integrationPool.QueryRow(context.Background(), `SELECT count(*) FROM waitlist_entries WHERE email_normalized = 'bot@example.test'`).Scan(&bots); err != nil || bots != 0 {
		t.Fatalf("honeypot rows=%d err=%v", bots, err)
	}

	unavailable := integrationRouter(errorStore{err: errors.New("database unavailable")}, errorReadiness{err: errors.New("database unavailable")})
	failedReady := httptest.NewRecorder()
	unavailable.ServeHTTP(failedReady, httptest.NewRequest(http.MethodGet, "/ready", nil))
	if failedReady.Code != http.StatusServiceUnavailable || strings.Contains(strings.ToLower(failedReady.Body.String()), "database") {
		t.Fatalf("unavailable readiness=%d %s", failedReady.Code, failedReady.Body.String())
	}
}

func TestIntegrationConcurrentDuplicateAndUnavailableStore(t *testing.T) {
	resetDatabase(t)
	r := integrationRouter(integrationPool, integrationPool)
	const attempts = 10
	results := make(chan int, attempts)
	var wg sync.WaitGroup
	for i := 0; i < attempts; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results <- integrationPost(r, `{"email":"Concurrent@Example.Test"}`, "198.51.100."+string(rune('1'+i))).Code
		}(i)
	}
	wg.Wait()
	close(results)
	for status := range results {
		if status != http.StatusAccepted {
			t.Fatalf("concurrent response = %d", status)
		}
	}
	var count int
	if err := integrationPool.QueryRow(context.Background(), `SELECT count(*) FROM waitlist_entries WHERE email_normalized = 'concurrent@example.test'`).Scan(&count); err != nil || count != 1 {
		t.Fatalf("count=%d err=%v", count, err)
	}

	unavailable := integrationRouter(errorStore{err: errors.New("database unavailable")}, errorReadiness{err: errors.New("database unavailable")})
	w := integrationPost(unavailable, `{"email":"unavailable@example.test"}`, "203.0.113.99")
	if w.Code != http.StatusInternalServerError || errorCode(t, w) != httpapi.InternalCode || strings.Contains(strings.ToLower(w.Body.String()), "database") {
		t.Fatalf("unavailable %d %s", w.Code, w.Body.String())
	}
}

type errorStore struct{ err error }

func (s errorStore) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, s.err
}

type errorReadiness struct{ err error }

func (r errorReadiness) Ping(context.Context) error { return r.err }
