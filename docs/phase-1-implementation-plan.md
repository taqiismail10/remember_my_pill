# Phase 1 implementation plan — accessible marketing foundation

> **Status: verified complete on 16 July 2026.** Evidence and final scope are recorded in [the Phase 1 completion report](phase-1-completion-report.md).

**Plan date:** 16 July 2026
**Milestone:** Frontend foundation and accessible marketing shell
**Objective:** Create the first runnable, responsive marketing-site foundation for Remember My Pill using the PRD-required Next.js App Router, TypeScript, Tailwind CSS, official Figma-derived brand assets, and reusable accessible primitives. The milestone intentionally remains a static product showcase; it does not accept or process waitlist data.

## Roadmap validation and selection

### Evidence validated

- `git ls-tree -r HEAD` contains only `.gitignore` and `README.md`; `README.md` contains only the title. There is no current frontend, backend, package manifest, test suite, database, deployment, or existing design system to preserve.
- The working tree did contain untracked `AGENTS.md`, this `docs/` directory, and a local PRD bundle (including the PRD and 13 supplementary PNG references). The audit and roadmap are correct that no product runtime exists.
- The official Figma file is accessible through MCP. Re-checked sources are Page 1 (`0:1`), RMP Palette (`6026:3`), and DP Option 2 + App Icon (`6021:115`). The palette confirms Manrope and the documented warm palette. The app-icon frame provides an official R/pill mark and warm-gradient application treatment.

### Corrections and qualifications

- The audit’s description of the PRD bundle and planning documentation as untracked/ignored is correct, but this is a source-control governance issue, not evidence that the local reference material is absent.
- The planning documents state that a Figma design-system search returned no published variables, components, or styles. This is compatible with the current Figma material being predominantly visual frames; implementation must use a documented repository token mapping instead of claiming a linked Figma variable library.
- The root README and the PRD bundle are still not tracked. This plan does not change `.gitignore`, stage files, or make a source-control decision.

### Why this is the first coherent milestone

The roadmap’s Phase 0 combines a user-owned source-control decision with export preparation. The highest-priority *implementable* milestone is Phase 1: it creates the browser runtime, token layer, shared primitives, semantic structure, and responsive shell required by all later marketing and waitlist work. It has clear UI-only acceptance criteria and does not require unresolved rules for pricing, subscriptions, referrals, persistence, or healthcare processing.

## Exact PRD sections addressed

- §3.1: premium brand presence and quick explanation of the product.
- §5: marketing-website and brand-system portions only; excluded features remain excluded.
- §8: brand-led hero and the product-story presentation only.
- §9.1–§9.10: visual-source hierarchy, color roles, typography, logo handling, surfaces, shape, touch targets, and reduced motion.
- §11 items 1–6: header, hero, problem, solution, product-demo introduction, and Snap → Understand → Act structure.
- §12.1–§12.3: Next.js App Router/TypeScript/Tailwind, responsive behavior, and the applicable accessibility requirements.
- §21: a limited Phase 1 asset-provenance record for the official R/pill mark and app-icon preview.
- §26.1 and §26.3: the applicable visual, responsive, and accessibility criteria.
- §28 Phase 1: frontend foundation and accessible marketing shell.

## Figma sources and assets

| Source | Purpose in this milestone |
| --- | --- |
| Page 1 (`0:1`) | Official brand-file root and asset provenance. |
| RMP Palette (`6026:3`) | Manrope, cream/terracotta/amber identity values, brand voice, and visual hierarchy reference. |
| DP Option 2 + App Icon (`6021:115`) | Official R/pill mark and warm peach/orange gradient app-icon composition. |

The official R/pill artwork will be exported directly from the Figma design context. No logo or mascot will be recreated from a local screenshot or typed as a substitute. The wordmark and mascot are excluded until separately exportable official variants are placed in the repository.

## Current implementation evidence

- No application source or package manager exists in the repository.
- Node.js 22.15.0 and npm 10.9.2 are available locally. Because no package manager is established, npm will be used only for the new `frontend/` package; this is recorded rather than presented as a change to an existing convention.
- No existing components, tokens, tests, or assets exist to reuse.

## Work included

1. Create a Next.js App Router + TypeScript + Tailwind frontend in `frontend/`.
2. Add central RMP CSS tokens, Manrope font loading, global reset/focus/reduced-motion behavior, and responsive layout rules.
3. Export and record the official R/pill icon and an official Figma app-icon preview.
4. Build small reusable primitives for brand link, button, card, section heading, and product-preview panel.
5. Build a single responsive home page with header, hero, problem, solution, Snap → Understand → Act story, and a clearly-labelled product preview.
6. Add unit tests for the primitives/content and a Playwright smoke test for the accessible page shell.
7. Update implementation-status documentation and README only with actual Phase 1 behavior and verification evidence.

