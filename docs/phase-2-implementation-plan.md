# Phase 2 implementation plan — minimal waitlist vertical slice

**Status:** Planned; blocked pending a usable local PostgreSQL test environment.

## Objective and user value

Allow a visitor to submit only a name and email through an accessible RMP waitlist form and receive a truthful success or friendly failure result backed by the PRD-required Go API and PostgreSQL persistence.

## PRD scope and revised boundary

Addresses PRD §§3.1, 5, 7.1, 8 (waitlist), 12.3–12.4, 13.1/13.3 (minimum applicable validation and rate limit), 14, 16–18 (minimal architecture/data model), 17 (`POST /api/waitlist`, `GET /health`), 22, 24, 26.2–26.4, and 28.

The previous roadmap placed all services in Phase 3. The approved revision moves only the minimal Go/PostgreSQL persistence needed to prove a genuine Phase 2 conversion flow. Phase 3 explicitly excludes this milestone and retains referral codes, attribution, rank, rewards, status lookup, sharing, email, admin export, and production-oriented operations.

## Figma and Phase 1 dependencies

Use palette `6026:3`, app-icon/identity `6021:115`, and R/pill artwork `6021:128`. Reuse Phase 1 CSS tokens, Manrope, logo, buttons, cards, semantic structure, focus treatment, and reduced-motion rules.

## Included scope

- Responsive name/email form with client validation, loading/disabled, success, duplicate, server, network, retry, and accessible live-status states.
- Go + Chi API: `POST /api/waitlist` and `GET /health`.
- PostgreSQL `waitlist_entries` migration with UUID id, bounded name, unique normalized email, and server timestamp.
- Trimmed/lowercased email, bounded field validation, malformed JSON handling, non-sensitive JSON errors, duplicate conflict handling, and 5 IP attempts/minute in-memory rate limit.
- `.env.example`, local Compose configuration, API/database/startup/migration documentation, and frontend API-base configuration.
- Unit, API, real disposable-PostgreSQL integration, and Playwright tests using fictional `example.com` data.

## Excluded Phase 3 scope

Referral, rank, rewards, status tokens/endpoints, sharing, email, admin authorization/export, billing, prescription upload/AI, health data, production deployment, and advanced abuse controls.

## API contract

`POST /api/waitlist` accepts `{ "name": "Alice Example", "email": "alice@example.com" }` and returns 201 `{ "status": "created" }`. It returns 400 `WAITLIST_VALIDATION_ERROR`, 409 `WAITLIST_EMAIL_EXISTS`, 429 `WAITLIST_RATE_LIMITED`, or 500 `WAITLIST_INTERNAL_ERROR`, all with `{ "error": { "code", "message" } }` and no internal detail. `GET /health` returns a minimal service status.

## Migration and security

Use an additive reversible up/down SQL migration, unique normalized-email constraint, parameterized queries, and transaction/constraint-driven duplicate safety. No public database port is published. Name/email are the only accepted fields; request bodies and database errors are not logged. The API base URL is configured by environment and no secret is exposed to the client.

## Expected files

`backend/go.mod`, `backend/cmd/server/main.go`, `backend/internal/*`, `backend/migrations/*`, `backend/tests/*`, `docker-compose.yml`, `.env.example`, `frontend/components/waitlist/*`, `frontend/lib/*`, `frontend/app/page.tsx`, frontend tests, API/database/environment/testing documentation, and `docs/evidence/phase-2/*`.

## Sequence and acceptance criteria

1. Start a disposable PostgreSQL environment and prove migration application.
2. Implement/test Go validation, rate limit, persistence, health, and migration.
3. Implement the form and API client against that real service.
4. Test success, invalid input, duplicate normalized email, API/database failure, network failure, loading, retry, keyboard flow, all required widths, and refresh-after-success behavior.
5. Verify browser console/network and Figma consistency; then update status docs only with evidence.

Acceptance requires a real frontend → Go → PostgreSQL successful submission, durable duplicate detection after restart, correct non-sensitive status responses, rate-limit behavior, passing checks, and sanitised evidence.

## Blocker and rollback

Go is installed, but `docker version` cannot connect to Docker Desktop’s Linux engine, so a disposable PostgreSQL database and required real integration test cannot be run. No implementation starts until that is available. All planned changes are additive; rollback is to stop local services and apply the migration down script before any non-local deployment.
