# Remember My Pill

## Current implementation

`frontend/` is a full Next.js marketing site built on an RMP design-token layer (Tailwind v4 `@theme`), translated from the official Figma brand identity: header/hero, problem framing, a Snap → Understand → Act story, an original tabbed product-preview (Today / Add medication / Reminder screens, illustrating the PRD's described mobile UI — not real screenshots, none exist in Figma), a caregiver section, benefit and privacy/trust cards, an honest "how the waitlist works" bridge, a waitlist CTA, an FAQ, and a footer. Motion uses Framer Motion with `prefers-reduced-motion` respected globally.

`/waitlist` is a dedicated, responsive email signup page (`frontend/app/waitlist/page.tsx`); every high-intent CTA on the landing page routes there instead of embedding the form inline. It presents an unchecked required waitlist consent and a separate optional marketing consent for the USA + Canada adult-pilot draft. It does not collect medication, prescription, or other health information.

The waitlist API stores a normalized email, optional name, and server-generated consent evidence. Required `waitlist-consent-v1` and optional `marketing-consent-v1` are persisted independently; historical rows are not backfilled. The legal-content activation flag defaults off, and the pilot is not publicly approved or launched. Status access, referrals, rank, Postmark delivery, and webhooks remain inactive.

Pilot-draft legal pages are available at [/privacy](frontend/app/privacy/page.tsx) and [/consumer-health-data-privacy](frontend/app/consumer-health-data-privacy/page.tsx). They are not final legal notices. A legal operator, business mailing address, final external review, and approved Terms of Use remain required before launch; `/terms` is intentionally absent rather than containing invented terms. See [the B2B-Prep record](docs/backend-phase-b2b-prep.md).

Product requirements are maintained in [the RMP PRD](docs/product/rmp-prd.md); supplementary design references are in [docs/design/references](docs/design/references/).

## Local development

```powershell
cd frontend
npm install
npm run dev
```

Available checks: `npm run format:check`, `npm run lint`, `npm run typecheck`, `npm test`, `npm run test:e2e`, and `npm run build`.

See [Phase 1 plan](docs/phase-1-implementation-plan.md) and the [brand export manifest](docs/brand-export-manifest.md) for scope and asset provenance.

## Backend local development

Start the repository-owned development database on port 5433, copy the placeholder values in `.env.example` into your local shell, then run the Go service from `backend/`. `ALLOWED_ORIGIN` is required and must be a complete `http` or `https` origin with no path. The API and testing details are in [local development](docs/local-development.md), [API](docs/api.md), [database](docs/database.md), and [backend testing](docs/backend-testing.md).
