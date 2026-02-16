// Full-page snapshots of the checkout flow: form and success page.
import { test, expect } from "@playwright/test";
import { loadConfig } from "../helpers/config";
import { takeSnapshot } from "../helpers/snapshot";
import {
  seedPlan,
  seedSecret,
  seedUser,
  seedCheckoutSession,
} from "../helpers/api-seed";
import { controlPost } from "../helpers/api";
import {
  fillAndSubmitCheckout,
  waitForCheckoutSuccess,
} from "../helpers/checkout-actions";

const cfg = loadConfig();
let checkoutUrl = "";
let externalSessionId = "";

test.describe.serial("Checkout full-page snapshots", () => {
  test.beforeAll(async () => {
    const key = await seedSecret(cfg);
    const suffix = Date.now();
    const planId = await seedPlan(cfg, `Snapshot Plan ${suffix}`, 4900);
    const userId = await seedUser(cfg, key, `snap-${suffix}@test.dev`);
    const session = await seedCheckoutSession(cfg, key, planId, userId);
    checkoutUrl = session.url;
    externalSessionId = session.externalSessionId;
  });

  test("snapshot checkout form", async ({ page }) => {
    await page.goto(checkoutUrl);
    const is404 = await page
      .getByText("404 page not found")
      .first()
      .isVisible({ timeout: 2_000 })
      .catch(() => false);
    if (is404) {
      test.skip(true, "Checkout frontend unavailable");
    }
    await expect(page.getByTestId("checkout-submit-button")).toBeVisible({
      timeout: 10_000,
    });
    await takeSnapshot(page, cfg, "checkout-form-full");
  });

  test("snapshot checkout success", async ({ page }) => {
    await page.goto(checkoutUrl);
    const is404 = await page
      .getByText("404 page not found")
      .first()
      .isVisible({ timeout: 2_000 })
      .catch(() => false);
    if (is404) {
      test.skip(true, "Checkout frontend unavailable");
    }
    await expect(page.getByTestId("checkout-submit-button")).toBeVisible({
      timeout: 10_000,
    });
    await fillAndSubmitCheckout(page);
    await waitForCheckoutSuccess(page);

    await controlPost(cfg, "/e2e/resolve-checkout", {
      external_session_id: externalSessionId,
    });
    await takeSnapshot(page, cfg, "checkout-success-full");
  });
});
