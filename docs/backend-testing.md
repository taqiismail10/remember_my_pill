# Backend testing

Run unit and handler tests from `backend/` with `go test ./...`; format with
`gofmt -d cmd/server/*.go`, vet with `go vet ./...`, and make the production
binary with `go build ./cmd/server`.

Integration tests are tagged and require `TEST_DATABASE_URL` for the disposable
PostgreSQL database named `remember_my_pill_test`. The intended command is
`go test -tags=integration ./...` after `docker compose -f
docker-compose.test.yml -p rmp-phase2-test up -d`.
Tear down only that project with
`docker compose -f docker-compose.test.yml -p rmp-phase2-test down -v`.

The canonical migration harness applies the embedded ordered sequence
`001 → 002 → 003 → 004 → 005`, validates legacy pre-consent rows, and checks down/up
reversal before resetting the disposable database for handler coverage.

## Later acceptance coverage

Before B2–B6 are declared complete, tests must cover:

- email-only compatibility during the migration window;
- true, missing, false, and stale-version consent;
- malformed/missing email, unknown fields, trailing JSON, and oversized body;
- populated honeypot without entry creation or disclosure;
- duplicate and concurrent duplicate submissions;
- invalid/valid/self/concurrent referrals and referral-code collision retry;
- status-token lookup and invalid-token rejection;
- deterministic signup-order rank and separately derived referral count;
- admin authorization, CSV escaping, and no public admin CORS;
- readiness with available/unavailable PostgreSQL;
- trusted and untrusted proxy-header behavior; and
- complete ordered migration lifecycle `001 → 002 → 003 → 004 → 005`, including
  paired down migrations, against disposable PostgreSQL.

Migration 004 is implemented for B3A internal schema coverage only. Tests must
continue to verify its referral and verification-token constraints without
enabling public B3 routes or real email.

B3B adapter tests use `EMAIL_PROVIDER=fake` or a local mocked Postmark HTTP
server. `go test ./...` must never call Postmark. A manual sandbox smoke test
may use the provider's documented test token only when explicitly requested;
it is not part of automated verification.

B2B-Prep tests cover independent required and marketing consent, compatibility
while the legal-content flag is off, enforcement when it is on, and migration
005 rollback/reapply. They do not send email or enable B3 routes.

On this Windows host, `go test -race ./...` is not executable because `cc1.exe`
reports `64-bit mode not compiled in`. Use a 64-bit-capable C toolchain in a
future environment; do not treat the race suite as passing here.
