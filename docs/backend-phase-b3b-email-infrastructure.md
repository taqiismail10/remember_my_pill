# Backend Phase B3B — email infrastructure

**Status:** Infrastructure only. No public status-access route, webhook route,
or production email delivery is active.

Postmark is the initial transactional-email provider. The `access.EmailSender`
boundary remains provider-neutral, so Amazon SES can be introduced later
without changing waitlist or access-domain callers. `EMAIL_PROVIDER=fake` is
the default and makes no network request; automated tests use it or local
mocked HTTP servers only.

## Configuration and secrets

`EMAIL_PROVIDER` accepts `fake` or `postmark`. Enabling `postmark` requires a
valid sender address/name, HTTPS `STATUS_ACCESS_BASE_URL`,
`POSTMARK_SERVER_TOKEN`, and `POSTMARK_MESSAGE_STREAM`. Startup rejects an
invalid explicit Postmark configuration without disclosing the bad value.
`POSTMARK_SERVER_TOKEN` is a secret and must be supplied only through the
runtime environment or deployment secret store. No secret belongs in Git.

`POSTMARK_STATUS_ACCESS_TEMPLATE` remains optional while final legal/product
copy is unapproved. The adapter rejects a send with no template alias before
making any HTTP request. A later approved template receives only the
short-lived verification URL as its template model; a long-lived status token
is never placed in email or a URL.

For an explicitly requested manual smoke test, Postmark's `POSTMARK_API_TEST`
server token can exercise the API without delivery. This is never used by
normal tests or application startup.

## Delivery, retry, and suppression

Postmark responses normalize internally to `accepted`, `temporary_failure`,
`rejected`, or `suppressed`. A Postmark inactive-recipient response (error
406) becomes `suppressed`; HTTP 429, 5xx, malformed responses, timeouts, and
network failures become `temporary_failure`. Future delivery work may retry a
temporary failure with a small bounded policy (for example, three attempts
with exponential backoff); it must not retry permanent rejection or
suppression, and it must not create an unbounded worker queue.

No delivery events or recipient data are persisted in B3B. When status access
is activated, suppression must affect delivery internally while its public
request response remains generic, avoiding an email-enumeration oracle.

## Webhook decision

No webhook receiver is registered in this phase. Postmark's current webhook
documentation states that it does not provide HMAC signature verification.
A future receiver therefore requires HTTPS, reverse-proxy Basic Authentication
embedded in the configured callback URL, Postmark IP allowlisting, strict
event-shape validation, and idempotent `MessageID` handling before it can be
mounted. It may retain only a correlation identifier, event category, and
timestamp; it must not retain payloads, message contents, URLs, tokens, IP,
or user-agent data. Webhook events remain operational signals only and cannot
change consent, eligibility, or referral attribution.

## Activation dependency

B2B consent approval and the approved legacy-entry policy remain prerequisites
for any status-access request or delivery flow. Pre-consent B2A records must
not receive status-access email. The production router still exposes no status
access, referral resolve, or Postmark webhook endpoint.

Postmark behavior is based on its official [Email API](https://postmarkapp.com/developer/api/email-api), [API sending guide](https://postmarkapp.com/developer/user-guide/send-email-with-api), and [webhook overview](https://postmarkapp.com/developer/webhooks/webhooks-overview).
