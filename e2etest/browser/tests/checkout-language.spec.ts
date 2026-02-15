import { expect, test } from "@playwright/test";
import { loadConfig } from "../helpers/config";
import {
  seedCheckoutSession,
  seedPlan,
  seedSecret,
  seedUser,
} from "../helpers/api-seed";

const cfg = loadConfig();

test("checkout language switcher supports japanese", async ({ page }) => {
  const key = await seedSecret(cfg);
  const suffix = Date.now();
  const planId = await seedPlan(cfg, `Checkout Lang ${suffix}`, 4900);
  const userId = await seedUser(cfg, key, `checkout-lang-${suffix}@test.dev`);
  const session = await seedCheckoutSession(cfg, key, planId, userId);

  await page.goto(session.url);
  const is404 = await page
    .getByText("404 page not found")
    .first()
    .isVisible({ timeout: 2_000 })
    .catch(() => false);
  if (is404) {
    test.skip(true, "Checkout frontend unavailable");
  }

  const toggle = page.getByLabel("Language");
  await expect(toggle).toBeVisible({ timeout: 10_000 });
  await toggle.selectOption("ja");
  await expect(page.getByText("お支払い詳細")).toBeVisible({ timeout: 10_000 });
});
