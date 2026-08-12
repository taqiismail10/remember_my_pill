# Remember My Pill — Project Instructions

## Project identity

This repository contains the Remember My Pill product and its supporting
website, application, APIs, documentation and deployment configuration.

All development, design and documentation work must remain consistent with
the approved Remember My Pill product requirements and brand identity.

## Authoritative sources

Use the following sources in this order:

1. Product and functional source of truth:
   `docs/product/rmp-prd.md`

2. Visual identity and design source of truth:
   `https://www.figma.com/design/LRx7E2phkTHL2DRV1gR86a/RMP-Brand-Identity?node-id=0-1&t=PGo3ZAXibgR7MOOP-1`

3. Local screenshots and visual references:
   `docs/design/references/`

4. Existing repository architecture, components, conventions and tests.

The PRD controls product scope, requirements, user flows, business rules,
privacy requirements, API behaviour, data models and acceptance criteria.

Figma controls approved brand colors, typography, logo, mascot, icon style,
visual hierarchy, spacing, component styling and overall visual identity.

The images under `docs/design/references/` are supplementary visual references.
Do not treat screenshots as a more authoritative source than the original
Figma design.

## Required reading before work

Before planning, coding or changing documentation:

1. Read `docs/product/rmp-prd.md`.
2. Identify the sections relevant to the requested task.
3. Inspect the existing implementation and repository conventions.
4. For visual work, use the Figma MCP server to retrieve the relevant design
   context, variables, components and assets.
5. Inspect the local reference images when they help clarify expected screens
   or marketing presentation.
6. Use Context7 for current framework and package documentation.

Do not begin implementation based only on a short prompt when the PRD or Figma
contains relevant requirements.

## Conflict handling

Do not silently resolve conflicts between the PRD, Figma, screenshots and
existing implementation.

When a conflict exists:

- Product behaviour and scope follow the PRD.
- Visual appearance follows Figma.
- Screenshots provide supplementary intent only.
- Existing code does not override an approved requirement merely because it
  was implemented first.
- Report meaningful conflicts before making irreversible or broad changes.
- When the user intentionally changes a requirement, update the PRD and all
  affected documentation so the repository has one current source of truth.

Do not invent missing product behaviour, pricing, health claims, legal claims,
brand assets or design tokens.

## Brand implementation requirements

- Use the official Remember My Pill logo and mascot from Figma or approved
  repository assets.
- Do not redraw, distort, recolor or replace the logo or mascot without an
  explicit approved requirement.
- Do not approximate brand colors by sampling screenshots when Figma variables
  or approved PRD tokens are available.
- Preserve the approved cream, green, terracotta and supporting brand system.
- Reuse existing design tokens and shared components.
- Avoid generic templates that make the product look unrelated to the RMP
  brand.
- Keep mobile app, website, admin interface and presentation materials visually
  related while respecting the needs of each platform.
- Maintain readable typography, sufficient contrast and accessible interaction
  states.

## UI implementation workflow

For every meaningful interface task:

1. Inspect the relevant Figma node or frame through Figma MCP.
2. Map the design to existing components and tokens.
3. Implement responsive behaviour for mobile, tablet and desktop where
   applicable.
4. Include loading, empty, success, error and disabled states.
5. Verify keyboard navigation and accessible labels.
6. Use Playwright to test important user flows.
7. Use Chrome DevTools to inspect console errors, failed requests, rendering
   issues and performance problems.
8. Compare the implementation against Figma and the approved local references.
9. Do not claim visual completion without browser verification.

Do not create a second design system when the repository already has one.

## Product and healthcare safety

- Treat medication information and user health-related information as
  sensitive.
- Never commit real patient data, prescriptions, credentials or identifiable
  health information.
- Use fictional data in fixtures, screenshots and tests.
- Do not present AI-generated medication interpretation as verified medical
  advice.
- Preserve confirmation, review and error-recovery steps around prescription
  extraction.
- Do not make unsupported HIPAA, PIPEDA, GDPR, clinical accuracy or regulatory
  claims.
- Follow the privacy and data-minimization requirements defined in the PRD.

## Documentation requirements

Documentation is part of the implementation.

Update relevant documentation whenever work changes:

- Product behaviour or scope
- User flows
- Screens or navigation
- API endpoints or payloads
- Database schemas or migrations
- Environment variables
- Authentication or authorization
- Privacy and security behaviour
- Deployment instructions
- Pricing, subscription or referral rules
- Supported devices or browsers
- Third-party services or AI models

Potentially affected files include:

- `docs/product/rmp-prd.md`
- Root `README.md`
- API documentation
- Environment-variable documentation
- Architecture documentation
- Deployment documentation
- Test documentation
- Implementation status or changelog documents

Documentation must describe the real implemented behaviour. Do not document
planned functionality as completed functionality.

When updating the PRD:

- Increment the document version when the change is material.
- Add the date and a concise change summary.
- Preserve prior decisions unless they are explicitly superseded.
- Remove contradictions and duplicated sections.
- Keep requirements testable.
- Update acceptance criteria when behaviour changes.

## Code quality

- Read the appropriate package manifests before selecting commands.
- Follow the repository's existing package manager and architecture.
- Prefer focused changes rather than broad rewrites.
- Do not add a dependency without explaining why it is needed.
- Use current official documentation through Context7 when library behaviour is
  version-sensitive.
- Preserve type safety.
- Do not suppress lint or TypeScript errors merely to make checks pass.
- Never hardcode secrets, tokens, passwords or production credentials.
- Do not modify `.env` files containing real secrets.
- Update `.env.example` when configuration requirements change.

## Verification

Before reporting completion, run all checks available for the affected project,
including where applicable:

- Formatting
- Linting
- Type checking
- Unit tests
- Integration tests
- Playwright end-to-end tests
- Production build
- Mobile and desktop visual checks
- Browser console and network inspection

Fix failures introduced by the work.

If a check cannot run, report exactly which check was skipped and why.

## Git and destructive operations

- Inspect `git status` before and after changes.
- Review the final diff.
- Do not overwrite unrelated user work.
- Do not force-push, rewrite history, delete branches or reset the repository
  without explicit approval.
- Do not apply destructive database changes without explicit approval.
- Do not deploy to production without explicit approval.

## Completion report

Every completion report must state:

1. What was changed
2. Which files were changed
3. Which PRD and Figma requirements were followed
4. Documentation updated
5. Tests and browser verification performed
6. Remaining limitations, conflicts or risks
