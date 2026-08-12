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
`001 → 002 → 003`, validates legacy pre-consent rows, and checks down/up
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
- complete ordered migration lifecycle `001 → 002 → 003 → 004`, including
  paired down migrations, against disposable PostgreSQL.

Migration 004 is not yet implemented. Add it to the one ordered migration list
when B3 introduces the referral-event schema.

On this Windows host, `go test -race ./...` is not executable because `cc1.exe`
reports `64-bit mode not compiled in`. Use a 64-bit-capable C toolchain in a
future environment; do not treat the race suite as passing here.
