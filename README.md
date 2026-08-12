# Remember My Pill

## Current implementation

`frontend/` is a full Next.js marketing site built on an RMP design-token layer (Tailwind v4 `@theme`), translated from the official Figma brand identity: header/hero, problem framing, a Snap → Understand → Act story, an original tabbed product-preview (Today / Add medication / Reminder screens, illustrating the PRD's described mobile UI — not real screenshots, none exist in Figma), a caregiver section, benefit and privacy/trust cards, an honest "how the waitlist works" bridge, a waitlist CTA, an FAQ, and a footer. Motion uses Framer Motion with `prefers-reduced-motion` respected globally.

`/waitlist` is a dedicated, single-viewport, email-only signup page (`frontend/app/waitlist/page.tsx`); every high-intent CTA on the landing page routes there instead of embedding the form inline.

The waitlist API stores a required, normalized email and an optional name (Phase 2 + the email-only `/waitlist` page): it does not process prescriptions, collect medication data, provide medical guidance, or implement referrals, rank, or consent persistence. The site's copy is written to match this exactly — it does not promise referral links, waitlist rank, or a consent record the backend does not have.

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
