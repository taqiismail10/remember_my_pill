# Waitlist database

The repository-owned development PostgreSQL 16 service is bound only to
`127.0.0.1:5433`. Current schema is supplied by
`backend/migrations/001_waitlist_entries.up.sql` and
`002_waitlist_entries_name_optional.up.sql`. `waitlist_entries` has a
PostgreSQL UUID, bounded optional nullable name, bounded unique normalized
email, and server-generated timestamp. It has no health or prescription data.

Backend Phase B1 locks the forward-only design for consent, referral, token,
rank, and legacy-entry handling in
[backend-phase-b1-contract.md](backend-phase-b1-contract.md). Migrations `003`
and `004` are intentionally not created in B1. Existing rows must remain
truthfully pre-consent; no consent values may be fabricated or backfilled.

Do not run down migrations against the development database. The disposable
integration database is separately bound to `127.0.0.1:5434` through
`docker-compose.test.yml` and must be torn down only with its own Compose
project and volume.
