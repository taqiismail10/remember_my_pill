# Waitlist database

The repository-owned development PostgreSQL 16 service is bound only to `127.0.0.1:5433`. It uses `backend/migrations/001_waitlist_entries.up.sql`, `002_waitlist_entries_name_optional.up.sql`, and their paired down migrations. The sole table, `waitlist_entries`, has a PostgreSQL-generated UUID, a bounded optional (nullable) name, bounded unique normalized email, and PostgreSQL-generated timestamp. `name` is nullable so the email-only `/waitlist` page can sign up without collecting a name; when provided it is still bounded to 1–100 characters by the existing check constraint. It has no health, prescription, referral, ranking, or administrative fields.

Do not run the down migration against the development database. The intended disposable integration database is separately named and bound only to 127.0.0.1:5434 through `docker-compose.test.yml`; it must be torn down with its own Compose project and volume only.
