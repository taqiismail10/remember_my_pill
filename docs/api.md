# Waitlist API

The Phase 2 API is intentionally limited to waitlist registration. It must not be used to submit medication, prescription, insurance, or other health data.

`POST /api/waitlist` accepts a JSON object with a required `email` and an optional `name`. The server trims both fields and stores email in lowercase. `name` may be omitted entirely (e.g. the email-only `/waitlist` page) or sent alongside `email` for backward compatibility with the original name+email contract; when omitted or blank, `name` is stored as SQL `NULL` rather than an empty string. A successful response is `201` with `{"status":"created"}`. Invalid JSON or values return `400` and `WAITLIST_VALIDATION_ERROR`; a normalized duplicate returns `409` and `WAITLIST_EMAIL_EXISTS`; more than five write attempts from one IP in a minute returns `429` and `WAITLIST_RATE_LIMITED`; storage failures return `500` and `WAITLIST_INTERNAL_ERROR`. Invalid write attempts consume a rate-limit attempt because the limiter deliberately protects the entire write endpoint before parsing. Unknown JSON fields are still rejected.

`GET /health` is a liveness endpoint only and returns exactly `{"status":"ok"}`. It does not disclose database state or configuration.

Browser access is allowlisted by the required `ALLOWED_ORIGIN` configuration. Only that exact origin receives CORS permission; preflight permits `POST, OPTIONS` and `Content-Type`, with no wildcard or credentials mode.
