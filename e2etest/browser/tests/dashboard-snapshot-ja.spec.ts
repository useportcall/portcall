import { expect, test } from "@playwright/test";
import { loadConfig } from "../helpers/config";
import { takeSnapshot } from "../helpers/snapshot";

const cfg = loadConfig();

test("snapshot dashboard home in japanese", async ({ page }) => {
  await page.goto(cfg.dashboard_url + "/");
  await page.evaluate(() => localStorage.setItem("i18nextLng", "ja"));
  await page.reload();
  await expect(
    page.getByText("プラン、ユーザー、サブスクリプション、請求をライブデータで管理します。"),
  ).toBeVisible({ timeout: 10_000 });
  await takeSnapshot(page, cfg, "dashboard-home-ja-full");
});
