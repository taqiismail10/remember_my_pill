# Design-system inventory

## Current-status update — 11 August 2026 (superseded — see below)

The full marketing redesign (11 August 2026) implemented the token system, component library, and every section this document previously listed as missing — see `frontend/app/globals.css` for the token layer, `frontend/components/` for the component library, and `docs/brand-export-manifest.md` for current asset provenance. The mascot is now exported and in use (`frontend/public/brand/mascot/mascot-wave.png`); the wordmark and supporting icon set remain unexported from Figma, so the documented decisions in the brand manifest apply. The remaining sections of this document are the original Phase 1 gap analysis and are kept for historical context — they no longer describe the current implementation.

## Current-status update — 16 July 2026 (historical)

The Phase 1 repository implementation now contains the documented CSS token baseline, Manrope loading, official R/pill SVG mark, app-icon preview, and shared marketing primitives. The asset inventory below remains incomplete for wordmark, mascot, supporting icons, illustrations, and product screenshots; no Phase 2 component is implemented.

**Inventory date:** 16 July 2026
**Product source of truth:** Remember My Pill PRD v1.2.0
**Visual identity source of truth:** Official RMP Brand Identity Figma file
**Implementation status:** No frontend or asset pipeline exists yet.

This inventory separates verified Figma identity material from PRD product/UI roles and from the supplementary local PNG references. It is not a substitute for the Figma export manifest required before visual implementation.

## Retrieval status

| Source | Status | Implication |
| --- | --- | --- |
| Official Figma brand file | Inspected through Figma MCP | Exact identity frames, palette board, typography, logo variants, mascot artwork, and application previews are available for implementation planning. |
| Published Figma design-system assets | No published variables/components/styles returned | Export and map the authoritative visual frames explicitly; do not assume a reusable library is available. |
| PRD brand/product specification | Read in full | Defines product color roles, accessibility, responsive behavior, asset requirements, and implementation constraints. |
| Local reference bundle | 13 PNGs inspected | Useful supplementary mobile and presentation references; not authoritative over Figma and not implementation assets. |
| Repository implementation | Not started | No public asset directory, tokens, components, package manifest, or frontend exists. |

## Figma frame inventory

| Frame | Node ID | Inventory use |
| --- | --- | --- |
| Page 1 | 0:1 | Root brand file canvas |
| RMP Pallete 1 | 6026:3 | Palette, typography, logo, mascot, icon, and application reference board |
| Youtube Banner | 6021:204 | Social/banner composition and approved brand messaging |
| Facebook Cover | 6023:220 | Social-cover composition and responsive safe-area reference |
| 2D Mascot | 6023:243 | Official 2D mascot artwork |
| 3D Mascot | 6023:241 | Official 3D mascot artwork |
| 3D Logo Render | 6023:242 | Official 3D logo treatment |
| Glass Logo | 6023:240 | Secondary logo treatment reference |
| DP Option 1 | 6019:73 | Orange-gradient display-picture option |
| DP Option 2 + App Icon | 6021:115 | Warm-gradient display-picture and app-icon option |
| Alt DP | 6021:2 | Green-gradient display-picture option |
| Flat Icon Option 1 | 6020:275 | Flat orange app-icon option |
| Flat Icon Option 2 | 6021:69 | Alternate flat orange app-icon option |

## Verified identity tokens

The following values are read from the Figma RMP Pallete 1 board and are the current identity baseline for implementation mapping.

| Token label | Value | Role | Repository status |
| --- | --- | --- | --- |
| Terracotta | #B5541F | Primary logo/brand-story color | Not implemented |
| Burnt Sienna | #C1592D | Deeper warm accent and illustration depth | Not implemented |
| Cream | #F7EAD9 | Warm marketing background/surface | Not implemented |
| Soft Cream | #FBF1E4 | Light panel/card surface | Not implemented |
| Amber Accent | #E08E3E | Supporting warm emphasis | Not implemented |
| Coral/Peach Accent | #F0A878 | Secondary warm accent | Not implemented |
| Typography | Manrope | Headings, subheadings, and body typography | Not implemented |

The Figma board also communicates the brand voice as warm, reassuring, simple, and trustworthy.

