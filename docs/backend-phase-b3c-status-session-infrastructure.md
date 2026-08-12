# Backend Phase B3C — browser status-session infrastructure

## Scope and activation state

B3C implements internal persistence and service primitives only. It does not
register public B3 routes, issue a cookie from an HTTP handler, send Postmark
email, change referral behavior, or alter B2B privacy/consent behavior. Those
steps remain separately gated by owner/legal approval.

## Session model

Migration 006 creates `status_access_sessions`. A session contains an opaque
256-bit random value generated with the existing access-token primitive; only
its SHA-256 hash is stored in PostgreSQL. The value is returned only from the
trusted internal creation/exchange boundary and must never be logged,
persisted in browser JavaScript storage, placed in a URL, or sent to analytics.

Sessions have a fixed seven-day lifetime. Lookup accepts only unrevoked,
unexpired hashes and never extends expiry. Current-session revocation,
all-session revocation, and expired-row cleanup are available for later HTTP
handlers. A future successful verification exchange will consume the
short-lived verification token and create the browser session in one database
transaction.

At most five active sessions may exist for an entry. Creation locks the entry
row, then deterministically revokes the oldest active session (`created_at`,
then `id`) before a sixth is inserted. A new verification does not revoke the
other active sessions unless the cap requires that eviction.

## Future cookie contract

The future same-origin production deployment uses `https://remembermypill.com/api/*`.
When public B3 handlers are approved, they must use:

| Attribute | Value |
| --- | --- |
| Name | `__Host-rmp-status` |
| Value | Opaque server-side session credential |
| HttpOnly | `true` |
| Secure | `true` |
| SameSite | `Lax` |
| Path | `/` |
| Domain | omitted (host-only) |
| Max-Age / expiry | seven days; fixed, not sliding |

The `__Host-` prefix requires `Secure` and no `Domain`, preventing the cookie
from being scoped to a parent or sibling subdomain. Session helpers provide a
distinct `rmp-status-dev` name only for explicit local HTTP development. This
is intentional: browsers reject an insecure `__Host-` cookie.

`RMP_ENVIRONMENT` defaults to `production`; `STATUS_SESSION_COOKIE_SECURE`
defaults to `true`. `false` is accepted only when the environment is exactly
`development`; startup rejects any insecure production configuration. Cookie
clearing must use the same name, path, security, HttpOnly, and SameSite scope.

## Future exchange and CSRF boundary

The approved future exchange is: email link with a short-lived one-time
verification value → frontend verification page → POST body to the API →
atomic consume/create → Set-Cookie → `history.replaceState`/redirect to the
status page. The verification value must be removed before analytics, external
navigation, or a referrer-bearing request. The long-lived session credential
never appears in the URL.

SameSite=Lax provides baseline CSRF protection for ordinary cross-site POSTs.
Future state-changing cookie-authenticated B3 endpoints must also enforce an
Origin check and a CSRF token or double-submit mechanism as appropriate. The
initial same-origin topology avoids credentialed cross-origin CORS. If an API
subdomain or separate frontend origin is chosen later, the owner must approve
the changed cookie/CORS/CSRF design before routes are enabled.

Native clients can later use their own securely stored bearer credential behind
the same status service; browser storage does not need to weaken to match that
future client model.
