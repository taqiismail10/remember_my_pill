# Backend testing

Run unit and handler tests from `backend/` with `go test ./...`; format with `gofmt -d cmd/server/*.go`, vet with `go vet ./...`, and make the production binary with `go build ./cmd/server`.

Integration tests are deliberately tagged and require `TEST_DATABASE_URL` for the repository-owned disposable PostgreSQL database. The intended command is `go test -tags=integration ./...` after `docker compose -f docker-compose.test.yml -p rmp-phase2-test up -d`. The suite owns its migration reset and exercises up/down/up lifecycle, real persistence, duplicate concurrency, schema constraints, and an unavailable connection. Tear down only that project with `docker compose -f docker-compose.test.yml -p rmp-phase2-test down -v`.

On this Windows host, `go test -race ./...` is not executable because `cc1.exe` reports `64-bit mode not compiled in`. Use a 64-bit-capable C toolchain in a future environment; do not treat the race suite as passing here.