## Explicit exclusions

- Waitlist modal/form, consent capture, referral codes, statuses, or API integration.
- Go service, PostgreSQL schema/migrations, rate limiting, admin export, deployment, Docker, or secrets.
- Pricing/subscription rules, billing, CRM, authentication, or medical-document processing.
- Product gallery, caregiver/privacy/FAQ/footer completion beyond the shell content needed for this milestone.
- Mascot, wordmark, social-preview, and product-screen asset exports that cannot be obtained as distinct approved Figma assets in this milestone.
- Any `.env` changes, `.gitignore` changes, staging, commit, or source-control policy decision.

## Files expected to change

- `frontend/package.json`, `frontend/package-lock.json`, Next/TypeScript/Tailwind configuration.
- `frontend/app/layout.tsx`, `frontend/app/page.tsx`, `frontend/app/globals.css`.
- `frontend/components/brand/*`, `frontend/components/marketing/*`, and `frontend/lib/*`.
- `frontend/public/brand/*` and `frontend/public/product-preview/*` (direct Figma exports only).
- `frontend/tests/*` and Playwright configuration.
- `docs/brand-export-manifest.md`, `docs/current-implementation-audit.md`, `docs/prd-gap-analysis.md`, `docs/implementation-roadmap.md`, `docs/design-system-inventory.md`, and `README.md`.

## Dependencies

- Node.js and npm are locally available.
- Next.js, React, TypeScript, Tailwind CSS, ESLint, Vitest, Testing Library, and Playwright are needed to create and verify the required frontend. They are new development dependencies because there was no package manifest; no backend or production-service dependency is introduced.
- Direct Figma MCP asset access for the approved R/pill mark and app-icon preview.

## Risks and controls

| Risk | Control |
| --- | --- |
| Asset provenance is lost | Record Figma node IDs, direct export type, date, and intended use in the export manifest. |
| Visual work uses a non-approved wordmark or mascot | Use only the official R/pill export; omit unavailable distinct assets. |
| Static CTA implies a working waitlist | Use preview/coming-soon language and no submission control. |
| New frontend scaffolding is mistaken for full PRD delivery | Keep the milestone exclusions explicit in code and all updated status documents. |
| Accessibility is deferred | Test headings, landmarks, focus, keyboard links/buttons, motion preference, and mobile layout now. |
| Existing source-control policy is changed accidentally | Do not alter `.gitignore` or stage/commit any file. |

## Implementation sequence

1. Create the plan and capture current repository/Figma evidence.
2. Create the isolated npm-based `frontend/` package and baseline tooling.
3. Export the two approved Figma assets and write their provenance manifest.
4. Define brand tokens and global accessible styles.
5. Implement semantic reusable primitives.
6. Implement the responsive landing-page shell and preview labelling.
7. Add unit and browser smoke coverage.
8. Run formatter, lint, type checking, unit tests, Playwright, production build, and viewport/console/network checks.
9. Update audit, gap, roadmap, inventory, and README with verified—not planned—status.

## Acceptance criteria

- `frontend/` runs as a Next.js App Router TypeScript application with Tailwind CSS and no extra competing design system.
- The marketing shell contains semantic header, main, navigation, and section landmarks with logical heading order.
- The official Figma R/pill asset is used and documented; no recreated logo/mascot is used.
- Tokens implement the verified cream, soft-cream, terracotta, burnt-sienna, amber, and Manrope identity baseline, with documented text/focus values.
- Hero, problem, solution, Snap → Understand → Act, and product-preview content remain readable and responsive from 320 through 1440 px without horizontal page overflow.
- Touch targets and focus styles are visible; the layout has a reduced-motion path; non-color visual cues accompany the three-step story.
- The page labels all product visuals as a preview and does not collect data or claim live prescription processing.
- Formatter, lint, TypeScript, unit tests, Playwright, and production build pass; browser console and required local-network requests show no failures.

## Verification plan

- Run the package’s format check, lint, TypeScript check, unit tests, Playwright test, and production build.
- Run the production application and use Playwright at 320, 375, 390, 430, 768, 1024, 1280, and 1440 px to check page overflow, visibility, and keyboard focus.
- Inspect the rendered page with Chrome DevTools for console errors, failed requests, accessibility tree/landmarks, and responsive rendering.
- Capture non-production screenshots under `frontend/test-results/` or Playwright output and compare the palette, mark, and warm/airy composition with Figma nodes `6026:3` and `6021:115`.

## Rollback considerations

The work is isolated to a new `frontend/` directory and untracked documentation. Removing the new directory and its documentation status entries returns the repository to its current tracked runtime state. No database, environment, deployment, source-control policy, or external service state is touched.
