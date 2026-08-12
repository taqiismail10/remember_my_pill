import { expect, test } from "@playwright/test";
import { mkdirSync } from "node:fs";
import { resolve } from "node:path";

const viewports = [
  { name: "320", width: 320, height: 720 },
  { name: "375", width: 375, height: 812 },
  { name: "390", width: 390, height: 844 },
  { name: "430", width: 430, height: 932 },
  { name: "768", width: 768, height: 1024 },
  { name: "1024", width: 1024, height: 900 },
  { name: "1280", width: 1280, height: 900 },
  { name: "1440", width: 1440, height: 1000 },
];

test("marketing site is responsive, keyboard-accessible, and preview-honest", async ({
  page,
}) => {
  const evidenceDir = resolve(process.cwd(), "../docs/evidence/redesign");
  mkdirSync(evidenceDir, { recursive: true });

  await page.emulateMedia({ reducedMotion: "reduce" });
  await page.goto("/");

  for (const viewport of viewports) {
    await page.setViewportSize(viewport);
    await expect(
      page.getByRole("heading", {
        level: 1,
        name: /turn confusing prescriptions into a routine you can trust/i,
      }),
    ).toBeVisible();
    await expect(
      page
        .getByRole("banner")
        .getByRole("link", { name: "Remember My Pill home" }),
    ).toBeVisible();
    expect(
      await page.locator("body").evaluate((node) => node.scrollWidth),
    ).toBe(await page.locator("body").evaluate((node) => node.clientWidth));
    await page.screenshot({
      path: resolve(evidenceDir, `marketing-${viewport.name}.png`),
      fullPage: true,
    });
  }

  await page.setViewportSize(viewports[0]);
  await page.keyboard.press("Tab");
  await expect(
    page
      .getByRole("banner")
      .getByRole("link", { name: "Remember My Pill home" }),
  ).toBeFocused();
  await expect(
    page.getByText(/does not collect prescriptions, medication details/i),
  ).toBeVisible();

  await page.setViewportSize(viewports[4]);
  const waitlistSection = page.locator("#waitlist");
  await expect(
    waitlistSection.getByRole("link", { name: "Join the waitlist" }),
  ).toHaveAttribute("href", "/waitlist");
});

test("high-intent CTAs route to the dedicated /waitlist page", async ({
  page,
}) => {
  await page.goto("/");
  for (const cta of await page
    .getByRole("link", { name: "Join the waitlist" })
    .all()) {
    await expect(cta).toHaveAttribute("href", "/waitlist");
  }
  await page
    .getByRole("banner")
    .getByRole("link", { name: "Join the waitlist" })
    .click();
  await expect(page).toHaveURL(/\/waitlist$/);
  await expect(
    page.getByRole("heading", {
      level: 1,
      name: /be first to know when remember my pill is ready/i,
    }),
  ).toBeVisible();
});
