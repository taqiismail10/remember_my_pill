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

B2A keeps the email-only client compatible: a payload containing only `email`
(and optional `name`) is accepted as a truthful pre-consent record. A future
client may instead send `consent: true` and `consentVersion`; the version must
exactly match non-empty `CONSENT_VERSION`, and the server—not the client—sets
the consent timestamp. Explicit `false`, missing paired consent fields, or a
mismatched version is rejected. The `company` field is the documented
honeypot: a non-empty value receives the same 202 response without an insert.
`referralCode` is accepted only as a reserved future field and is not used or
returned in B2A.

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