## Logo inventory

### Verified assets and variants

- Primary R/pill logo.
- Horizontal Remember My Pill wordmark/lockup.
- Black-and-white logo variant.
- Inverted logo variant.
- Warm gradient logo variant.
- R monogram/app icon.
- Flat and gradient app-icon options.
- Glass and 3D logo treatments for approved campaign/application contexts.

### Implementation rules

- Use official exported artwork; do not rebuild the mark with text, emoji, or utility icons.
- Preserve proportions, clear space, and approved variants.
- Do not stretch, rotate, crop, recolor, or add unapproved effects.
- Use SVG wherever the export supports it.
- Provide accessible text for linked logos, such as Remember My Pill home.
- Use only the appropriate light/dark/background variant for each surface.
- Record source node, variant, export date, and intended use in the export manifest.

**Current status:** Missing from the repository. No implementation-owned logo files exist.

## Mascot inventory

### Verified assets

- Official 2D mascot artwork.
- Official 3D mascot artwork.
- Mascot used in application previews and reminder/support illustrations.
- Mascot imagery associated with friendly, low-risk educational or success contexts.

### Implementation rules

- Use the official artwork rather than recreating or approximating the mascot.
- Use suitable alt text when informative and empty alt text when decorative.
- Do not place a playful mascot pose beside urgent medication warnings.
- Do not use the mascot as a replacement for a status label or warning icon.
- Avoid continuous animation and respect prefers-reduced-motion.

**Current status:** Missing from the repository. No implementation-owned mascot files exist.

## Supporting icon and illustration language

The Figma identity board shows a warm line/filled icon family for:

- reminder;
- schedule;
- completed;
- pill;
- notification; and
- care.

Use the approved identity icon treatment where available. Lucide may be used for utility controls only when no brand-specific icon exists, with stroke weight and optical size normalized to the approved language.

**Current status:** Missing from the repository. No public icon or illustration pipeline exists.

## Product UI token baseline

The PRD defines additional product/UI roles that must be mapped alongside the Figma identity palette. These values remain a PRD baseline until the implementation token mapping is documented and reviewed against the Figma direction.

| Token | PRD baseline | Intended role | Status |
| --- | --- | --- | --- |
| brand.cream | #F7EAD9 | Marketing background/warm surface | Pending token implementation |
| brand.softCream | #FBF1E4 | Cards/panels/secondary surfaces | Pending token implementation |
| brand.terracotta | #B5541F | Logo/storytelling/marketing emphasis | Pending token implementation |
| brand.burntSienna | #C1592D | Stronger warm accent and illustration depth | Pending token implementation |
| brand.amber | #E08E3E | Supporting emphasis and schedule labels | Pending token implementation |
| product.deepGreen | #155E55 | Primary product actions and selected navigation | Pending token implementation |
| product.deepGreenDark | #0B4C47 | High-contrast product surfaces | Pending token implementation |
| product.coral | #FF6757 | Central add/high-attention action | Pending token implementation |
| surface.app | #F4F8F8 | Mobile-app background | Pending token implementation |
| surface.white | #FFFFFF | Cards, inputs, and modal surfaces | Pending token implementation |
| text.ink | #202033 | Primary text | Pending token implementation |
| text.muted | #646774 | Secondary text | Pending token implementation |
| state.success | #00C96B | Completed dose and success | Pending token implementation |
| state.warning | #FFB400 | Due/attention indicator | Pending token implementation |
| state.error | #FF3348 | Error or missed-dose state | Pending token implementation |

The warm Figma identity palette and the PRD product-state palette are complementary roles, not permission to mix competing primary CTAs. Deep green should anchor product actions; terracotta/orange should carry brand storytelling and marketing emphasis; coral should remain limited to clearly defined high-attention actions.

## Layout, shape, and motion requirements

These requirements come from the PRD and must be mapped to Figma measurements before implementation:

