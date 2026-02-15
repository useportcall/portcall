import { expect, test } from "@playwright/test";
import { addIncludedFeature } from "../helpers/feature-actions";
import { loadConfig } from "../helpers/config";
import { createPlan, setPlanName, setPlanPrice } from "../helpers/plan-actions";
import { goToPlans } from "../helpers/dashboard-nav";
import { recordDashboardMutations } from "../helpers/request-recorder";

const cfg = loadConfig();

test("plan view updates fixed and metered settings through API calls", async ({ page }) => {
  const requests = recordDashboardMutations(page);
  const featureId = `api_calls_${Date.now()}`;

  await page.goto(cfg.dashboard_url + "/plans");
  await goToPlans(page);
  await createPlan(page);
  await setPlanName(page, "Plan View Coverage");
  await setPlanPrice(page, "12.34");
  await addIncludedFeature(page, "priority_support");

  await page.getByTestId("add-metered-feature-button").click();
  const meteredCard = page.getByTestId("metered-feature-card").first();
  await meteredCard.getByTestId("metered-title-input").fill("API Calls");
  await meteredCard.getByTestId("metered-title-input").press("Tab");
  await meteredCard.getByTestId("metered-feature-button").click();
  await page.getByTestId("metered-feature-input").fill(featureId);
  await page.getByTestId("metered-feature-input").press("Enter");

  await meteredCard.getByTestId("metered-limit-button").click();
  await page.getByTestId("metered-limit-input").fill("250");
  await page.getByTestId("metered-limit-input").press("Enter");
  await meteredCard.getByTestId("metered-reset-button").click();
  await page.getByRole("option", { name: "monthly" }).click();
  await meteredCard.getByTestId("metered-rollover-button").click();
  await page.getByRole("option", { name: "Yes" }).click();

  await page.waitForTimeout(300);
  await page.reload();
  await expect(page.getByTestId("plan-price-input")).toHaveValue("12.34");
  await expect(page.getByTestId("feature-badges")).toContainText("priority_support");
  await expect(page.getByTestId("metered-title-input").first()).toHaveValue("API Calls");
  await expect(page.getByTestId("metered-limit-button").first()).toContainText("250");

  expect(requests.hasPath("/plans/")).toBe(true);
  expect(requests.hasPath("/plan-items")).toBe(true);
  expect(requests.hasPath("/plan-items/")).toBe(true);
  expect(requests.hasPath("/features")).toBe(true);
  expect(requests.hasPath("/plan-features/")).toBe(true);
});
