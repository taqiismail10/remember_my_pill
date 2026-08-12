# Backend Phase B2B-Prep — USA + Canada pilot privacy and consent foundation

**Status:** Technical and pilot-draft content foundation only; not approved for
public pilot or production launch.

## Scope and legal blockers

The intended pilot is for adults aged 18+ in the United States and Canada.
Remember My Pill is being prepared for medication-information organization,
extraction, user confirmation or correction, reminder setup, and reminder
notifications. It does not diagnose, prescribe, recommend medication changes,
or replace a doctor or pharmacist.

The legal operator, business mailing address, final external legal review, and
approved Terms of Use are unresolved. `/terms` is deliberately not created:
the supplied material did not include approved Terms, so inventing them would
be unsafe. These are production blockers. The detailed owner privacy/consent
specification referenced for this phase was not available in the repository or
attachment set; the draft routes therefore state only supplied scope and
verified implementation facts, without inventing legal promises.

## Consent model

Required `waitlist-consent-v1` is recorded in `consent_version` and
server-generated `consented_at`. Optional `marketing-consent-v1` is recorded
independently only when affirmatively selected, with its own version and
server-generated timestamp. A marketing opt-in without required waitlist
consent is rejected. Declining marketing consent stores neither version nor
timestamp. Future withdrawal may set `marketing_withdrawn_at` without
rewriting historical evidence.

Legacy/pre-consent records remain NULL for both consent types. They receive no
status-access or promotional email; an owner-approved re-consent policy is
still required.

## Activation control

`PILOT_LEGAL_CONTENT_APPROVED=false` is the safe default. When false, the
existing server compatibility path can still accept historical email-only
clients as truthful pre-consent records. When true, new requests require exact
`waitlist-consent-v1`. The flag is server-side only. It does **not** approve
legal content, launch the pilot, add Terms, enable B3 status routes, referral
flows, Postmark delivery, or webhooks.

`/privacy` and `/consumer-health-data-privacy` visibly identify themselves as
pilot drafts and not production-ready. They do not fabricate a legal entity,
mailing address, privacy contact, retention policy, or health-data advertising
promise.
