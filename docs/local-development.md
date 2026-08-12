# Local development

Use the existing Compose files without changing unrelated Docker resources. Development PostgreSQL is `rmp-phase2-dev-db` on localhost:5433. Copy only placeholder values from `.env.example` into the process environment; do not place production credentials in this repository. Start the backend from `backend/` with `go run ./cmd/server` after setting `DATABASE_URL`, `PORT`, and a valid absolute `ALLOWED_ORIGIN` such as `http://localhost:3000`.

The separate test Compose project is `rmp-phase2-test` and only uses localhost:5434. It must never use port 5432 or the development service's port 5433.

## Locked deployment assumptions for later backend phases

The initial public deployment is one API instance, so the current in-memory
five-attempts/IP/minute limiter remains acceptable until horizontal scaling is
introduced. The server must use `RemoteAddr` by default. Proxy headers such as
`Forwarded` and `X-Forwarded-For` may only be trusted after explicit trusted
proxy infrastructure/configuration is introduced; they must never be trusted
blindly from public clients. A distributed/shared limiter is required before
multiple API instances serve traffic.

The future admin credential is environment-provided and header-based; never
place credentials in a query string or local documentation. See
[backend-phase-b1-contract.md](backend-phase-b1-contract.md) for the locked
contract and legacy-consent policy.
