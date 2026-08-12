# PRD gap analysis

## Current-status update — 16 July 2026

The pre-implementation baseline below is superseded for PRD Phase 1 only: the frontend foundation, token layer, approved R/pill asset use, responsive static marketing shell, and its recorded checks are verified complete. The approved Phase 2 scope combines minimal waitlist UI, Go API, PostgreSQL persistence, validation, duplicate safety, basic rate limiting, and end-to-end checks; referral and administration features remain Phase 3 gaps.

**Baseline:** Remember My Pill PRD v1.2.0, updated 16 July 2026
**Assessment date:** 16 July 2026
**Assessment rule:** A requirement is counted as implemented only when repository evidence exists. The PRD controls product scope and behavior; the official Figma file controls visual identity.

## Executive summary

The repository contains the product requirements and reference material, but no application implementation. Every product, infrastructure, and QA capability defined by the PRD remains missing. The only implementation-adjacent work present is planning documentation in the working tree.

The official Figma file has now been inspected through Figma MCP. Its brand identity is available for implementation planning, but approved exported assets, a token mapping, and a Figma export manifest have not yet been added to the repository.

## Current repository evidence

- Tracked files are limited to .gitignore and README.md.
- README.md contains only the project title.
- The working tree contains AGENTS.md, docs/, and the ignored PRD/reference bundle.
- There is no frontend/, backend/, package.json, go.mod, Docker configuration, database migration, test suite, CI workflow, or environment template.
- The local reference bundle contains nine mobile PNGs and four presentation PNGs under docs/design/references/.
- The latest commit is 0d53b11, Remove PRD assets from repository.

## Gap matrix

| PRD area | Required outcome | Status | Repository evidence | Gap and next decision |
| --- | --- | --- | --- | --- |
| Release scope | Responsive marketing site, product showcase, waitlist, referral system, backend, and VPS deployment; no public prescription processing or full browser app. | Missing | PRD scope sections; no application directories or manifests. | Establish the frontend/backend scaffold without expanding waitlist scope. |
| Brand foundation | Official logo, mascot, typography, colors, icons, illustrations, tokens, and export manifest. | Partially specified / not implemented | Figma file inspected; no implementation asset or token directories. | Export approved assets and record Figma node provenance before visual implementation. |
| Marketing information architecture | Header, hero, problem, solution, product demo, Snap–Understand–Act, caregiver, privacy, plans preview, waitlist, FAQ, and footer. | Missing | PRD marketing structure; no page or component source. | Build the accessible marketing shell after asset/token setup. |
| Product showcase | Approved phone mock-ups, product gallery, and interactive/static prescription-to-reminder story. | Missing | Local screenshots only; no public asset pipeline or showcase code. | Keep all demos clearly labelled as previews; do not imply live prescription processing. |
| Responsive behavior | Supported widths from 320 px through 1440+ px, no overflow, deliberate stacking, usable CTAs and navigation. | Missing | No browser app or responsive CSS. | Implement and verify at every PRD-required width. |
| Accessibility | Semantic structure, keyboard navigation, focus states, labelled errors, status announcements, reduced motion, and non-color state cues. | Missing | No UI or accessibility test setup. | Treat accessibility as a foundation requirement, not a final pass. |
| Waitlist form | Name, email, consent, optional referral code, honeypot, and all required validation/loading/success/error states. | Missing | No form, route, or state model. | Implement without collecting medication or other health information. |
| Referral behavior | Unique code, safe resolution, deterministic rank, successful-referral credit, duplicate protection, and self-referral prevention. | Missing | No backend or persistence. | Define and test rules at the service layer; do not hard-code future rewards in the frontend. |
| API | Public waitlist/referral endpoints, opaque status-token flow, admin export, and health endpoint. | Missing | No Go server or API documentation beyond the PRD. | Implement the documented endpoints with safe duplicate behavior and rate limits. |
| Database | PostgreSQL schema for waitlist entries and referral events, constraints, token hashing, timestamps, and migrations. | Missing | No database files or migrations. | Create migrations before public intake; protect the database from public exposure. |
| Privacy/security | Data minimization, TLS, secrets outside source control, restricted CORS, secure headers, request limits, logs, deletion explanation, and protected export. | Missing | No infrastructure or security configuration. | Complete a privacy/security design before enabling waitlist submissions. |
| Deployment | Separate frontend/backend services, internal PostgreSQL network, proxy/TLS, health checks, restart policies, backups, and restore procedure. | Missing | No Compose, proxy, Dockerfile, or deployment docs. | Add production configuration after service boundaries and health behavior exist. |
| QA | Unit, component, API, migration, Playwright, accessibility, visual, console/network, and Lighthouse verification. | Missing | No package tooling, test runner, or runnable application. | Add test harnesses alongside implementation and record real results. |
| Documentation | README, architecture, API, environment, privacy, deployment, rollback, test, and implementation-status documentation. | Partial | Planning docs exist; README is title-only. | Document actual behavior incrementally; do not mark planned work complete. |

