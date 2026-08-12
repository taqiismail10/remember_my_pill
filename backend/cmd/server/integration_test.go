//go:build integration

package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var integrationPool *pgxpool.Pool

func TestMain(m *testing.M) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		os.Exit(0)
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		panic(err)
	}
	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}
	integrationPool = pool
	code := m.Run()
	pool.Close()
	os.Exit(code)
}
func migration(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "migrations", name))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
func resetDatabase(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	_, err := integrationPool.Exec(ctx, migration(t, "001_waitlist_entries.down.sql"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = integrationPool.Exec(ctx, migration(t, "001_waitlist_entries.up.sql")); err != nil {
		t.Fatal(err)
	}
}
func integrationRouter() http.Handler {
	return newRouter(integrationPool, "http://localhost:3000", newLimiter(nil))
}
func integrationPost(r http.Handler, name, email, ip string) *httptest.ResponseRecorder {
	q := httptest.NewRequest(http.MethodPost, "/api/waitlist", strings.NewReader(`{"name":"`+name+`","email":"`+email+`"}`))
	q.RemoteAddr = ip + ":1234"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, q)
	return w
}

func TestIntegrationMigrationLifecycleAndSchema(t *testing.T) {
	ctx := context.Background()
	_, _ = integrationPool.Exec(ctx, migration(t, "001_waitlist_entries.down.sql"))
	if _, err := integrationPool.Exec(ctx, migration(t, "001_waitlist_entries.up.sql")); err != nil {
		t.Fatal(err)
	}
	var exists bool
	if err := integrationPool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name='waitlist_entries')`).Scan(&exists); err != nil || !exists {
		t.Fatalf("table exists=%v err=%v", exists, err)
	}
	var prohibited int
	if err := integrationPool.QueryRow(ctx, `SELECT count(*) FROM information_schema.columns WHERE table_name='waitlist_entries' AND column_name IN ('health_data','referral_code','referral_count')`).Scan(&prohibited); err != nil || prohibited != 0 {
		t.Fatalf("prohibited=%d err=%v", prohibited, err)
	}
	if _, err := integrationPool.Exec(ctx, migration(t, "001_waitlist_entries.down.sql")); err != nil {
		t.Fatal(err)
	}
	var after bool
	if err := integrationPool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name='waitlist_entries')`).Scan(&after); err != nil || after {
		t.Fatalf("removed=%v err=%v", after, err)
	}
	if _, err := integrationPool.Exec(ctx, migration(t, "001_waitlist_entries.up.sql")); err != nil {
		t.Fatal(err)
	}
}
func TestIntegrationPersistenceDuplicateAndConstraints(t *testing.T) {
	resetDatabase(t)
	r := integrationRouter()
	w := integrationPost(r, "  Grace Hopper  ", " GRACE@EXAMPLE.TEST ", "203.0.113.1")
	if w.Code != 201 {
		t.Fatalf("first %d %s", w.Code, w.Body.String())
	}
	var id string
	var name, email string
	var created time.Time
	if err := integrationPool.QueryRow(context.Background(), `SELECT id::text,name,email_normalized,created_at FROM waitlist_entries`).Scan(&id, &name, &email, &created); err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`).MatchString(id) || name != "Grace Hopper" || email != "grace@example.test" || created.IsZero() {
		t.Fatalf("stored id=%s name=%q email=%q created=%v", id, name, email, created)
	}
	for _, address := range []string{"grace@example.test", "GrAcE@Example.Test"} {
		w = integrationPost(r, "Grace", address, "203.0.113."+string(rune('2'+len(address)%5)))
		if w.Code != 409 || errorCode(t, w) != "WAITLIST_EMAIL_EXISTS" {
			t.Fatalf("dup %d %s", w.Code, w.Body.String())
		}
	}
	var count int
	if err := integrationPool.QueryRow(context.Background(), `SELECT count(*) FROM waitlist_entries`).Scan(&count); err != nil || count != 1 {
		t.Fatalf("count=%d err=%v", count, err)
	}
	// name is nullable (migration 002): a NULL name must be accepted, not
	// rejected, at the database level.
	if _, err := integrationPool.Exec(context.Background(), `INSERT INTO waitlist_entries (name,email_normalized) VALUES (NULL,'null-name@example.test')`); err != nil {
		t.Fatalf("NULL name should be allowed: %v", err)
	}
	for _, q := range []string{`INSERT INTO waitlist_entries (name,email_normalized) VALUES ('valid',NULL)`, `INSERT INTO waitlist_entries (name,email_normalized) VALUES ('` + strings.Repeat("a", 101) + `','long-name@example.test')`, `INSERT INTO waitlist_entries (name,email_normalized) VALUES ('valid','` + strings.Repeat("a", 246) + `@example.test')`, `INSERT INTO waitlist_entries (name,email_normalized) VALUES ('again','grace@example.test')`} {
		if _, err := integrationPool.Exec(context.Background(), q); err == nil {
			t.Fatalf("constraint did not reject %q", q)
		}
	}
}
func TestIntegrationEmailOnlySignupStoresNullName(t *testing.T) {
	resetDatabase(t)
	r := integrationRouter()
	q := httptest.NewRequest(http.MethodPost, "/api/waitlist", strings.NewReader(`{"email":"nameless@example.test"}`))
	q.RemoteAddr = "203.0.113.50:1234"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, q)
	if w.Code != 201 {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	var name *string
	if err := integrationPool.QueryRow(context.Background(), `SELECT name FROM waitlist_entries WHERE email_normalized='nameless@example.test'`).Scan(&name); err != nil {
		t.Fatal(err)
	}
	if name != nil {
		t.Fatalf("expected NULL name, got %q", *name)
	}
}
func TestIntegrationConcurrentDuplicateAndUnavailableStore(t *testing.T) {
	resetDatabase(t)
	r := integrationRouter()
	const n = 10
	results := make(chan int, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results <- integrationPost(r, "Concurrent", "Concurrent@Example.Test", "198.51.100."+string(rune('1'+i))).Code
		}(i)
	}
	wg.Wait()
	close(results)
	created, conflict, other := 0, 0, 0
	for status := range results {
		switch status {
		case 201:
			created++
		case 409:
			conflict++
		default:
			other++
		}
	}
	if created != 1 || conflict != n-1 || other != 0 {
		t.Fatalf("created=%d conflict=%d other=%d", created, conflict, other)
	}
	var count int
	if err := integrationPool.QueryRow(context.Background(), `SELECT count(*) FROM waitlist_entries WHERE email_normalized='concurrent@example.test'`).Scan(&count); err != nil || count != 1 {
		t.Fatalf("count=%d err=%v", count, err)
	}
	config, err := pgxpool.ParseConfig("postgres://rmp_test:rmp_test_only@127.0.0.1:1/unavailable?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	config.ConnConfig.ConnectTimeout = 100 * time.Millisecond
	unavailable, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	defer unavailable.Close()
	w := integrationPost(newRouter(unavailable, "http://localhost:3000", newLimiter(nil)), "Unavailable", "unavailable@example.test", "203.0.113.99")
	if w.Code != 500 || errorCode(t, w) != internalCode || strings.Contains(strings.ToLower(w.Body.String()), "connect") {
		t.Fatalf("unavailable %d %s", w.Code, w.Body.String())
	}
}
