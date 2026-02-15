import { expect, test } from "@playwright/test";
import { createSecret } from "../helpers/secret-actions";
import { goToDeveloper } from "../helpers/dashboard-nav";
import { apiPost, controlGet, controlPost, dashPost, dashPostRoot } from "../helpers/api";
import { loadConfig } from "../helpers/config";

const cfg = loadConfig();

async function createPlan(name: string, cents: number): Promise<string> {
  const plan = (await dashPost(cfg, "plans")) as any;
  await dashPost(cfg, `plans/${plan.id}`, { name, currency: "usd" });
  await dashPost(cfg, `plan-items/${plan.items[0].id}`, { unit_amount: cents });
  await dashPost(cfg, `plans/${plan.id}/publish`);
  return plan.id as string;
}

async function waitForCount(kind: string, want: number): Promise<any> {
  for (let i = 0; i < 40; i += 1) {
    const data = (await controlGet(cfg, `/e2e/discord/messages?kind=${kind}`)) as any;
    if ((data.count ?? 0) >= want) return data;
    await new Promise((r) => setTimeout(r, 100));
  }
  throw new Error(`discord ${kind} did not reach count ${want}`);
}

test("discord notifications fire for account signup and billing flows", async ({ page }) => {
  const suffix = Date.now();
  await controlPost(cfg, "/e2e/discord/reset", {});
  await controlPost(cfg, "/e2e/signup/prepare", {});
  await dashPostRoot(cfg, "apps", { name: `Signup Project ${suffix}` });

  await page.goto(cfg.dashboard_url + "/developer");
  await goToDeveloper(page);
  const key = await createSecret(page);
  const apiUser = (await apiPost(cfg, key, "/v1/users", {
    name: "API Notify",
    email: `api-${suffix}@test.dev`,
  })) as any;

  const starterID = await createPlan(`Notify Starter ${suffix}`, 2900);
  const proID = await createPlan(`Notify Pro ${suffix}`, 6900);
  await apiPost(cfg, key, `/v1/users/${apiUser.id}/billing-address`, {
    line1: "1 Test Blvd", city: "San Francisco", postal_code: "94105", country: "US",
  });
  const session = (await apiPost(cfg, key, "/v1/checkout-sessions", {
    plan_id: starterID,
    user_id: apiUser.id,
    redirect_url: cfg.dashboard_url + "/subscriptions",
    cancel_url: cfg.dashboard_url + "/plans",
  })) as any;
  await controlPost(cfg, "/e2e/resolve-checkout", {
    external_session_id: session.external_session_id,
  });
  await controlPost(cfg, "/e2e/upgrade-subscription", {
    user_id: apiUser.id,
    plan_id: proID,
  });

  const mode = (await controlGet(cfg, "/e2e/discord/mode")) as any;
  if (mode.mode === "live") {
    expect(mode.mode).toBe("live");
    return;
  }
  const signup = await waitForCount("signup", 1);
  const billing = await waitForCount("billing", 1);
  expect((signup.messages as string[]).join(" ")).toContain("account signed up");
  expect((billing.messages as string[]).join(" ")).toContain("upgraded");
});
