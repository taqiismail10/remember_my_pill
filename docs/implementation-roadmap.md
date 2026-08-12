# Implementation roadmap

## Current-status update — 16 July 2026

Phase 1 — frontend foundation and accessible marketing shell is verified complete; see [the completion report](phase-1-completion-report.md). The approved Phase 2 boundary now includes the minimal Go/PostgreSQL waitlist vertical slice so a real end-to-end flow can be verified. Phase 3 is reserved for referrals, ranking, status, sharing, confirmation email, administration, and operational hardening. Phase 2 and later remain not started. The original roadmap below records the pre-implementation plan.

**Roadmap date:** 16 July 2026
**Baseline:** Remember My Pill PRD v1.2.0
**Current state:** Not started; the repository contains requirements, references, and planning documents but no runnable product implementation.

This roadmap is an evidence-based sequence. A phase is not complete until its exit criteria are demonstrated in the repository and, where applicable, in browser or service verification.

## Guiding decisions

- The PRD is the product and functional source of truth.
- The official Figma Brand Identity file is the visual identity source of truth.
- Local PNG references support interpretation but do not override Figma.
- The waitlist release must not collect medication or other health information.
- Product showcases must be labelled as previews and must not imply live public prescription processing.
- No production secrets, real patient data, or unsupported healthcare claims may enter the repository.
- Preserve the PRD-proposed separation between frontend, backend, and PostgreSQL unless a documented architecture decision changes it.

## Current baseline

### Repository evidence

- Tracked files: .gitignore and README.md only.
- Working-tree materials: AGENTS.md, docs/, and the ignored PRD/reference bundle.
- Missing: frontend/, backend/, package manifests, Go module, migrations, Docker/Compose, CI, environment template, assets pipeline, and tests.
- README.md contains only the project title.
- Latest commit: 0d53b11 Remove PRD assets from repository.

### Design evidence

Figma MCP access is working for the supplied official file. Inspected frames include:

- RMP Pallete 1 (6026:3);
- Youtube Banner (6021:204);
- Facebook Cover (6023:220);
- 2D Mascot (6023:243);
- 3D Mascot (6023:241);
- 3D Logo Render (6023:242).

The Figma identity board confirms Manrope, the warm cream/terracotta/orange palette, official R/pill logo variants, mascot artwork, supporting icons, and the warm/reassuring/simple/trustworthy brand voice. The Figma design-system search returned no published variables, components, or styles, so the repository must receive an explicit export and token-mapping record.

## Dependency flow

Figma asset export and provenance
            |
            v
brand tokens, typography, logo, mascot, and shared primitives
            |
            v
responsive marketing shell and product showcase
            |
            +-----------------------------+
            v                             v
accessible waitlist UI             Go API and PostgreSQL
            |                             |
            +-------------+---------------+
                          v
             referral/status experience
                          |
                          v
             QA, deployment, and documentation

## Phase 0 — source control, brand export, and scaffold decisions

**Status:** Not started
**Purpose:** Remove ambiguity before implementation work begins.

### Work

- Decide whether the PRD bundle and current documentation should be tracked; do not change .gitignore without that decision.
- Confirm the frontend/backend service boundaries from the PRD.
- Select the frontend package manager and Go toolchain version.
- Export the approved horizontal/stacked/monogram logo variants, mascot artwork, icons, illustrations, product screens, and social-preview assets from Figma.
- Create the PRD-required brand export manifest containing Figma node/component name, export date, variant, intended use, background compatibility, and source URL.
- Map Figma typography and colors into a documented repository token plan.
- Confirm that no export contains real patient data or unsupported medical claims.

### Exit criteria

- The intended repository structure and source-control treatment are documented.
- Approved Figma assets are available for implementation with provenance.
- Token mapping distinguishes marketing identity colors from product UI and state roles.
- No visual implementation uses redrawn or placeholder logo/mascot artwork.

## Phase 1 — frontend foundation and accessible marketing shell

**Status:** Not started
**Purpose:** Establish the reusable visual and semantic foundation.

### Work

- Create the Next.js App Router and TypeScript frontend specified by the PRD.
- Add the selected styling/tooling configuration and centrally defined design tokens.
- Implement reusable logo, button, card, navigation, section, dialog, form-field, status, and icon primitives.
- Add semantic page structure, keyboard navigation, visible focus states, contrast-safe text, touch targets, and reduced-motion behavior.
- Build the header, hero, problem, solution, Snap–Understand–Act, and initial product-preview sections.
- Add responsive layouts for the PRD-required widths.

### Exit criteria

- The site renders at 320, 375, 390, 430, 768, 1024, 1280, and 1440+ px without horizontal overflow.
- Official Figma assets and mapped tokens are used.
- Keyboard focus and reduced-motion behavior are observable.
- No section claims that unavailable waitlist-release functionality is live.

## Phase 2 — product showcase, trust, and conversion UI

