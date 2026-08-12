# Brand export manifest

**Export date:** 16 July 2026 (Phase 1); updated 11 August 2026 (full marketing redesign).
**Status:** Assets used by the current marketing site are recorded here.

| Asset | Repository path | Figma source | Format | Intended use | Background compatibility |
| --- | --- | --- | --- | --- | --- |
| R/pill mark | `frontend/public/brand/r-pill-mark.svg` | `DP Option 2 + App Icon`, node `6021:115`, nested official R/pill artwork node `6021:128` | SVG exported by Figma MCP | Header/footer logo tile, favicon | Placed on a terracotta rounded-square tile (matching the Figma "Primary Logo"), so it works on any surface |
| App-icon composition | `frontend/public/product-preview/app-icon.png` | `DP Option 2 + App Icon`, node `6021:115` | PNG Figma MCP frame export | OpenGraph/Twitter social preview image | Cream canvas / warm peach-orange icon treatment |
| 2D mascot artwork | `frontend/public/brand/mascot/mascot-wave.png` | `2D Mascot`, node `6023:243` | PNG Figma MCP export (flattened; no transparent cutout is available from the source file) | Caregiver section, waitlist success state, final CTA | Cream background baked into the export; matches the site's cream/cream-soft surfaces |

The Figma file is the source of truth: [RMP Brand Identity](https://www.figma.com/design/LRx7E2phkTHL2DRV1gR86a/RMP-Brand-Identity?node-id=0-1&t=PGo3ZAXibgR7MOOP-1).

## Documented decisions

- **Wordmark:** the Figma file does not contain a standalone wordmark vector (the "Remember My Pill" lettering only exists baked into the flattened `RMP Pallete 1` board, or as white text sized for the dark-gradient social banners). Per PRD §9.5 the wordmark should be an exported asset rather than typed text, but no light-background, terracotta wordmark asset exists to export. The site therefore sets "Remember My Pill" in Manrope Extrabold next to the untouched R/pill logo mark (the actual drawn logomark), matching the letterform, weight, and color shown on the brand board. This should be revisited if/when a dedicated wordmark vector is added to Figma.
- **Supporting icon set** (reminder/schedule/completed/pill/notification/care) is only available flattened inside the `RMP Pallete 1` raster board, not as separate vector layers. Per PRD §9.8, Lucide icons are used for these utility roles instead of cropping the raster board, normalized to a consistent stroke weight and rounded style.
- **Mascot** has no transparent-background export available from Figma (the source node is a flattened raster with a cream backdrop baked in). It is used at sizes/on surfaces where the baked-in cream backdrop reads correctly rather than being recolored or cut out.

Still not exported: a distinct horizontal/stacked wordmark vector, additional mascot poses, and real mobile product screenshots (the current product-preview screens are original illustrations built from the PRD's §10.1–10.4 descriptions, not exported Figma screens — no such screens exist in the Figma file).
