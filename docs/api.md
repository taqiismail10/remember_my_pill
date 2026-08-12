# Waitlist API

The current implementation is intentionally limited to waitlist registration
and must not accept medication, prescription, insurance, or other health data.

## Current implemented contract

`POST /api/waitlist` accepts required `email` and optional `name`. The server
trims both fields and stores email in lowercase. Omitted/blank `name` is stored
as SQL `NULL`. A success is `201 {"status":"created"}`. Invalid data returns
`400 WAITLIST_VALIDATION_ERROR`; a duplicate returns
`409 WAITLIST_EMAIL_EXISTS`; more than five write attempts/IP/minute returns
`429 WAITLIST_RATE_LIMITED`; storage failures return
`500 WAITLIST_INTERNAL_ERROR`. Invalid write attempts consume a limit attempt.
Unknown JSON fields and trailing JSON are rejected; the request body limit is
4 KiB.

`GET /health` is liveness-only and returns exactly `{"status":"ok"}`. It does
not disclose database state or configuration. Browser access is allowlisted by
`ALLOWED_ORIGIN`; preflight permits `POST, OPTIONS` and `Content-Type`, with
no wildcard or credential mode.

## Locked target contract

Backend Phase B1 has locked, but not implemented, the consent, duplicate
privacy, status, referral, admin-export, readiness, proxy, and token contracts
in [backend-phase-b1-contract.md](backend-phase-b1-contract.md). New post-B2
signups require explicit consent; the email-only request remains compatible
only during the documented migration window; and duplicate handling will become
success-like/non-enumerating.
