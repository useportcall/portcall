// Full-page snapshots of the invoice view (light + dark mode).
import { test, expect } from "@playwright/test";
import { loadConfig } from "../helpers/config";
import { takeSnapshot } from "../helpers/snapshot";
import { seedInvoice } from "../helpers/api-seed";

const cfg = loadConfig();
let invoiceId = "";

test.describe.serial("Invoice page snapshots", () => {
  test("seed invoice for snapshot", async () => {
    invoiceId = await seedInvoice(cfg);
    expect(invoiceId).toBeTruthy();
  });

  test("snapshot invoice light mode", async ({ page }) => {
    await page.goto(`${cfg.file_url}/invoices/${invoiceId}/view`);
    const frame = page.frameLocator("#invoice-frame");
    await expect(frame.locator(".invoice-container")).toBeVisible({
      timeout: 10_000,
    });
    await takeSnapshot(page, cfg, "invoice-view-light");
  });

  test("snapshot invoice dark mode", async ({ page }) => {
    await page.goto(`${cfg.file_url}/invoices/${invoiceId}/view`);
    const frame = page.frameLocator("#invoice-frame");
    await expect(frame.locator(".invoice-container")).toBeVisible({
      timeout: 10_000,
    });
    // Toggle dark mode via the theme button
    await page.click("#theme-btn");
    // Wait for theme to propagate to iframe
    await page.waitForTimeout(500);
    await takeSnapshot(page, cfg, "invoice-view-dark");
  });
});
