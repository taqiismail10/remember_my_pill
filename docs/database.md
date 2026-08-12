# Waitlist database

The repository-owned development PostgreSQL 16 service is bound only to
`127.0.0.1:5433`. Current schema is supplied by
`backend/migrations/001_waitlist_entries.up.sql`,
`002_waitlist_entries_name_optional.up.sql`, and
`003_waitlist_referral_consent.up.sql`. `waitlist_entries` has a PostgreSQL
UUID, bounded optional nullable name, bounded unique normalized email, and
server-generated timestamps. Migration 003 adds nullable staged
`referral_code`, `referred_by_id`, `status_token_hash`, `consent_version`, and
`consented_at` fields, plus `updated_at`. It has no health or prescription data.

Migration 003 is forward-only and deliberately leaves all newly introduced
consent, referral, and status columns NULL for existing records. No consent,
referral, or token data is fabricated or backfilled. B2A only writes consent
columns when a valid future-style consent payload is supplied; it does not yet
generate or expose referral codes or status tokens. Migration 004 remains
reserved for the later referral-event work.

Do not run down migrations against the development database. The disposable
integration database is separately bound to `127.0.0.1:5434` through
`docker-compose.test.yml` and must be torn down only with its own Compose
project and volume.
