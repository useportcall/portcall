// Full-page snapshots of dashboard billing & settings views:
// subscriptions, invoices, developer, integrations, and company.
import { test, expect } from "@playwright/test";
import { loadConfig } from "../helpers/config";
import { takeSnapshot } from "../helpers/snapshot";

const cfg = loadConfig();

test.describe.serial("Dashboard billing & settings snapshots", () => {
  test("snapshot subscriptions page", async ({ page }) => {
    await page.goto(cfg.dashboard_url + "/subscriptions");
    await expect(
      page.getByText("Review and manage subscriptions here."),
    ).toBeVisible({ timeout: 10_000 });
    await takeSnapshot(page, cfg, "dashboard-subscriptions-full");
  });

  test("snapshot invoices page", async ({ page }) => {
    await page.goto(cfg.dashboard_url + "/invoices");
    await expect(
      page.locator("p.font-semibold", { hasText: "Invoices" }),
    ).toBeVisible({ timeout: 10_000 });
    await takeSnapshot(page, cfg, "dashboard-invoices-full");
  });

  test("snapshot developer page", async ({ page }) => {
    await page.goto(cfg.dashboard_url + "/developer");
    await expect(page.getByRole("heading", { name: "Developer" })).toBeVisible({
      timeout: 10_000,
    });
    await takeSnapshot(page, cfg, "dashboard-developer-full");
  });

  test("snapshot integrations page", async ({ page }) => {
    await page.goto(cfg.dashboard_url + "/integrations");
    await expect(
      page.getByText("Payment Integrations", { exact: true }),
    ).toBeVisible({ timeout: 10_000 });
    await takeSnapshot(page, cfg, "dashboard-integrations-full");
  });

  test("snapshot company page", async ({ page }) => {
    await page.goto(cfg.dashboard_url + "/company");
    await expect(
      page.locator("p.font-semibold", { hasText: "Company details" }),
    ).toBeVisible({ timeout: 10_000 });
    await takeSnapshot(page, cfg, "dashboard-company-full");
  });

  test("snapshot quotes list page", async ({ page }) => {
    await page.goto(cfg.dashboard_url + "/quotes");
    await expect(page.getByRole("heading", { name: "Quotes" })).toBeVisible({
      timeout: 10_000,
    });
    await takeSnapshot(page, cfg, "dashboard-quotes-list-full");
  });

  test("snapshot usage page", async ({ page }) => {
    await page.goto(cfg.dashboard_url + "/usage");
    // The usage view is hidden for billing-exempt (dogfood) apps.
    // If the Billing heading appears, snapshot it; otherwise skip.
    const visible = await page
      .getByRole("heading", { name: "Billing" })
      .isVisible({ timeout: 5_000 })
      .catch(() => false);
    if (!visible) {
      test.skip(true, "Usage page hidden for billing-exempt apps");
    }
    await takeSnapshot(page, cfg, "dashboard-usage-full");
  });
});
