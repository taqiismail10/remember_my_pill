# Waitlist API

The current implementation is intentionally limited to waitlist registration
and must not accept medication, prescription, insurance, or other health data.

## Current implemented contract

`POST /api/waitlist` accepts required `email` and optional `name`. The server
trims both fields and stores email in lowercase. Omitted/blank `name` is stored
as SQL `NULL`. A newly accepted request and an existing-email request both
return exactly `202 {"status":"accepted"}`; the API never discloses which
case occurred. Invalid data returns `400 WAITLIST_VALIDATION_ERROR`; more than
five write attempts/IP/minute returns `429 WAITLIST_RATE_LIMITED`; storage
failures return `500 WAITLIST_INTERNAL_ERROR`. Invalid write attempts consume a
limit attempt. Unknown JSON fields and trailing JSON are rejected; the request
body limit is 4 KiB.

For the pilot draft, the waitlist UI submits `consent: true` with exact
`waitlist-consent-v1`, plus a separate `marketingConsent` choice. The server,
not the client, sets both consent timestamps. A `marketingConsent: true`
request requires exact `marketing-consent-v1`; a declined marketing choice
stores neither marketing timestamp nor version. Explicit false required
consent, missing paired consent fields, or a mismatched version is rejected.

`PILOT_LEGAL_CONTENT_APPROVED=false` keeps the historical email-only API
compatibility path available as a truthful pre-consent record. It cannot
record a marketing opt-in. When the server-only flag is `true`, all new
requests require `waitlist-consent-v1`. The `company` field is the documented
honeypot: a non-empty value receives the same 202 response without an insert.
`referralCode` is reserved and is not used or returned.

`GET /health` is liveness-only and returns exactly `{"status":"ok"}`. `GET
/ready` pings PostgreSQL and returns `200 {"status":"ready"}` only when the
dependency is usable; otherwise it returns a generic `503
{"status":"unavailable"}`. Neither endpoint discloses connection details.

Every response has an `X-Request-ID`. Incoming IDs and `Forwarded` or
`X-Forwarded-For` addresses are honored only when the direct peer belongs to
validated `TRUSTED_PROXY_CIDRS`; otherwise client identity is `RemoteAddr` and
the server generates the request ID. Browser access is allowlisted by
`ALLOWED_ORIGIN`; preflight permits `POST, OPTIONS` and `Content-Type`, with no
wildcard or credential mode.

## Locked target contract

Backend Phase B1 locked the later status, referral, token, and admin-export
contracts in [backend-phase-b1-contract.md](backend-phase-b1-contract.md).
B2A implements only the consent transition, duplicate privacy, readiness, and
proxy/request-ID foundations. Status, referral, token issuance, and admin APIs
remain unimplemented.

B3A adds no public endpoints. Its referral-event, verification-token, and
status-token primitives are internal only and remain disabled until B2B is
approved and later B3 activation work is authorized.

B3B adds a disabled transactional-email adapter only. The router does not
register `POST /api/waitlist/status-access/request`,
`POST /api/waitlist/status-access/exchange`, `GET /api/waitlist/status`, a
referral-resolve route, or a Postmark webhook endpoint. See
[B3B email infrastructure](backend-phase-b3b-email-infrastructure.md).

B3C adds internal, server-side browser status-session primitives only. It does
not register a verification exchange, logout, status, referral, or webhook
route. See [B3C browser session infrastructure](backend-phase-b3c-status-session-infrastructure.md).

B3D-Prep registers the future request, exchange, status, and logout routes,
but every one returns feature-unavailable while `STATUS_ACCESS_ENABLED=false`
or pilot legal approval is false. They are not public functionality. See
[B3D status-access foundation](backend-phase-b3d-status-access-foundation.md).

## B2B-Prep pilot consent foundation

The USA + Canada adult-pilot foundation recognizes `waitlist-consent-v1` and
separate optional `marketing-consent-v1`. Required consent evidence uses the
existing `consent_version` and server-generated `consented_at`; marketing
evidence is independent and only recorded when `marketingConsent: true` is
paired with its exact version. `marketingConsent: false` creates no marketing
timestamp or version. Public responses remain generic `202` for new,
duplicate, and honeypot requests.

`PILOT_LEGAL_CONTENT_APPROVED` defaults to `false`. While false, legacy
email-only compatibility remains in place; when true, it requires the exact
waitlist consent version. This technical flag does not approve legal content,
enable B3 routes, referrals, Postmark delivery, or any public pilot launch.
