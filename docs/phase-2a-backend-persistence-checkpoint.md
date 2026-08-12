# Phase 2A backend/persistence checkpoint

## Implemented foundation

The Go/Chi service provides `POST /api/waitlist` and liveness-only `GET /health`. The write endpoint accepts only name/email JSON, trims values, lowercases email, validates lengths/syntax, rejects unknown/trailing JSON, writes parameterized PostgreSQL data, maps unique conflicts to 409, and has a five-attempts-per-IP-per-minute in-memory rate limit. Error responses use `{ "error": { "code", "message" } }` without database detail.

Frontend origin is configured by `NEXT_PUBLIC_API_BASE_URL`; the backend reads `ALLOWED_ORIGIN`. CORS is explicit (no wildcard), permits only `POST, OPTIONS` and `Content-Type`, and does not enable credentials. Production can omit the public base URL for same-origin reverse-proxy routing; local values are placeholders in `.env.example`.

## Database and lifecycle

PostgreSQL 16 development is repository-owned at localhost:5433. Disposable repository-owned test PostgreSQL 16 uses localhost:5434. The test lifecycle was verified on the isolated database: health passed; up migration created the extension/table; down migration dropped the table; up migration created it again. Development data and unrelated port-5432 resources were untouched.

## Checks

`gofmt`, `go vet ./...`, `go test ./...`, and `go build ./cmd/server` passed. Compose test configuration validation passed.

## Remaining Phase 2B work

Dedicated Go handler/integration tests, migration metadata tooling, frontend state tests, real frontend-to-Go execution, Playwright, Chrome DevTools, final documentation reconciliation, and final Phase 2 verdict remain. This checkpoint does not implement Phase 3 functionality.
