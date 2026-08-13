import { expect, test } from "@playwright/test";

for (const viewport of [
  { width: 320, height: 720 },
  { width: 1440, height: 1000 },
]) {
  test(`status routes remain safely unavailable at ${viewport.width}px`, async ({
    page,
  }) => {
    await page.setViewportSize(viewport);
    await page.goto("/waitlist/verify#v=short-lived-token");
    await expect(page).toHaveURL(/\/waitlist\/verify$/);
    await expect(page.getByText(/not available yet/i)).toBeVisible();
    await page.goto("/waitlist/status");
    await expect(page.getByText(/not available yet/i)).toBeVisible();
    expect(
      await page.locator("body").evaluate((node) => node.scrollWidth),
    ).toBe(await page.locator("body").evaluate((node) => node.clientWidth));
  });
}