**Status:** Not started
**Purpose:** Complete the public-facing marketing experience and accessible waitlist surface.

### Work

- Add approved phone mock-ups and product-screen gallery components.
- Build caregiver support, privacy/trust, plans-preview, FAQ, and footer sections.
- Keep plans and future capabilities explicitly subject to change.
- Implement the waitlist form with name, email, consent, optional referral-code state, and a hidden honeypot.
- Implement initial, validation, submitting, success, existing-email, invalid-referral, rate-limited, network-failure, and copied-link states.
- Add accessible labels, associated errors, status announcements, and dialog focus return/trapping where a modal is used.
- Ensure the form never requests medication, diagnosis, prescription, insurance, or other health information.

### Exit criteria

- The full conversion flow is usable by keyboard and touch.
- All required form states have deterministic UI and non-technical feedback.
- Product previews are visually coherent with the Figma identity and local references.
- Responsive and accessibility checks pass before backend integration is treated as complete.

## Phase 3 — waitlist/referral API and PostgreSQL

**Status:** Not started
**Purpose:** Implement the privacy-preserving service behind conversion and referral behavior.

### Work

- Create the Go API and PostgreSQL connection layer.
- Add migrations for waitlist_entries and referral_events.
- Implement:
  - POST /api/waitlist;
  - GET /api/waitlist/status;
  - POST /api/waitlist/referral/resolve;
  - GET /api/admin/export;
  - GET /health.
- Normalize and validate names/emails; store consent version and timestamp.
- Generate unique referral codes and opaque status tokens; store only the status-token hash.
- Implement safe duplicate handling without email enumeration.
- Implement valid-referral credit exactly once, self-referral prevention, deterministic rank, and configurable future thresholds.
- Add IP rate limiting, request-size limits, honeypot handling, restricted CORS, secure headers, request IDs, and minimal security logs.
- Protect admin export with stronger authentication than a URL token.

### Exit criteria

- New entries, duplicates, invalid referrals, valid referrals, self-referrals, and rate limits are covered by automated tests.
- Database constraints prevent duplicate referral credit and duplicate normalized emails.
- No health-related fields are accepted by the waitlist API.
- The status endpoint uses an opaque token rather than a raw email query.

## Phase 4 — integration, infrastructure, and release QA

**Status:** Not started
**Purpose:** Prove the complete system on a production-like environment.

### Work

- Connect the frontend to the API with loading, success, duplicate, error, and retry behavior.
- Add Dockerfiles and Docker Compose for proxy, frontend, backend, PostgreSQL, and backup process.
- Configure internal database networking, TLS termination, health checks, restart policies, non-root containers where practical, and environment-based secrets.
- Document database backup, restore, rollback, and deployment procedures.
- Add frontend unit/component tests, backend validation/referral/rate-limit tests, migration tests, Playwright conversion tests, and accessibility checks.
- Run production-like builds and inspect browser console, network requests, rendering, image loading, reduced motion, and Core Web Vitals.
- Run Lighthouse and visual comparison against approved Figma frames and local supplementary references.
- Complete privacy, security, and legal-review placeholders.

### Exit criteria

- PRD acceptance criteria are demonstrably met on a production-like build.
- No required network requests fail and no uncaught console errors remain.
- Lighthouse performance and accessibility targets are measured and documented.
- Backup restoration is tested, not merely documented.
- Known limitations and review requirements are recorded before launch.

## Workstream checklist

| Workstream | Depends on | Completion evidence |
| --- | --- | --- |
| Source control and documentation | Repository decision | Tracked/ignored-file policy and README updated |
| Brand assets and tokens | Figma inspection | Export manifest, assets, typography, and token mapping |
| Frontend shell | Brand foundation | Responsive, semantic, accessible page structure |
| Product showcase | Frontend shell/assets | Approved previews with clear capability labels |
| Waitlist UI | Frontend shell | Complete state model and keyboard-tested form |
| API and data model | Service decision | Endpoints, migrations, constraints, and tests |
| Referral flow | API/data model | Deterministic rank and exactly-once referral credit |
| Deployment | Stable services | Compose, proxy, health, secrets, backups, rollback |
| Release QA | All implementation work | Playwright, accessibility, Lighthouse, visual, console/network evidence |

## Sequencing constraints

- Do not finalize brand assets or typography from screenshots alone; Figma is authoritative.
- Do not enable public waitlist intake before duplicate-safe, rate-limited, consent-aware backend behavior exists.
- Do not expose PostgreSQL port 5432 publicly.
- Do not add billing, CRM synchronization, live prescription processing, medical advice, or health-data collection to close roadmap gaps.
- Do not mark a phase complete because files exist; require the relevant exit evidence.
- Update the PRD and affected documentation if an approved requirement changes.

## Verification status

No application commands are currently available to run. Formatting, lint, type-check, unit, integration, Playwright, Lighthouse, production-build, deployment, and restore verification remain future work because the implementation scaffold does not yet exist.
