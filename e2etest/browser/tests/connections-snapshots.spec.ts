// Snapshot: connections page showing the Add Provider dialog with Braintree.
import { test, expect } from "@playwright/test";
import { loadConfig } from "../helpers/config";
import { takeSnapshot } from "../helpers/snapshot";

const cfg = loadConfig();

test.describe.serial("Connection form snapshots", () => {
  test("snapshot connections form", async ({ page }) => {
    await page.goto(cfg.dashboard_url + "/integrations");
    await expect(
      page.getByText("Payment Integrations", { exact: true }),
    ).toBeVisible({ timeout: 10_000 });

    // Open the Add Provider dialog.
    await page.getByTestId("add-provider-button").click();
    await expect(page.getByText("Add Payment Provider")).toBeVisible();

    // Select Braintree to show provider-specific fields.
    await page.getByRole("combobox").click();
    await page
      .locator("[role=option]")
      .filter({ hasText: "Braintree" })
      .click();
    await expect(page.getByText("Manual webhook setup")).toBeVisible();

    await takeSnapshot(page, cfg, "connections-form-braintree");
  });

  test("snapshot connections braintree card", async ({ page }) => {
    await page.goto(cfg.dashboard_url + "/integrations");
    await expect(
      page.getByText("Payment Integrations", { exact: true }),
    ).toBeVisible({ timeout: 10_000 });
    await takeSnapshot(page, cfg, "connections-braintree-card");
  });
});
