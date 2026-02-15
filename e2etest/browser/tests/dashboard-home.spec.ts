import { expect, test } from "@playwright/test";
import { loadConfig } from "../helpers/config";

const cfg = loadConfig();

test("home shows quick start for empty apps and overview for active apps", async ({
  page,
}) => {
  await page.goto(cfg.dashboard_url + "/");
  const mode = await Promise.race([
    page
      .getByTestId("home-quick-start")
      .waitFor({ state: "visible", timeout: 10_000 })
      .then(() => "quick-start"),
    page
      .getByTestId("home-overview")
      .waitFor({ state: "visible", timeout: 10_000 })
      .then(() => "overview"),
  ]);
  if (mode === "quick-start") {
    await expect(page.getByTestId("home-quick-start")).toBeVisible({ timeout: 10_000 });
    await expect(page.getByText("Add or manage plans")).toBeVisible();
    return;
  }

  await expect(page.getByTestId("home-overview")).toBeVisible({ timeout: 10_000 });
  await expect(page.getByTestId("home-stat-users")).toBeVisible();
  await expect(page.getByTestId("home-stat-plans")).toBeVisible();
  await expect(page.getByTestId("home-stat-subscriptions")).toBeVisible();
  await expect(page.getByTestId("home-recent-users")).toBeVisible();
  await expect(page.getByTestId("home-recent-subscriptions")).toBeVisible();
});
