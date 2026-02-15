// Full-page snapshots of core dashboard views:
// home, plans, plan editor, users, and user detail.
import { test, expect } from "@playwright/test";
import { loadConfig } from "../helpers/config";
import { takeSnapshot } from "../helpers/snapshot";
import { seedPlan, seedSecret, seedUser } from "../helpers/api-seed";

const cfg = loadConfig();
let planId = "";
let userId = "";

test.describe.serial("Dashboard core page snapshots", () => {
  test("seed data for dashboard views", async () => {
    const key = await seedSecret(cfg);
    const suffix = Date.now();
    planId = await seedPlan(cfg, `Dashboard Plan ${suffix}`, 2900);
    userId = await seedUser(cfg, key, `dash-snap-${suffix}@test.dev`);
  });

  test("snapshot home page", async ({ page }) => {
    await page.goto(cfg.dashboard_url + "/");
    const loaded = await Promise.race([
      page
        .getByTestId("home-quick-start")
        .waitFor({ state: "visible", timeout: 10_000 })
        .then(() => true),
      page
        .getByTestId("home-overview")
        .waitFor({ state: "visible", timeout: 10_000 })
        .then(() => true),
    ]);
    expect(loaded).toBe(true);
    await takeSnapshot(page, cfg, "dashboard-home-full");
  });

  test("snapshot plans list", async ({ page }) => {
    await page.goto(cfg.dashboard_url + "/plans");
    await expect(page.getByRole("heading", { name: "Plans" })).toBeVisible({
      timeout: 10_000,
    });
    await takeSnapshot(page, cfg, "dashboard-plans-list-full");
  });

  test("snapshot plan editor", async ({ page }) => {
    await page.goto(cfg.dashboard_url + `/plans/${planId}`);
    await expect(page.getByTestId("plan-price-input")).toBeVisible({
      timeout: 10_000,
    });
    await takeSnapshot(page, cfg, "dashboard-plan-editor-full");
  });

  test("snapshot users list", async ({ page }) => {
    await page.goto(cfg.dashboard_url + "/users");
    await expect(page.getByTestId("add-user-button")).toBeVisible({
      timeout: 10_000,
    });
    await takeSnapshot(page, cfg, "dashboard-users-list-full");
  });

  test("snapshot user detail", async ({ page }) => {
    await page.goto(cfg.dashboard_url + `/users/${userId}`);
    await expect(
      page.getByRole("heading", { name: "User Information" }),
    ).toBeVisible({ timeout: 10_000 });
    await takeSnapshot(page, cfg, "dashboard-user-detail-full");
  });
});
