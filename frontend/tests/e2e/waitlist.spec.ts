import { expect, test } from "@playwright/test";

const viewports = [
  { name: "320", width: 320, height: 720 },
  { name: "375", width: 375, height: 812 },
  { name: "390", width: 390, height: 844 },
  { name: "430", width: 430, height: 932 },
  { name: "768", width: 768, height: 1024 },
  { name: "1024", width: 1024, height: 900 },
  { name: "1440", width: 1440, height: 1000 },
];

const oneViewportTargets = [
  { name: "1366x768", width: 1366, height: 768 },
  { name: "1440x900", width: 1440, height: 900 },
  { name: "1920x1080", width: 1920, height: 1080 },
];

test("waitlist page is responsive with no horizontal overflow", async ({
  page,
}) => {
  await page.emulateMedia({ reducedMotion: "reduce" });
  await page.goto("/waitlist");

  for (const viewport of viewports) {
    await page.setViewportSize(viewport);
    await expect(
      page.getByRole("heading", {
        level: 1,
        name: /be first to know when remember my pill is ready/i,
      }),
    ).toBeVisible();
    await expect(page.getByLabel("Email")).toBeVisible();
    await expect(
      page.getByRole("button", { name: "Join the waitlist" }),
    ).toBeVisible();
    expect(
      await page.locator("body").evaluate((node) => node.scrollWidth),
    ).toBe(await page.locator("body").evaluate((node) => node.clientWidth));
  }
});

test("waitlist page fits within one viewport on common desktop sizes", async ({
  page,
}) => {
  await page.goto("/waitlist");
  for (const target of oneViewportTargets) {
    await page.setViewportSize(target);
    const { scrollHeight, innerHeight } = await page.evaluate(() => ({
      scrollHeight: document.body.scrollHeight,
      innerHeight: window.innerHeight,
    }));
    expect(scrollHeight).toBeLessThanOrEqual(innerHeight);
  }
});

test("waitlist page has only an email field, validates it, and supports keyboard submission", async ({
  page,
}) => {
  await page.goto("/waitlist");

  await expect(page.getByLabel("Email")).toHaveAttribute("type", "email");
  await expect(page.getByLabel("Email")).toHaveAttribute(
    "autocomplete",
    "email",
  );
  await expect(page.getByLabel("Name")).toHaveCount(0);

  // Invalid email is rejected inline without a network call.
  await page.getByLabel("Email").fill("not-an-email");
  await page.getByRole("button", { name: "Join the waitlist" }).click();
  await expect(page.getByText(/enter a valid email address/i)).toBeVisible();

  // A valid email submits via the keyboard; the local dev backend isn't
  // running for this suite, so this also exercises the network-error path.
  await page.getByLabel("Email").fill("ada@example.com");
  await page.keyboard.press("Enter");
  await expect(
    page.getByText(/couldn't reach the waitlist service/i),
  ).toBeVisible();
});

test("back to home link returns to the landing page", async ({ page }) => {
  await page.goto("/waitlist");
  await page.getByRole("link", { name: /back to home/i }).click();
  await expect(page).toHaveURL(/\/$/);
  await expect(
    page.getByRole("heading", {
      level: 1,
      name: /turn confusing prescriptions into a routine you can trust/i,
    }),
  ).toBeVisible();
});
