# Phase 2C backend test checkpoint

## Baseline and focused testability refactor

The existing Go/Chi server retained its architecture. Handler construction now accepts a small PostgreSQL `Exec` interface, the limiter accepts an injectable clock, CORS is a dedicated middleware function, and startup rejects missing or malformed `ALLOWED_ORIGIN`. These focused changes allow deterministic full-router tests without a broad backend rewrite.

## Passed handler coverage

`backend/cmd/server/main_test.go` verifies successful `POST /api/waitlist` behavior (201, JSON contract, trim/lowercase normalization, exactly one insertion); missing, whitespace, malformed, oversized, unknown-field, malformed, trailing, multi-document, and non-object JSON validation failures (400, no insertion, safe JSON error); PostgreSQL unique violation mapping (409); generic repository failure mapping (500 without connection detail); five-write rate limiting, independent IP buckets, and health exemption; exact liveness response `{"status":"ok"}`; and CORS through the real router (exact allowlist, 204 preflight, restricted method/header set, no wildcard or credentials, denied origin, and no-Origin client). Invalid writes intentionally consume capacity because limiting occurs before parsing.

`gofmt -d`, `go vet ./...`, `go test ./...`, and `go build ./cmd/server` passed. `go test -race ./...` was attempted but is **NOT EXECUTED — WINDOWS C TOOLCHAIN LIMITATION**: `cc1.exe: sorry, unimplemented: 64-bit mode not compiled in`. A future 64-bit-capable C toolchain is required.

## Disposable PostgreSQL and integration status

`docker-compose.test.yml` defines PostgreSQL 16, repository-owned container/network/volume names, localhost-only port 5434, health checking, distinct test credentials, and deterministic `down -v` cleanup. A tagged real integration suite was added in `backend/cmd/server/integration_test.go`; it automates migration down/up/down/up, schema checks, real HTTP persistence, direct SQL constraints, concurrent duplicate handling, and unavailable-connection behavior.

This host currently reserves Windows TCP ports 5434–5533. Docker configuration validation passed, but startup failed binding `127.0.0.1:5434` with `ports are not available`; no service listened on 5434. The created repository-owned container, network, and volume were immediately removed with `docker compose -f docker-compose.test.yml -p rmp-phase2-test down -v`. Therefore the real integration suite, lifecycle execution, repeatability execution, and live database-unavailable execution remain unverified and must not be marked passed. Port 5433 development PostgreSQL and unrelated port 5432 resources were preserved.

## Defects fixed and remaining work

The refactor replaced text matching for duplicate errors with PostgreSQL SQLSTATE `23505`, made required origin configuration fail safely at startup, made health output exact, and stopped silently discarding `ListenAndServe` errors. No Phase 3 functionality was added.

Phase 2D still requires a host where isolated 5434 can bind, real integration execution, frontend/API flow tests, Playwright, Chrome DevTools, and final Phase 2 reconciliation. Phase 2 remains incomplete.
