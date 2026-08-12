# Local development

Use the existing Compose files without changing unrelated Docker resources. Development PostgreSQL is `rmp-phase2-dev-db` on localhost:5433. Copy only placeholder values from `.env.example` into the process environment; do not place production credentials in this repository. Start the backend from `backend/` with `go run ./cmd/server` after setting `DATABASE_URL`, `PORT`, and a valid absolute `ALLOWED_ORIGIN` such as `http://localhost:3000`.

The separate test Compose project is `rmp-phase2-test` and only uses localhost:5434. It must never use port 5432 or the development service's port 5433.
