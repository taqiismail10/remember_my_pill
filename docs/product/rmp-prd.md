# Product Requirements Document

## Remember My Pill (RMP) — Web Platform, Product Showcase & Waitlist

### Document Control

| Field | Detail |
| --- | --- |
| **Product name** | Remember My Pill (RMP) |
| **Document version** | 1.2.0 |
| **Last updated** | 16 July 2026 |
| **Primary deliverable** | Premium marketing website, product showcase and viral waitlist |
| **Frontend** | Next.js App Router, TypeScript, Tailwind CSS, shadcn/ui, Framer Motion, Lucide Icons |
| **Backend** | Go with Chi or Gin |
| **Database** | PostgreSQL |
| **Deployment** | Linux VPS, Docker Compose, Caddy or Nginx, automatic TLS |
| **Primary markets** | United States and Canada |
| **Privacy posture** | Data-minimizing, privacy-first, HIPAA-aware, PIPEDA-aware and GDPR-aligned; final legal review required |
| **Design source of truth** | [RMP Brand Identity — Figma](https://www.figma.com/design/LRx7E2phkTHL2DRV1gR86a/RMP-Brand-Identity?node-id=0-1&t=PGo3ZAXibgR7MOOP-1) |

### Version 1.2.0 Update Summary

- Added the official Figma brand file as the source of truth for colors, logo, mascot, typography and visual assets.
- Added a complete brand identity and UI design specification derived from the supplied app and presentation references.
- Replaced the previous dark-mode-first direction with a light-first, warm-cream brand experience.
- Clarified the division between deep green product UI actions and terracotta/orange brand accents.
- Added logo, mascot, image-export, accessibility and visual-fidelity requirements.
- Added explicit responsive, screenshot-comparison and browser-testing acceptance criteria.
- Added design-token and brand-asset folders to the repository structure.
- Removed the duplicated raw-markdown copy that previously appeared at the end of the document.
- Replaced hard-coded production secrets with safe placeholders and strengthened deployment guidance.

---

## 1. Product Overview

**Remember My Pill** is a privacy-first medication support product designed to help people understand prescription instructions, organise medication schedules and follow daily treatment routines with less confusion.

The product experience centres on three actions:

1. **Snap** — capture or upload a prescription.
2. **Understand** — convert the prescription into clear, structured instructions.
3. **Act** — organise medicines into grouped reminders and mark doses as taken.

This PRD defines the public-facing **Remember My Pill web platform**. The website serves as the brand lighthouse for the mobile product, communicates the product’s value, demonstrates key workflows, captures high-intent waitlist registrations and supports a referral-driven early-access campaign.

The web platform must look and feel like the same product shown in the supplied mobile UI screens and presentation deck. It must not use a generic health-tech template or introduce a competing visual language.

---

## 2. Problem Statement

Medication instructions can be difficult to read, interpret and coordinate, especially when a person manages several medicines or helps multiple family members. The experience becomes more difficult when instructions are handwritten, use clinical shorthand or are spread across different prescriptions and schedules.

Before mobile-store launch, Remember My Pill needs a trustworthy digital presence that can:

- explain the product clearly;
- demonstrate the prescription-to-reminder workflow;
- present the brand consistently;
- collect qualified early users;
- create a simple referral loop; and
- protect visitor data through strict minimisation and secure infrastructure.

---

## 3. Product Goals

### 3.1 Primary Goals

- **Establish a premium brand presence:** Create a high-fidelity website consistent with the official RMP Figma identity and product UI.
- **Explain the product quickly:** Communicate the “Snap, Understand, Act” flow within the first screen and first 30 seconds of browsing.
- **Show, not only describe:** Use real product mock-ups, interactive demonstrations and approved mobile screens.
- **Capture high-intent waitlist entries:** Collect a minimal set of validated user details.
- **Create organic growth:** Generate unique referral links and rank improvements for successful referrals.
- **Protect user trust:** Avoid behavioural advertising, unnecessary trackers and collection of health information during the waitlist phase.

### 3.2 Non-Goals for the Waitlist Release

- Processing real prescriptions through the public marketing website.
- Providing medical advice or diagnosis.
- Collecting medical histories, insurance information or medication details from waitlist users.
- Charging users through Stripe, PayPal or app-store billing.
- Launching a full CRM or marketing-automation integration.
- Reproducing the complete mobile app as a browser application.

---

## 4. Target Users

### 4.1 Proactive Patient

A health-conscious user who manages one or more daily medicines and wants a simple, reliable routine.

### 4.2 Family Caregiver

A person coordinating medicines for parents, children, partners or other family members and needing grouped schedules and clear ownership.

### 4.3 Older or Low-Confidence Digital User

A user who benefits from plain language, large touch targets, predictable navigation and accessible reminders.

### 4.4 Privacy-Conscious User

A user who wants clear information about where data is stored, what is collected and what is never collected.

---

## 5. Scope

| Area | Included in v1 | Excluded from v1 |
| --- | --- | --- |
| **Marketing website** | Responsive landing page, solution overview, product demo, feature sections, trust/privacy messaging, FAQ and waitlist CTAs | Full authenticated web app |
| **Product showcase** | Approved phone mock-ups, interactive prescription flow, reminder examples and mobile screen gallery | Live medical-document processing |
| **Brand system** | Official logo, mascot, Figma colors, typography, icon language, gradients, cards and motion rules | Unapproved logo or mascot redesigns |
| **Waitlist** | Name, email, consent, unique referral code, referral rank and share link | Phone number, medical profile or billing data |
| **Backend** | Validation, duplicate handling, referral accounting, rate limits, health check and restricted export | CRM sync and production billing |
| **Deployment** | Docker Compose on a single Linux VPS with TLS and database backups | Kubernetes or multi-region infrastructure |

---

## 6. Success Metrics

| Metric | Target |
| --- | --- |
| Landing-page Lighthouse Performance | 95 or higher on representative mobile and desktop runs |
| Lighthouse Accessibility | 95 or higher, with no critical accessibility violations |
| Core waitlist form completion | At least 70% after the user opens the form |
| Invalid or duplicate submissions | Returned with clear, non-technical feedback |
| Referral link generation | 100% for successfully created waitlist entries |
| Mobile layout defects | No critical overflow, clipped CTAs or unusable controls at supported widths |
| Visual fidelity | Major brand screens and marketing components match the approved Figma direction |
| Privacy | No health information collected by the waitlist form |

Conversion-rate targets should be calibrated after real traffic is available rather than invented before launch.

---

## 7. Primary User Journey

```text
Visitor opens remembermypill.com
        |
        v
Sees approved RMP logo, product promise and mobile-product preview
        |
        v
Explores the Snap -> Understand -> Act demonstration
        |
        v
Reviews benefits for patients and caregivers
        |
        v
Clicks Join the Waitlist / Secure Early Access
        |
        v
Enters name, email and consent
        |
        v
Backend validates, deduplicates and stores the entry
        |
        v
Visitor receives rank, referral link and sharing options
```

### 7.1 Required Conversion States

The interface must provide designed states for:

- initial form;
- field validation;
- submitting/loading;
- success;
- existing email;
- invalid referral code;
- rate-limited request;
- network failure; and
- copied referral link.

---

## 8. Core Features

| Feature | Requirement | Priority |
| --- | --- | --- |
| **Brand-led hero** | Official logo, concise value proposition, product mock-up and primary waitlist CTA | Critical |
| **Interactive product story** | Demonstrate prescription capture, AI interpretation, review and grouped reminders | High |
| **Product screen gallery** | Use approved RMP mobile screens inside consistent device mock-ups | High |
| **Caregiver value section** | Explain multi-person medication support and sharing | High |
| **Privacy section** | Clearly state data-minimisation and waitlist-data boundaries | Critical |
| **Waitlist modal/page** | Accessible form with validation, consent and success state | Critical |
| **Referral system** | Unique code, rank, referral count and shareable URL | High |
| **FAQ** | Explain availability, privacy, supported markets and early access | Medium |
| **Language architecture** | English-first; content structure prepared for future localisation | Medium |
| **Secure admin export** | Authenticated CSV export for authorised operators | Critical |

---

## 9. Brand Identity and Design System

### 9.1 Source of Truth

The official Figma file is authoritative for:

- logo artwork and permitted lock-ups;
- mascot artwork and poses;
- primary and secondary color values;
- typography and font weights;
- spacing and component proportions;
- illustration style;
- icon treatment;
- corner radii;
- shadows;
- gradients; and
- presentation/marketing composition.

**Figma:** [RMP Brand Identity](https://www.figma.com/design/LRx7E2phkTHL2DRV1gR86a/RMP-Brand-Identity?node-id=0-1&t=PGo3ZAXibgR7MOOP-1)

When this PRD, a screenshot and the Figma file differ, the latest approved Figma component or variable wins. A developer must not approximate the logo, mascot or brand mark with text, emoji, a Lucide icon or a hand-drawn substitute.

### 9.2 Brand Character

The visual personality should feel:

- warm and reassuring;
- modern but not clinical or cold;
- clear and accessible;
- family-friendly;
- privacy-conscious;
- calm during reminders and error states; and
- premium without looking luxury-oriented or decorative.

### 9.3 Color Roles

The following palette is the implementation baseline visible in the supplied references. Exact production values must be exported from the Figma variables before final implementation.

| Token | Baseline value | Intended use |
| --- | --- | --- |
| `brand.cream` | `#F7EAD9` | Marketing background and warm brand surface |
| `brand.softCream` | `#FBF1E4` | Cards, panels and secondary surfaces |
| `brand.terracotta` | `#B5541F` | Brand emphasis, logo applications and marketing accents |
| `brand.burntSienna` | `#C1592D` | Stronger orange accents, illustration depth and selected highlights |
| `brand.amber` | `#E08E3E` | Supporting emphasis, schedule labels and warm indicators |
| `product.deepGreen` | `#155E55` | Primary product actions, selected navigation and positive brand anchors |
| `product.deepGreenDark` | `#0B4C47` | High-contrast product surfaces and streak cards |
| `product.coral` | `#FF6757` | Central add action and limited high-attention controls |
| `surface.app` | `#F4F8F8` | Main mobile-app background |
| `surface.white` | `#FFFFFF` | Cards, inputs and modal surfaces |
| `text.ink` | `#202033` | Primary text |
| `text.muted` | `#646774` | Secondary text |
| `state.success` | `#00C96B` | Completed dose and success state |
| `state.warning` | `#FFB400` | Due/attention indicators |
| `state.error` | `#FF3348` | Error or missed-dose state |

#### Color Rules

- Use **deep green** for primary product actions and selected product navigation.
- Use **terracotta/orange** for the logo, brand storytelling, illustrations and marketing emphasis.
- Use **coral** sparingly for the central add button or one clearly defined high-attention action.
- Do not use orange and green as competing primary CTAs within the same visual hierarchy.
- Do not place low-contrast gray text on cream backgrounds.
- All text and interactive states must meet WCAG AA contrast requirements.
- Decorative gradients must not reduce readability or become full-page mesh backgrounds.

### 9.4 Backgrounds and Surfaces

- The primary marketing background is a soft cream or near-white surface.
- The product UI uses a cool, very light gray-blue/green background.
- The centre of each section should remain light and airy with generous negative space.
- An oversized translucent R mark may be used as a subtle watermark behind marketing content.
- Cards use white or soft-cream fills, fine neutral borders and restrained shadows.
- Avoid dark-mode-first presentation unless a separate approved dark theme is added to Figma.

### 9.5 Typography

- Use the exact font family and weights defined in the Figma file.
- If the official web font is not available at implementation time, use a modern system sans fallback while preserving Figma metrics as closely as possible.
- Headings should be bold, compact and highly legible.
- Body copy should use comfortable line height and avoid overly narrow measure.
- Sentence case is preferred for headings and controls.
- The wordmark must remain an exported brand asset rather than typed text.

### 9.6 Logo Requirements

Required exports from Figma:

- primary horizontal logo;
- stacked logo;
- R monogram;
- light-background variant;
- dark-background variant, when approved;
- single-color variant, when approved;
- favicon/app-icon source; and
- social preview mark.

Rules:

- Preserve the official proportions and clear space.
- Never stretch, rotate, crop or rebuild the logo.
- Never recolor the logo outside approved variants.
- Never add unapproved shadows, outlines or gradients.
- Use SVG for the website wherever possible.
- Provide meaningful accessible text for the linked logo, such as `Remember My Pill home`.

### 9.7 Mascot Requirements

The official mascot in the Figma file is part of the RMP identity and should be exported as approved assets rather than recreated.

Approved use cases include:

- onboarding support;
- empty states;
- successful waitlist registration;
- friendly explanations;
- FAQ illustrations; and
- low-risk educational moments.

Mascot restrictions:

- Do not place the mascot beside urgent medication warnings in a playful pose.
- Do not use the mascot as a substitute for a clear label or status icon.
- Do not animate the mascot continuously.
- Motion must respect `prefers-reduced-motion`.
- All mascot images need suitable alt text, or empty alt text when purely decorative.

### 9.8 Icons and Illustrations

- Use the approved icon style shown in the Figma file and supplied designs.
- Lucide may be used for utility icons when no brand-specific icon exists.
- Use consistent stroke width, optical size and corner treatment.
- Pill, calendar, family, reminder and prescription illustrations should use the approved green/terracotta language.
- Avoid mixed icon libraries on the same page unless visually normalised.

### 9.9 Shape, Radius and Spacing

- Use generous rounded cards and controls consistent with the supplied mobile designs.
- Large panels should generally use a radius in the 20–32 px range.
- Pills, tags and segmented controls should use fully rounded shapes.
- Touch targets must be at least 44 × 44 CSS pixels.
- Use a consistent spacing scale based on the Figma variables.
- Maintain generous vertical rhythm; do not compress sections to imitate dashboard density on the marketing website.

### 9.10 Motion

- Motion should clarify hierarchy and progression rather than decorate every element.
- Preferred effects: short fades, small upward transitions, gentle mock-up movement and controlled progress transitions.
- Avoid large parallax, continuous floating or motion that can distract from medication information.
- Provide a reduced-motion path with no essential information dependent on animation.

---

## 10. Product UI Reference Requirements

The supplied app screens establish the following recurring patterns:

### 10.1 Home Screen

- Friendly greeting and RMP wordmark.
- Person/profile selector card.
- Summary chips for weight, daily pill count and sleep.
- Grouped medication schedule by time of day.
- Clear dose-state indicators.
- One large state-aware completion control.
- Persistent bottom navigation with a central add action.

### 10.2 Add Medication

- Clear back action and page title.
- Person context chip.
- Segmented control for `Scan Prescription` and `Add Manually`.
- Large camera/preview area.
- Guidance text for prescription framing.
- Full-width capture CTA.

### 10.3 Profile

- Profile identity card.
- Subscription badge.
- Adherence/streak summary.
- Four summary metric cards.
- Settings rows below the metrics.
- Persistent bottom navigation.

### 10.4 Medication Reminder

- Calm full-screen reminder presentation.
- Clear due time and number of medicines/people.
- Individual medicine cards with person, medicine, strength/form and quantity.
- Strong primary action to mark medicines as taken.
- Secondary dismiss and snooze actions.
- Close affordance that does not compete with the primary action.

### 10.5 State Consistency

- Green indicates the approved/primary product state.
- Orange/terracotta remains available for brand emphasis and legacy/exploratory references.
- Completed, due, missed and disabled states must remain visually distinct without relying on color alone.

---

## 11. Marketing Website Information Architecture

Recommended page structure:

1. **Header** — logo, navigation and primary waitlist CTA.
2. **Hero** — value proposition, product preview and CTA.
3. **Problem** — common prescription and adherence friction.
4. **Solution** — AI prescription reader, grouped reminders and family sharing.
5. **Product demo** — Add Prescription, Analysing, Review and My Reminders.
6. **How it works** — Snap, Understand, Act.
7. **Caregiver support** — multi-person and family use.
8. **Privacy and trust** — data boundaries and user control.
9. **Plans preview** — free and future premium positioning, clearly marked as subject to change.
10. **Waitlist/referral** — conversion module.
11. **FAQ**.
12. **Footer** — contact, privacy, terms and social links.

The product mock-ups and design references must be presented as product previews, not as claims that every shown capability is already available in production.

---

## 12. Frontend Functional Requirements

### 12.1 Framework

- Next.js App Router with TypeScript.
- Tailwind CSS and centrally defined CSS variables/design tokens.
- shadcn/ui primitives customised to RMP rather than used with default appearance.
- Framer Motion only where motion meaningfully improves comprehension.
- Next Image or equivalent optimisation for product screens and brand illustrations.

### 12.2 Responsive Behaviour

Required testing widths:

- 320 px;
- 375 px;
- 390 px;
- 430 px;
- 768 px;
- 1024 px;
- 1280 px; and
- 1440 px or wider.

Requirements:

- No horizontal page scrolling.
- Primary CTAs remain visible and usable.
- Phone mock-ups scale without clipping critical UI.
- Multi-column sections stack in a deliberate order.
- Text does not overlap the watermark or decorative artwork.
- Navigation supports both keyboard and touch interaction.

### 12.3 Accessibility

- Semantic heading hierarchy.
- Keyboard-accessible navigation, dialog and segmented controls.
- Visible focus indicators.
- Programmatically associated labels and errors.
- Screen-reader announcement for form status and copied links.
- Color is not the only state indicator.
- Reduced-motion support.
- Decorative images use empty alt text.
- Product screenshots use useful alt text or nearby captions.

### 12.4 Waitlist Form

Required fields:

- `name`;
- `email`;
- consent acknowledgement; and
- optional referral code from URL state.

The form must not collect medication or health-condition details.

### 12.5 Visual Verification

After meaningful UI changes, use Playwright and Chrome DevTools to verify:

- desktop and mobile rendering;
- console errors;
- failed network requests;
- keyboard flow;
- form states;
- reduced motion;
- image loading;
- Core Web Vitals; and
- screenshot comparison against approved Figma references.

---

## 13. Waitlist and Referral Backend

### 13.1 Entry Creation

The frontend sends a POST request containing:

```json
{
  "name": "Alice",
  "email": "alice@example.com",
  "referralCode": "rmp-bob-12",
  "consent": true
}
```

The API must:

- trim and validate name;
- normalise email to lowercase;
- validate email format;
- reject disposable or malformed requests only according to a documented policy;
- check for an existing email;
- validate an optional referral code;
- generate a unique referral code;
- store consent timestamp and policy version;
- return a stable, non-sensitive response; and
- avoid exposing whether unrelated email addresses exist through enumeration-prone endpoints.

### 13.2 Referral Rules

- One waitlist account maps to one referral code.
- A user cannot refer their own email.
- A referral is counted only after the referred entry is successfully created.
- Duplicate registrations do not increase referral count.
- Rank calculation must be deterministic and documented.
- Any future reward thresholds must be configurable rather than hard-coded in the frontend.

### 13.3 Rate Limiting and Abuse Protection

- Default submission limit: 5 attempts per IP per minute.
- Add a hidden honeypot field.
- Add server-side request-size limits.
- Record minimal security logs without storing prescription or health data.
- Introduce CAPTCHA only when abuse justifies the privacy and conversion trade-off.

---

## 14. Security and Privacy Requirements

### 14.1 Data Minimisation

The waitlist must not collect:

- diagnoses;
- medicine names;
- prescription images;
- insurance details;
- government identifiers;
- medical-record numbers; or
- phone numbers unless a later approved requirement is added.

### 14.2 Security Controls

- TLS for all public traffic.
- Database inaccessible directly from the public internet.
- Secrets provided through environment variables or a secret store.
- No production secrets committed to Git.
- Parameterised database queries.
- CSRF protection where cookie-based state is used.
- Secure headers through Caddy/Nginx and application middleware.
- Restricted CORS.
- Admin export protected by stronger authentication than a token embedded in a URL.
- Regular backups and tested restore procedure.
- Dependency and container-image scanning before deployment.

### 14.3 Privacy Communication

The website must explain in plain language:

- what the waitlist collects;
- why it is collected;
- how long it is kept;
- whether data is shared;
- how users can request deletion; and
- that the public website does not process real prescriptions during the waitlist stage.

---

## 15. Non-Functional Requirements

| Category | Requirement |
| --- | --- |
| **Performance** | Lighthouse Performance target of 95 or higher on representative mobile and desktop runs |
| **Accessibility** | WCAG 2.1 AA minimum; target WCAG 2.2 AA where practical |
| **Availability** | Health-check endpoint and container restart policy |
| **Scalability** | Support initial waitlist traffic on a small VPS without architectural rewrite |
| **Observability** | Structured logs, request IDs, error-rate monitoring and uptime checks |
| **Maintainability** | Typed frontend, documented API, migrations and automated tests |
| **Privacy** | No third-party advertising trackers or unnecessary cookies |
| **Internationalisation** | Copy architecture ready for later localisation |

---

## 16. System Architecture

```text
Browser
  |
  | HTTPS
  v
Caddy / Nginx
  |--------------------------|
  v                          v
Next.js frontend         Go API
                              |
                              v
                         PostgreSQL
```

Deployment boundaries:

- Caddy/Nginx terminates TLS and proxies traffic.
- The frontend and API run as separate containers or clearly separated services.
- PostgreSQL is reachable only on the internal Docker network.
- The public marketing site does not directly access the database.
- Static brand and reference assets are served through the frontend.

---

## 17. API Requirements

| Method | Endpoint | Purpose | Authentication |
| --- | --- | --- | --- |
| `POST` | `/api/waitlist` | Create or safely resolve a waitlist entry | Public, rate-limited |
| `GET` | `/api/waitlist/status` | Retrieve status using a non-guessable status token | Status token |
| `POST` | `/api/waitlist/referral/resolve` | Validate a referral code without exposing personal data | Public, rate-limited |
| `GET` | `/api/admin/export` | Export authorised waitlist data | Admin authentication |
| `GET` | `/health` | Liveness/readiness response | Internal/public minimal response |

### 17.1 Example Success Response

```json
{
  "status": "created",
  "rank": 213,
  "referralCount": 0,
  "referralCode": "rmp-alice-7k2m",
  "referralUrl": "https://remembermypill.com/?ref=rmp-alice-7k2m",
  "statusToken": "opaque-non-guessable-token"
}
```

The status endpoint should use the opaque status token rather than a raw email query parameter.

---

## 18. Data Model

### 18.1 Table: `waitlist_entries`

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `id` | UUID | Yes | Primary key |
| `name` | VARCHAR(100) | Yes | Trimmed and validated |
| `email_normalized` | VARCHAR(255) | Yes | Unique, lowercase |
| `referral_code` | VARCHAR(50) | Yes | Unique public code |
| `referred_by_id` | UUID | No | Foreign key to referring entry |
| `referral_count` | INTEGER | Yes | Default `0`; may also be derived |
| `status_token_hash` | TEXT | Yes | Hash of opaque status token |
| `consent_version` | VARCHAR(30) | Yes | Privacy/consent version accepted |
| `consented_at` | TIMESTAMPTZ | Yes | Consent timestamp |
| `created_at` | TIMESTAMPTZ | Yes | Server-generated timestamp |
| `updated_at` | TIMESTAMPTZ | Yes | Server-generated timestamp |

### 18.2 Table: `referral_events`

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `id` | UUID | Yes | Primary key |
| `referrer_id` | UUID | Yes | Referring waitlist entry |
| `referred_entry_id` | UUID | Yes | Newly created entry |
| `created_at` | TIMESTAMPTZ | Yes | Event timestamp |

A unique constraint on `referred_entry_id` prevents duplicate referral credit.

---

## 19. AI and Model Requirement

AI parsing is represented in the product showcase but is not required to run inside the waitlist backend.

Current implementation placeholders:

- **Primary parsing model:** Gemini 2.5 Flash-Lite.
- **Fallback model:** Claude Haiku 4.5 when confidence or review rules require a second pass.

Before implementation, the engineering team must verify current model availability, pricing, data-processing terms, regional availability and medical-safety limitations. No model output should be presented as medical advice. Prescription interpretation requires clear review and confirmation before any schedule is saved.

---

## 20. Repository Structure

```text
remembermypill-platform/
├── AGENTS.md
├── .github/
│   └── workflows/
│       ├── ci.yml
│       └── deploy.yml
├── backend/
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── api/
│   │   ├── config/
│   │   ├── database/
│   │   ├── middleware/
│   │   ├── referrals/
│   │   └── waitlist/
│   ├── migrations/
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── app/
│   ├── components/
│   │   ├── brand/
│   │   ├── marketing/
│   │   ├── product-demo/
│   │   ├── waitlist/
│   │   └── ui/
│   ├── lib/
│   ├── public/
│   │   ├── brand/
│   │   │   ├── logos/
│   │   │   ├── mascot/
│   │   │   ├── icons/
│   │   │   └── illustrations/
│   │   └── product-screens/
│   ├── styles/
│   │   └── tokens.css
│   ├── tests/
│   └── package.json
├── docs/
│   ├── product/
│   │   └── rmp-prd.md
│   ├── brand-export-manifest.md
│   └── design/
│       └── references/
├── .env.example
├── docker-compose.yml
└── README.md
```

---

## 21. Brand Asset Export Manifest

Before frontend implementation, export and document the following from Figma:

| Asset | Preferred format | Notes |
| --- | --- | --- |
| Horizontal logo | SVG | Primary website header/footer |
| Stacked logo | SVG | Narrow or campaign layouts |
| R monogram | SVG | Favicon, watermark and compact contexts |
| Mascot poses | SVG, PNG fallback | Preserve official artwork and names |
| Product illustrations | SVG/WebP | Prescription, pill, calendar, family and reminder scenes |
| App screenshots | WebP/PNG | Export at 2× where raster is required |
| Color variables | CSS/JSON | Map Figma variable names to web tokens |
| Typography styles | Documentation/CSS | Family, size, weight, line height and tracking |
| Social preview | PNG/JPEG | 1200 × 630 target |

Each exported asset should record:

- Figma node/component name;
- export date;
- variant;
- intended use;
- light/dark background compatibility; and
- source node URL.

---

## 22. Environment Variables

Use `.env.example` with placeholders only:

```dotenv
PORT=8080
DATABASE_URL=postgres://rmp_user:CHANGE_ME@db:5432/remembermypill?sslmode=disable
ADMIN_EXPORT_TOKEN=CHANGE_ME_WITH_A_LONG_RANDOM_SECRET
APP_BASE_URL=https://remembermypill.com
NEXT_PUBLIC_API_URL=https://remembermypill.com/api
LOG_LEVEL=info
```

Production values must be stored outside source control. The production database password and admin credential must never appear in the PRD, repository or deployment logs.

---

## 23. Docker and Deployment Requirements

Minimum production services:

- reverse proxy;
- frontend;
- backend;
- PostgreSQL; and
- backup process.

Key requirements:

- Do not publish PostgreSQL port `5432` to the public host interface.
- Use internal Docker networking.
- Add health checks.
- Use restart policies.
- Pin major image versions and review updates.
- Persist database data in an encrypted or access-controlled volume.
- Run containers as non-root where practical.
- Back up the database on a defined schedule and test restoration.

Illustrative structure:

```yaml
services:
  proxy:
    image: caddy:2-alpine
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    depends_on:
      - frontend
      - backend

  frontend:
    build: ./frontend
    restart: unless-stopped
    expose:
      - "3000"

  backend:
    build: ./backend
    restart: unless-stopped
    expose:
      - "8080"
    env_file:
      - .env
    depends_on:
      db:
        condition: service_healthy

  db:
    image: postgres:16-alpine
    restart: unless-stopped
    env_file:
      - .env
    volumes:
      - rmp_db_volume:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U rmp_user -d remembermypill"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  rmp_db_volume:
```

The final Compose file must include actual proxy routes, production build configuration and secret handling appropriate to the deployment environment.

---

## 24. Testing and Quality Assurance

### 24.1 Frontend

- Unit tests for validation and rank-display logic.
- Component tests for waitlist states.
- Playwright tests for the full conversion flow.
- Mobile and desktop visual regression snapshots.
- Keyboard and screen-reader checks.
- Lighthouse testing on production-like builds.

### 24.2 Backend

- Validation tests.
- Duplicate-entry tests.
- Referral-credit tests.
- Self-referral prevention tests.
- Rate-limit tests.
- Database migration tests.
- Admin export authentication tests.
- Health-check tests.

### 24.3 Brand QA

- Compare implementation with the approved Figma frames.
- Confirm official logo and mascot assets are used.
- Confirm colors come from design tokens rather than scattered hard-coded values.
- Confirm product screenshots are current and not distorted.
- Confirm orange, green and coral roles are consistent.
- Confirm reduced-motion behaviour.

---

## 25. Deliverables

1. Production-ready Next.js marketing website.
2. Approved RMP brand assets exported from Figma.
3. Reusable design-token implementation.
4. Product-demo and mobile-screen showcase.
5. Accessible waitlist and referral experience.
6. Go API with PostgreSQL migrations.
7. Docker Compose production configuration.
8. Automated test suite.
9. Privacy, terms and waitlist-consent copy placeholders for legal review.
10. Deployment and rollback documentation.

---

## 26. Acceptance Criteria

### 26.1 Brand and Visual

- The official Figma logo is used in the header, footer, favicon and social preview as appropriate.
- The official mascot is used only in approved contexts and is not redrawn.
- Design tokens are based on Figma exports.
- The site follows the warm cream, terracotta and deep-green identity shown in the references.
- The interface is light-first and does not revert to the previous dark glowing-gradient concept.
- Product mock-ups preserve aspect ratio and remain legible.
- Figma comparison reveals no critical differences in layout, typography, color or asset usage.

### 26.2 Functional

- A new user can complete the waitlist form and receive a unique referral URL.
- A duplicate email receives a clear safe response without creating another row.
- A valid referral increments credit exactly once.
- An invalid referral does not block legitimate registration.
- Rank and referral count remain consistent across refreshes.
- Admin export rejects unauthorised requests.

### 26.3 Responsive and Accessible

- All required viewport widths pass without horizontal overflow or clipped controls.
- The complete waitlist flow is usable by keyboard.
- Dialog focus is trapped and returned correctly.
- Focus states are visible.
- Errors are announced to assistive technologies.
- Color is not the only indicator of medication or waitlist state.
- Reduced-motion mode removes non-essential transitions.

### 26.4 Performance and Reliability

- Representative Lighthouse scores meet the targets.
- No uncaught console errors on supported pages.
- No failed required network requests.
- Images are compressed and correctly sized.
- Health checks report service readiness.
- Database backup and restore steps are documented and tested before launch.

---

## 27. Risks and Mitigation

| Risk | Mitigation |
| --- | --- |
| Brand inconsistency between deck, app and website | Treat the latest approved Figma components and variables as authoritative; maintain an export manifest |
| Generic AI-generated visual style | Require review against Figma and use official assets rather than approximations |
| Too many competing accent colors | Enforce semantic color roles through design tokens |
| Waitlist spam | Rate limiting, honeypot, request limits and monitoring |
| Email enumeration | Use opaque status tokens and safe duplicate responses |
| Secret leakage | `.env.example` placeholders, secret management and automated scanning |
| Misleading product claims | Label previews and future features clearly; require product/legal review |
| Accessibility regression | Automated checks plus keyboard and manual review |
| Heavy animations reduce performance | Use restrained motion, lazy loading and reduced-motion support |
| Outdated model/provider assumptions | Verify AI model availability and data terms before implementation |

---

## 28. Implementation Sequence

### Phase 1 — Brand Foundation

- Review and approve the Figma file.
- Export logo, mascot, icons and illustrations.
- Create the brand export manifest.
- Implement CSS/design tokens.
- Create reusable logo, button, card and section primitives.

### Phase 2 — Marketing Experience

- Build header, hero and solution sections.
- Build product mock-up/gallery.
- Build privacy, FAQ and footer sections.
- Verify responsive behaviour and visual fidelity.

### Phase 3 — Waitlist

- Implement form states and consent.
- Build Go API and migrations.
- Add duplicate handling, status tokens and referral rules.
- Add rate limiting and abuse protection.

### Phase 4 — QA and Deployment

- Add automated frontend and backend tests.
- Run Playwright, Chrome DevTools and Lighthouse checks.
- Complete accessibility and privacy review.
- Deploy to staging.
- Test backups, rollback and production configuration.

---

## 29. Design Reference Appendix

The accompanying design-reference bundle contains the supplied visual references under `docs/design/references/`.

| File | Reference purpose |
| --- | --- |
| `mobile-home-orange.png` | Earlier orange-accent home exploration |
| `mobile-add-medication-orange.png` | Earlier orange-accent scan flow |
| `mobile-profile-orange.png` | Earlier orange-accent profile screen |
| `mobile-reminder-orange.png` | Earlier orange-accent reminder screen |
| `mobile-home-green-state-1.png` | Deep-green home screen and completed-state reference |
| `mobile-home-green-state-2.png` | Deep-green home screen and disabled/completed control reference |
| `mobile-add-medication-green.png` | Deep-green prescription capture reference |
| `mobile-profile-green.png` | Deep-green profile and metric-card reference |
| `mobile-reminder-green.png` | Deep-green medication reminder reference |
| `deck-data-sources.png` | Marketing deck background, logo and information-layout reference |
| `deck-our-solution.png` | Marketing illustration and phone-mock-up reference |
| `deck-product-demo.png` | Multi-screen product showcase reference |
| `deck-business-model.png` | Marketing card, typography and icon reference |

These screenshots support interpretation, but the latest approved Figma file remains the final authority for production assets and exact design values.
