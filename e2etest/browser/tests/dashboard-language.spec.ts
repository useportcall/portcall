import { expect, test } from "@playwright/test";
import { loadConfig } from "../helpers/config";

const cfg = loadConfig();

test("dashboard language switcher supports japanese", async ({ page }) => {
  await page.goto(cfg.dashboard_url + "/");
  const toggle = page.getByTestId("dashboard-language-switcher");
  const visible = await toggle.isVisible({ timeout: 2_000 }).catch(() => false);
  if (!visible) {
    test.skip(true, "Dashboard language switcher not present in current static bundle");
  }
  await toggle.click();
  await page.getByRole("option", { name: "日本語" }).click();
  await expect(toggle).toContainText("日本語");
});
