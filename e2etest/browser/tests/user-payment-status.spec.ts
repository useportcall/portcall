import { expect, test } from "@playwright/test";
import { createSecret } from "../helpers/secret-actions";
import { goToDeveloper } from "../helpers/dashboard-nav";
import { apiPost, controlPost, dashPost } from "../helpers/api";
import { loadConfig } from "../helpers/config";

const cfg = loadConfig();

async function createPlan(name: string, cents: number): Promise<string> {
  const plan = (await dashPost(cfg, "plans")) as any;
  await dashPost(cfg, `plans/${plan.id}`, { name, currency: "usd" });
  await dashPost(cfg, `plan-items/${plan.items[0].id}`, { unit_amount: cents });
  await dashPost(cfg, `plans/${plan.id}/publish`);
  return plan.id as string;
}

test("user view shows payment health state for past due invoices", async ({ page }) => {
  const suffix = Date.now();

  await page.goto(cfg.dashboard_url + "/developer");
  await goToDeveloper(page);
  const key = await createSecret(page);
  const starterPlanID = await createPlan(`Past Due Plan ${suffix}`, 2900);

  const user = (await apiPost(cfg, key, "/v1/users", {
    name: "Past Due User",
    email: `past-due-${suffix}@test.dev`,
  })) as any;
  await apiPost(cfg, key, `/v1/users/${user.id}/billing-address`, {
    line1: "1 Test Blvd", city: "San Francisco", postal_code: "94105", country: "US",
  });
  const session = (await apiPost(cfg, key, "/v1/checkout-sessions", {
    plan_id: starterPlanID,
    user_id: user.id,
    redirect_url: cfg.dashboard_url + "/users",
    cancel_url: cfg.dashboard_url + "/users",
  })) as any;
  await controlPost(cfg, "/e2e/resolve-checkout", { external_session_id: session.external_session_id });
  await controlPost(cfg, "/e2e/user/payment-status", {
    user_id: user.id,
    subscription_status: "past_due",
    invoice_status: "past_due",
  });

  await page.goto(cfg.dashboard_url + `/users/${user.id}`);
  const card = page.getByTestId("payment-health-card");
  await expect(card).toBeVisible();
  await expect(page.getByTestId("payment-health-badge")).toContainText("Action required");
  await card.screenshot({ path: test.info().outputPath("user-payment-status-lowres.png") });
});
