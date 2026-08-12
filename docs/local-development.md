# Local development

Use the existing Compose files without changing unrelated Docker resources. Development PostgreSQL is `rmp-phase2-dev-db` on localhost:5433. Copy only placeholder values from `.env.example` into the process environment; do not place production credentials in this repository. Start the backend from `backend/` with `go run ./cmd/server` after setting `DATABASE_URL`, `PORT`, and a valid absolute `ALLOWED_ORIGIN` such as `http://localhost:3000`. B2A also recognizes optional `CONSENT_VERSION` and `TRUSTED_PROXY_CIDRS`. Leave `CONSENT_VERSION` blank until an approved legal policy identifier exists: email-only compatibility remains active, while consent-bearing payloads are rejected. Set trusted proxy CIDRs only for the direct reverse proxies that terminate connections to this service; a blank value means `RemoteAddr` is always used.

The separate test Compose project is `rmp-phase2-test` and only uses localhost:5434. It must never use port 5432 or the development service's port 5433.

## Locked deployment assumptions for later backend phases

The initial public deployment is one API instance, so the current in-memory
five-attempts/IP/minute limiter remains acceptable until horizontal scaling is
introduced. The server uses `RemoteAddr` by default. Proxy headers such as
`Forwarded` and `X-Forwarded-For` are trusted only when the direct peer is in
configured `TRUSTED_PROXY_CIDRS`; they must never be trusted blindly from
public clients. A distributed/shared limiter is required before multiple API
instances serve traffic.

The future admin credential is environment-provided and header-based; never
place credentials in a query string or local documentation. See
[backend-phase-b1-contract.md](backend-phase-b1-contract.md) for the locked
contract and legacy-consent policy.

B3A has no email provider configuration and sends no email. Its internal token
and referral foundations remain inactive while B2B consent approval is pending.

B3B adds a disabled provider foundation. Keep `EMAIL_PROVIDER=fake` for normal
local development and tests; it needs no credentials and never contacts a
provider. `EMAIL_PROVIDER=postmark` is an explicit opt-in and requires the
validated placeholders in `.env.example`, including a secret
`POSTMARK_SERVER_TOKEN`. It does not activate routes or send email by itself.
The final status-access template and B2B/legacy-consent approval are still
required before delivery can be enabled.

The B2B-Prep privacy and consent foundation is a USA + Canada, adults-18+
pilot draft. `PILOT_LEGAL_CONTENT_APPROVED=false` is the required safe default.
Changing it to `true` only enforces the configured required consent version;
it does not establish legal approval or activate email, B3, referrals, or the
pilot. The legal operator, mailing address, final review, and Terms of Use
remain launch blockers.
