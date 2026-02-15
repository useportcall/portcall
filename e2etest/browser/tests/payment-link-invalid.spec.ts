import { expect, test } from "@playwright/test";
import { loadConfig } from "../helpers/config";
import { dashPost } from "../helpers/api";
import { seedPlan, seedSecret, seedUser } from "../helpers/api-seed";

const cfg = loadConfig();

test("payment-link with invalid token is rejected by checkout", async ({
  page,
}) => {
  const key = await seedSecret(cfg);
  const userID = await seedUser(cfg, key, `bad-link-${Date.now()}@test.dev`);
  const planID = await seedPlan(cfg, `Bad Link Plan ${Date.now()}`, 1900);

  const link = (await dashPost(cfg, "payment-links", {
    user_id: userID,
    plan_id: planID,
  })) as any;

  const url = new URL(link.url as string);
  url.searchParams.set("pt", "invalid-token");

  await page.goto(url.toString());
  await expect(page.getByText("Invalid checkout link.")).toBeVisible({
    timeout: 10_000,
  });
});
