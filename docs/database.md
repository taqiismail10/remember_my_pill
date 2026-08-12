# Waitlist database

The repository-owned development PostgreSQL 16 service is bound only to
`127.0.0.1:5433`. Current schema is supplied by
`backend/migrations/001_waitlist_entries.up.sql`,
`002_waitlist_entries_name_optional.up.sql`,
`003_waitlist_referral_consent.up.sql`,
`004_referral_events_access_tokens.up.sql`, and
`005_waitlist_marketing_consent.up.sql`. `waitlist_entries` has a PostgreSQL
UUID, bounded optional nullable name, bounded unique normalized email, and
server-generated timestamps. Migration 003 adds nullable staged
`referral_code`, `referred_by_id`, `status_token_hash`, `consent_version`, and
`consented_at` fields, plus `updated_at`. It has no health or prescription data.

Migration 003 is forward-only and deliberately leaves all newly introduced
consent, referral, and status columns NULL for existing records. No consent,
referral, or token data is fabricated or backfilled. The B2B-Prep required
consent flow writes only a valid configured version and a server-generated
timestamp; it does not generate or expose referral codes or status tokens.

Migration 004 now provides B3A-only internal foundations: constrained
`referral_events` and hash-only `waitlist_verification_tokens`. It does not
enable referral attribution, status access, public B3 endpoints, real email,
or generate values for historical rows.

Migration 005 adds independent nullable marketing-consent evidence:
`marketing_consent_version`, `marketing_consented_at`, and a future
`marketing_withdrawn_at`. A check constraint prohibits a consent version
without its timestamp (or the reverse). Existing and historical rows remain
NULL/marketing-unconsented; no consent is backfilled or inferred.

Do not run down migrations against the development database. The disposable
integration database is separately bound to `127.0.0.1:5434` through
`docker-compose.test.yml` and must be torn down only with its own Compose
project and volume.
