# Backend Phase B1 — contract and schema lock

**Status:** Approved design only. This document does not implement application
or database behavior.

This is the canonical B1 record for the PRD-aligned backend contract. The
current runtime contract remains separately documented in [API](api.md).

## Compatibility and privacy policy

- Email-only signup is the approved primary UX. `name` remains optional, and
  legacy name-plus-email clients remain supported.
- New post-B2 signups require `consent: true` and `consentVersion` equal to a
  server-configured active consent-policy identifier. The server generates
  `consentedAt`.
- The current email-only request remains accepted during the B2 migration
  window. B2 must document the cutover and then reject missing consent for new
  signups.
- Existing rows are pre-consent records. They must not be backfilled with
  invented consent values and must remain distinguishable from consented rows.
- Target duplicate handling is success-like and non-enumerating. A duplicate
  must not reveal that an email exists or return its existing token, code,
  count, rank, or personal data.

## Target public API

### `POST /api/waitlist`

```json
{
  "name": "Alice",
  "email": "alice@example.com",
  "consent": true,
  "consentVersion": "LEGAL_POLICY_VERSION",
  "referralCode": "rmp-example-code",
  "company": ""
}
```

| Field | Target rule |
| --- | --- |
| `email` | Required; trim, lowercase, validate, maximum 255 characters. |
| `name` | Optional; trim; blank becomes SQL `NULL`; maximum 100 characters. |
| `consent` | Required and `true` after B2 cutover. |
| `consentVersion` | Required after B2 cutover; exact active server-configured value. |
| `referralCode` | Optional; invalid code does not block a valid signup. |
| `company` | Optional honeypot; a non-empty value creates no entry and discloses nothing. |

Unknown JSON fields, trailing JSON, and oversized bodies remain rejected. The
target body limit remains 4 KiB unless a later approved contract changes it.

Uniform post-B2 public response, for both a newly created entry and a
duplicate:

```json
{
  "status": "accepted"
}
```

This response uses the same `202 Accepted` status for both cases. It does not
return a status token, referral code, rank, referral count, or personal data.
Returning a fresh status token only to a new entry would itself be a duplicate
enumeration side channel. A privacy-preserving status-token delivery mechanism
(for example, an approved verified out-of-band channel) is therefore required
before B3 exposes status lookup; it is an owner/product decision, not an
assumption B1 can invent.

### `GET /api/waitlist/status`

```text
Authorization: Bearer <opaque-status-token>
```

No email, raw token query parameter, or guessable identifier is accepted.

```json
{
  "rank": 213,
  "referralCount": 2,
  "referralCode": "rmp-example-code",
  "referralUrl": "https://remembermypill.com/?ref=rmp-example-code"
}
```

The response excludes email, name, database IDs, consent metadata, referrer
identity, and referred-user identity.

### `POST /api/waitlist/referral/resolve`

Public and rate-limited. It may disclose referral validity only, for example
`{"status":"valid"}`. It must not return referrer personal data or counts.

### `GET /api/admin/export`

The initial design requires a strong environment-provided credential in a
request header. Query-string credentials are prohibited. Production also needs
reverse-proxy/network restriction, while preserving future SSO or
identity-aware-proxy replacement. It has no browser/public CORS. Successful
and failed authorization attempts are auditable without logging credentials.

Permitted CSV fields: `created_at`, `email_normalized`, `name`,
`consent_version`, `consented_at`, `referral_code`, derived `referral_count`,
and derived `rank`. Exclude token hashes/tokens, IP addresses, user agents,
database IDs, and referral-event relationships.

### Health contracts

- `GET /health` remains liveness-only and minimal.
- Future `GET /ready` returns success only when required dependencies,
  including PostgreSQL, are usable; it returns a generic non-200 response
  otherwise and never leaks database/configuration details.

## Stable public errors

```json
{
  "error": {
    "code": "WAITLIST_VALIDATION_ERROR",
    "message": "Enter a valid email."
  }
}
```

Safe categories are validation, rate-limit, unauthorized admin, and generic
internal errors. Never expose SQL errors, connection strings, stack traces,
internal paths, credentials, token values, or unrelated-email existence.

## Referral, rank, and token rules

- Referral codes are random, non-derivable, unique, and retried on a database
  uniqueness collision.
- One referred entry creates at most one credit; self-referral creates none.
- Attribution/event creation is transactionally and concurrently safe.
- Referral count is initially derived from events, not a mutable counter.
- Rank v1 is referral-independent:

```text
rank = 1 + count(entries ordered ahead by created_at ASC, id ASC)
```

Future reward thresholds require product approval and server configuration,
never frontend hard-coding.

Status tokens are cryptographically random with high entropy. Store only a
hash, return the raw value only through the approved privacy-preserving
issuance mechanism, and never log or place it in a query string.

## Request/proxy and scaling policy

The initial deployment is one API instance. The current in-memory limit remains
five attempts/IP/minute. Use `RemoteAddr` by default. `Forwarded` and
`X-Forwarded-For` are trusted only behind explicitly configured trusted proxy
infrastructure; public client headers are never trusted blindly. Horizontal
scaling requires a shared/distributed limiter before multiple API instances
serve traffic.

## Forward-only schema design

Do not alter migrations `001` or `002`.

`003_waitlist_referral_consent` will add `referral_code`, `referred_by_id`,
`status_token_hash`, `consent_version`, `consented_at`, and `updated_at` to
`waitlist_entries`. Consent fields are initially nullable for truthfulness of
legacy pre-consent rows. B2 application writes require consent; a later
constraint migration is allowed only after approved legacy-entry disposition.

`004_referral_events` will create:

| Item | Design |
| --- | --- |
| `id` | UUID primary key. |
| `referrer_id` | Required FK to `waitlist_entries`. |
| `referred_entry_id` | Required FK to `waitlist_entries`. |
| `created_at` | Server-generated timestamp. |
| `UNIQUE(referred_entry_id)` | One credit per referred signup. |
| `CHECK(referrer_id <> referred_entry_id)` | Database self-referral guard. |
| Index on `referrer_id` | Efficient derived counts. |

Migration tests must run ordered `001 → 002 → 003 → 004` upgrades and paired
down migrations against disposable PostgreSQL.
