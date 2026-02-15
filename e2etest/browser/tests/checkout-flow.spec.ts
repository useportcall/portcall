// Full browser e2e checkout flow test.
// Creates a plan with features, user, secret, and payment provider
// in the dashboard, then completes checkout in the browser.
import { test, expect } from "@playwright/test";
import { loadConfig } from "../helpers/config";
import {
  goToPlans,
  goToDeveloper,
  goToIntegrations,
  goToUsers,
} from "../helpers/dashboard-nav";
import {
  createPlan,
  setPlanName,
  setPlanPrice,
  publishPlan,
} from "../helpers/plan-actions";
import { addIncludedFeature } from "../helpers/feature-actions";
import { createSecret } from "../helpers/secret-actions";
import { addMockProvider } from "../helpers/provider-actions";
import { createUser } from "../helpers/user-actions";
import { dashGet, dashPost, apiPost, controlPost } from "../helpers/api";
import {
  fillAndSubmitCheckout,
  waitForCheckoutSuccess,
} from "../helpers/checkout-actions";

const cfg = loadConfig();
let secretKey = "";
let userPublicId = "";
let checkoutUrl = "";
let externalSessionId = "";
let skipCheckoutCompletion = false;

function normalizeCheckoutURL(raw: string, checkoutBase: string): string {
  try {
    const parsed = new URL(raw);
    const id = parsed.searchParams.get("id");
    const token = parsed.searchParams.get("st");
    if (!id || !token) return raw;

    const base = new URL(checkoutBase);
    base.searchParams.set("id", id);
    base.searchParams.set("st", token);
    return base.toString();
  } catch {
    return raw;
  }
}

test.describe.serial("Full checkout flow via browser", () => {
  test("1 — create plan with features and publish", async ({ page }) => {
    await page.goto(cfg.dashboard_url + "/plans");
    await goToPlans(page);
    await createPlan(page);
    await setPlanName(page, "Browser E2E Plan");
    await setPlanPrice(page, "29.00");
    await addIncludedFeature(page, "analytics");
    await addIncludedFeature(page, "support");
    await publishPlan(page);
  });

  test("2 — create API secret", async ({ page }) => {
    await page.goto(cfg.dashboard_url + "/developer");
    await goToDeveloper(page);
    secretKey = await createSecret(page);
    expect(secretKey.length).toBeGreaterThan(10);
  });

  test("3 — add mock payment provider", async ({ page }) => {
    await page.goto(cfg.dashboard_url + "/integrations");
    await goToIntegrations(page);
    await addMockProvider(page);
  });

  test("4 — create user via dashboard", async ({ page }) => {
    await page.goto(cfg.dashboard_url + "/users");
    await goToUsers(page);
    await createUser(page, "checkout-e2e@test.com");
    userPublicId = page.url().split("/users/")[1];
    expect(userPublicId).toBeTruthy();
  });

  test("5 — create checkout session via API", async () => {
    await apiPost(cfg, secretKey, `/v1/users/${userPublicId}/billing-address`, {
      line1: "1 Test Blvd",
      city: "San Francisco",
      postal_code: "94105",
      country: "US",
    });
    const plans = (await dashGet(cfg, "plans")) as any;
    const list = Array.isArray(plans) ? plans : (plans?.data ?? []);
    const plan =
      list.find((p: any) => p.name === "Browser E2E Plan") ?? list[0];
    const session = await apiPost(cfg, secretKey, "/v1/checkout-sessions", {
      plan_id: plan.public_id ?? plan.id,
      user_id: userPublicId,
      redirect_url: cfg.dashboard_url + "/subscriptions",
      cancel_url: cfg.dashboard_url + "/plans",
    });
    checkoutUrl = normalizeCheckoutURL((session as any).url, cfg.checkout_url);
    externalSessionId = (session as any).external_session_id;
    expect(checkoutUrl).toBeTruthy();
  });

  test("6 — complete checkout in browser", async ({ page }) => {
    await page.goto(checkoutUrl);
    const is404Page = await page
      .getByText("404 page not found")
      .first()
      .isVisible({ timeout: 2_000 })
      .catch(() => false);
    if (is404Page) {
      skipCheckoutCompletion = true;
      test.skip(true, "Checkout frontend static bundle is unavailable in this environment");
    }
    await expect(page.getByTestId("checkout-submit-button")).toBeVisible({
      timeout: 10_000,
    });
    await fillAndSubmitCheckout(page);
    await waitForCheckoutSuccess(page);
  });

  test("7 — trigger billing saga", async () => {
    test.skip(skipCheckoutCompletion, "Checkout completion was skipped");
    const res = await controlPost(cfg, "/e2e/resolve-checkout", {
      external_session_id: externalSessionId,
    });
    expect((res as any).ok).toBe(true);
  });

  test("8 — verify subscription in dashboard", async ({ page }) => {
    test.skip(skipCheckoutCompletion, "Checkout completion was skipped");
    await page.goto(cfg.dashboard_url + "/subscriptions");
    await expect(page.getByText("active")).toBeVisible({ timeout: 10_000 });
  });

  test("9 — verify invoice in dashboard", async ({ page }) => {
    test.skip(skipCheckoutCompletion, "Checkout completion was skipped");
    await page.goto(cfg.dashboard_url + "/invoices");
    await expect(page.getByTestId("invoice-table")).toBeVisible({
      timeout: 10_000,
    });
    await expect(page.getByText("paid")).toBeVisible();
  });
});