| Area | Required direction | Current status |
| --- | --- | --- |
| Marketing surfaces | Light cream/near-white, airy, generous negative space | Not implemented |
| Product surfaces | Cool, very light gray-blue/green app background | Not implemented |
| Cards | White/soft-cream fills, fine neutral borders, restrained shadows | Not implemented |
| Large panel radius | Generally 20–32 px | Not implemented |
| Pills/tags/segmented controls | Fully rounded | Not implemented |
| Touch targets | At least 44 × 44 CSS px | Not implemented |
| Typography | Exact Figma family/weights; readable line-height and measure | Not implemented |
| Motion | Short fades/progression; no distracting continuous motion | Not implemented |
| Reduced motion | No essential information depends on animation | Not implemented |

## Component and behavior inventory

| Component/area | Required characteristics | Status |
| --- | --- | --- |
| Header/navigation | Semantic navigation, linked official logo, responsive touch/keyboard behavior, primary waitlist CTA | Missing |
| Hero | Value proposition, product preview, primary waitlist CTA, brand-led composition | Missing |
| Product demo | Capture, analysing, review, grouped reminders; preview language | Missing |
| Product gallery | Approved screens in consistent mock-ups, preserved aspect ratio, useful captions/alt text | Missing |
| Caregiver section | Multi-person support and family value explanation | Missing |
| Privacy section | Plain-language collection boundaries, retention, sharing, deletion request path | Missing |
| Waitlist form | Name, email, consent, optional referral code, honeypot, accessible validation | Missing |
| Waitlist states | Loading, success, duplicate, invalid referral, rate limit, network failure, copied link | Missing |
| Dialog | Keyboard focus trap/return, labelled title, status/error announcement | Missing |
| Product status visuals | Green primary state, distinct due/missed/completed/disabled states, not color-only | Missing |
| FAQ/footer | Accessible disclosure and contact/privacy/terms placeholders | Missing |
| Motion | Purposeful transitions and reduced-motion path | Missing |

## Supplementary local reference register

All files are under docs/design/references/ and remain subordinate to Figma.

| Group | Files | Intended use |
| --- | --- | --- |
| Earlier orange exploration | mobile-home-orange.png; mobile-add-medication-orange.png; mobile-profile-orange.png; mobile-reminder-orange.png | Historical/exploratory orange-accent mobile states |
| Green product references | mobile-home-green-state-1.png; mobile-home-green-state-2.png; mobile-add-medication-green.png; mobile-profile-green.png; mobile-reminder-green.png | Product-state, capture, profile, and reminder patterns |
| Marketing/deck references | deck-data-sources.png; deck-our-solution.png; deck-product-demo.png; deck-business-model.png | Marketing composition, showcase, cards, typography, and icon context |

These PNGs are not exported production assets and must not be treated as final logo, mascot, typography, or token sources.

## Required export manifest

Before visual implementation is considered complete, create a manifest covering:

- horizontal, stacked, monogram, and approved light/dark/single-color logo variants;
- favicon/app icon and social preview;
- mascot poses and variants;
- prescription, pill, calendar, family, and reminder illustrations;
- approved product screenshots/mock-ups;
- color variables and typography styles;
- source Figma node/component URL;
- export date;
- intended use;
- background compatibility;
- format and resolution.

Preferred formats are SVG for logos/icons, SVG or PNG fallback for mascot artwork, SVG/WebP for illustrations, and WebP/PNG at appropriate density for screenshots.

## Implementation gaps

1. No asset files are present in a frontend public pipeline.
2. No CSS/JSON token file exists.
3. No font loading or typography configuration exists.
4. No shared components exist.
5. No product or marketing screen is implemented.
6. No asset provenance/export manifest exists.
7. No automated visual or accessibility verification exists.

## Next actions

1. Resolve the source-control decision for the ignored PRD bundle and untracked documentation.
2. Export and document the approved Figma assets listed above.
3. Implement the token mapping with separate marketing, product, surface, text, and state roles.
4. Build reusable accessible primitives before page-specific sections.
5. Validate responsive, keyboard, reduced-motion, contrast, and visual fidelity requirements as each component is introduced.
6. Update this inventory when an asset, token, component, or behavior is actually implemented.

## Verification status

Figma MCP inspection succeeded for the official file and relevant brand frames. The Figma design-system search returned no published variables, components, or styles. No frontend implementation or asset pipeline exists, so no build, lint, browser, Playwright, Lighthouse, or visual-regression checks were available to run.
