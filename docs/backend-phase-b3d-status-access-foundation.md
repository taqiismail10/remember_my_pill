# Backend Phase B3D-Prep — status-access foundation

## Status: implemented but disabled pending legal/product activation

`STATUS_ACCESS_ENABLED=false` is the default server gate. The feature becomes
internally active only when that flag and `PILOT_LEGAL_CONTENT_APPROVED=true`
are both set, with a valid HTTPS `STATUS_ACCESS_BASE_URL` and a base64-encoded
32-byte-or-longer `STATUS_ACCESS_RATE_LIMIT_KEY`. The frontend has its own
pre-launch display gate, but it cannot enable the server.

While inactive, `POST /api/waitlist/status-access/request`,
`POST /api/waitlist/status-access/exchange`, `GET /api/waitlist/status`, and
`POST /api/waitlist/status/logout` return feature-unavailable and set no
status cookie. Postmark remains inactive in normal configuration; automated
tests use only the fake sender.

## Future active contracts

The request endpoint accepts `{ "email": "..." }` and always returns
`202 {"status":"accepted"}` for valid, unknown, legacy, provider-failure,
suppressed, duplicate, and rate-limited cases. Eligibility is internal:
entries must have explicit required-consent evidence; legacy NULL-consent
records cannot obtain a verification token or session. Request protection is a
bounded single-instance in-memory limiter: three keyed-HMAC IP requests per
hour and three keyed-HMAC normalized-email requests per day. A distributed
limiter is required before horizontal scaling. Raw normalized email is never a
rate-limit map key or a log field.

Verification links use `https://remembermypill.com/waitlist/verify#v=...`.
The fragment is not sent with the initial request. The browser route uses
`Referrer-Policy: no-referrer`, clears history immediately, then POSTs
`{ "verificationToken": "..." }` to the exchange route. No localStorage,
sessionStorage, IndexedDB, URL query, or JS-readable cookie carries the
long-lived authorization.

The active exchange atomically consumes the one-time 15-minute token and
creates the B3C hash-only seven-day server session. It sets only the matching
`__Host-rmp-status` production cookie (Secure, HttpOnly, SameSite=Lax,
Path=/, no Domain). Failed exchange returns generic authorization failure and
sets no cookie.

The active status route returns only rank, independent referral count, and an
optional referral code/validated canonical `/waitlist?ref=<code>` URL; it
never returns email, names, consent, IDs, or token material. Logout requires
exact configured Origin, `Sec-Fetch-Site` same-origin/same-site, and JSON
content type; it revokes only the current session, clears its cookie, and is
idempotent. Future sensitive cookie-authenticated mutations require a fuller
CSRF synchronizer/double-submit design.

Preferred production topology remains same-origin:
`https://remembermypill.com/api/*` reverse-proxied to Go. Separate local
origins are configuration-only development support, not a production CORS
architecture.