## Figma and brand gap analysis

The official Figma file is the visual identity source of truth. Figma MCP inspection confirmed:

- Manrope typography;
- Terracotta #B5541F;
- Burnt Sienna #C1592D;
- Cream #F7EAD9;
- Soft Cream #FBF1E4;
- Amber Accent #E08E3E;
- Coral/Peach Accent #F0A878;
- R/pill logo and app-icon variants;
- official 2D and 3D mascot artwork;
- supporting reminder, schedule, completed, pill, notification, and care icons;
- warm, reassuring, simple, trustworthy brand voice.

Relevant Figma frames include RMP Pallete 1 (6026:3), Youtube Banner (6021:204), Facebook Cover (6023:220), 2D Mascot (6023:243), 3D Mascot (6023:241), and 3D Logo Render (6023:242).

The Figma design-system search returned no published variables, components, or styles. The remaining gap is therefore not Figma access; it is exporting and mapping the authoritative visual material into the repository. The implementation must not redraw the logo/mascot or use local screenshots as final production assets.

## Functional dependency chain

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

## Material gaps by implementation phase

### Phase 1: Brand foundation

Missing:

- frontend project scaffold and package tooling;
- CSS/design-token implementation;
- official asset export directory;
- brand export manifest;
- reusable logo, button, card, navigation, and section primitives.

Exit evidence:

- approved Figma assets are present with source-node provenance;
- tokens are mapped and documented;
- shared primitives render responsively and accessibly.

### Phase 2: Marketing and conversion experience

Missing:

- marketing page sections;
- product preview/mock-up system;
- caregiver, privacy, FAQ, plans-preview, and footer sections;
- waitlist form and complete state model.

Exit evidence:

- all required sections exist;
- previews are clearly distinguished from live product behavior;
- form states are keyboard-accessible and do not collect health data;
- required viewport checks pass without horizontal overflow.

### Phase 3: Waitlist and referral services

Missing:

- Go API;
- PostgreSQL migrations and data access;
- validation, duplicate handling, status tokens, referral accounting, and rank;
- rate limiting, honeypot, request-size limits, restricted CORS, and admin export protection.

Exit evidence:

- new registrations, duplicates, valid/invalid referrals, self-referrals, and rate limits are covered by automated tests;
- no endpoint exposes sensitive email-existence information;
- consent version and timestamp are persisted securely.

### Phase 4: Release readiness

Missing:

- Docker Compose and reverse proxy;
- health checks, restart policies, backups, restore testing, and rollback documentation;
- frontend/backend/API/migration/Playwright/accessibility/visual/Lighthouse checks;
- setup, architecture, API, privacy, deployment, and test documentation.

Exit evidence:

- PRD acceptance criteria are demonstrated on a production-like build;
- no required network failures or uncaught console errors remain;
- legal/privacy copy is clearly marked for required review.

## Source conflicts and decision rules

1. The PRD is authoritative for product scope, functional behavior, privacy requirements, API behavior, data models, and acceptance criteria.
2. Figma is authoritative for visual identity, including logo, mascot, typography, colors, icons, composition, and visual assets.
3. Local screenshots are supplementary references only. Orange and green screen variants should not override the PRD/Figma hierarchy.
4. The PRD’s product UI roles, including deep-green actions and state colors, must be mapped alongside the warm Figma marketing identity. Any token difference must be documented rather than silently approximated.
5. The waitlist release must not add live public prescription processing, medical advice, diagnosis, billing, CRM automation, or collection of medication/health data.

## Recommended next actions

1. Decide whether the currently ignored PRD bundle and untracked documentation should be included in source control.
2. Establish frontend/ and backend/ with the selected package manager and Go toolchain.
3. Export the approved Figma logo, mascot, icons, illustrations, product screens, typography, and color references.
4. Create the PRD-required brand export manifest and repository token mapping.
5. Build the shared accessible primitives and marketing shell.
6. Implement the waitlist and referral backend before enabling public conversion.
7. Add automated and browser-based verification, then update the README and operational documentation with actual behavior.

## Verification status

- Repository structure, Git state, manifests, README, planning documents, PRD, and local references were inspected.
- Figma MCP inspection succeeded for the supplied official file and key brand frames.
- No application tests, builds, lint, type checks, Playwright flows, Lighthouse runs, or production checks were available because no implementation exists.
