import { expect, test } from "@playwright/test";
import { loadConfig } from "../helpers/config";
import { seedPlan } from "../helpers/api-seed";
import { controlGet, controlPost } from "../helpers/api";
import { createUser } from "../helpers/user-actions";
import {
  autofillAndSubmitMockCheckout,
  waitForCheckoutSuccess,
} from "../helpers/checkout-actions";
import { createPaymentLinkFromUserPage } from "../helpers/payment-link-actions";
import { uploadArtifactIfLive } from "../helpers/artifact-upload";
import { takeSnapshot } from "../helpers/snapshot";

const cfg = loadConfig();

test("payment-link flow via dashboard creates checkout and subscription", async ({
  browser,
}, testInfo) => {
  const context = await browser.newContext({
    recordVideo: {
      dir: testInfo.outputDir,
      size: { width: 1280, height: 720 },
    },
  });
  const page = await context.newPage();

  const planId = await seedPlan(cfg, `Payment Link Plan ${Date.now()}`, 2900);
  await page.goto(cfg.dashboard_url + "/users");
  await createUser(page, `payment-link-${Date.now()}@test.dev`);
  await takeSnapshot(page, cfg, "payment-link-user-page");

  const paymentLinkURL = await createPaymentLinkFromUserPage(page, planId);
  expect(paymentLinkURL).toContain("pl=");
  expect(paymentLinkURL).toContain("pt=");

  await page.goto(paymentLinkURL);
  const checkoutSessionID = new URL(page.url()).searchParams.get("id") ?? "";
  expect(checkoutSessionID).toContain("cs_");
  const checkoutSession = (await controlGet(
    cfg,
    `/e2e/checkout-session?public_id=${encodeURIComponent(checkoutSessionID)}`,
  )) as any;
  const externalSessionID = checkoutSession.external_session_id as string;
  expect(externalSessionID).toBeTruthy();

  await autofillAndSubmitMockCheckout(page);
  await waitForCheckoutSuccess(page);
  await takeSnapshot(page, cfg, "payment-link-checkout-success");

  const resolved = await controlPost(cfg, "/e2e/resolve-checkout", {
    external_session_id: externalSessionID,
  });
  expect((resolved as any).ok).toBe(true);

  await page.goto(cfg.dashboard_url + "/subscriptions");
  await expect(page.getByText("active").first()).toBeVisible({
    timeout: 10_000,
  });

  await context.close();

  const videoPath = await page.video()?.path();
  if (videoPath) {
    await uploadArtifactIfLive(
      cfg,
      videoPath,
      `payment-link-demo-${Date.now()}.webm`,
      "video/webm",
      "payment-link dashboard demo video",
    );
  }
});
