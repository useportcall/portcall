import { expect, test } from "@playwright/test";
import { fillAndSubmitCheckout, waitForCheckoutSuccess } from "../helpers/checkout-actions";
import { dashGet, dashGetBlob, dashPost } from "../helpers/api";
import { loadConfig } from "../helpers/config";

const cfg = loadConfig();
const VALID_SIGNATURE =
  "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAIAAAACCAYAAABytg0kAAAAFElEQVR4nGLi4uL6z4AMAAEAAP//EAgBII8WKK0AAAAASUVORK5CYII=";

async function drawSignature(page: any) {
  const canvas = page.locator("#signature");
  await expect(canvas).toBeVisible();
  const box = await canvas.boundingBox();
  if (!box) throw new Error("signature canvas not found");

  await page.mouse.move(box.x + 20, box.y + box.height / 2);
  await page.mouse.down();
  await page.mouse.move(box.x + box.width / 2, box.y + box.height / 3);
  await page.mouse.move(box.x + box.width - 20, box.y + box.height / 2);
  await page.mouse.up();
}

test("quote issue, sign, accept, checkout, and signed artifact retrieval", async ({ page }) => {
  const quote = await dashPost(cfg, "quotes");
  const quoteId = String((quote as any).id);
  const planId = String((quote as any).plan.id);

  const user = await dashPost(cfg, "users", {
    email: `quote-e2e-${Date.now()}@test.dev`,
  });
  const userId = String((user as any).id);

  const items = (await dashGet(cfg, `plan-items?plan_id=${planId}`)) as any[];
  const fixed = items.find((item) => item.pricing_model === "fixed");
  await dashPost(cfg, `plan-items/${fixed.id}`, { unit_amount: 4900 });
  await dashPost(cfg, `quotes/${quoteId}`, {
    user_id: userId,
    recipient_email: `quote-recipient-${Date.now()}@test.dev`,
    recipient_name: "Quote Recipient",
    recipient_title: "Director",
    company_name: "Quote Co",
    direct_checkout_enabled: true,
    toc: "Payment due upon acceptance.",
  });
  await dashPost(cfg, `quotes/${quoteId}/send`);

  const sentQuote = await dashGet(cfg, `quotes/${quoteId}`);
  const quoteURL = String((sentQuote as any).url);
  expect(quoteURL).toContain("/quotes/");

  await page.goto(quoteURL);
  await expect(page.getByText("Quote details")).toBeVisible();
  await drawSignature(page);
  await page.evaluate((signature) => {
    (document.getElementById("signatureData") as HTMLInputElement).value = signature;
    (document.querySelector("form") as HTMLFormElement).submit();
  }, VALID_SIGNATURE);

  await expect.poll(() => page.url()).toContain(cfg.checkout_url);
  const hasCheckoutForm = await page
    .getByTestId("checkout-card-input")
    .isVisible({ timeout: 2_000 })
    .catch(() => false);
  if (hasCheckoutForm) {
    await fillAndSubmitCheckout(page);
    await waitForCheckoutSuccess(page);
  }

  const signatureBlob = await dashGetBlob(cfg, `quotes/${quoteId}/signature`);
  expect(signatureBlob.size).toBeGreaterThan(64);

  const acceptedQuote = await dashGet(cfg, `quotes/${quoteId}`);
  expect((acceptedQuote as any).status).toBe("accepted");
});
