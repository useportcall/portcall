import { expect, test } from "@playwright/test";
import { dashGet, dashPost } from "../helpers/api";
import { loadConfig } from "../helpers/config";

const cfg = loadConfig();

test("quote view supports japanese toggle", async ({ page }) => {
  const quote = await dashPost(cfg, "quotes");
  const quoteId = String((quote as any).id);

  const user = await dashPost(cfg, "users", {
    email: `quote-lang-${Date.now()}@test.dev`,
  });
  const userId = String((user as any).id);

  await dashPost(cfg, `quotes/${quoteId}`, {
    user_id: userId,
    recipient_email: `quote-lang-recipient-${Date.now()}@test.dev`,
    recipient_name: "Quote Recipient",
    recipient_title: "Director",
    company_name: "Quote Co",
    direct_checkout_enabled: true,
    toc: "Payment due upon acceptance.",
  });
  await dashPost(cfg, `quotes/${quoteId}/send`);

  const sentQuote = await dashGet(cfg, `quotes/${quoteId}`);
  const quoteURL = new URL(String((sentQuote as any).url));
  quoteURL.searchParams.set("lang", "ja");

  await page.goto(quoteURL.toString());
  await expect(page.getByText("見積もり詳細")).toBeVisible({ timeout: 10_000 });
  await page.getByRole("link", { name: /English/ }).click();
  await expect(page.getByText("Quote details")).toBeVisible({ timeout: 10_000 });
});
