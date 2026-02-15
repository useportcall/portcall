import { Page, expect } from "@playwright/test";

export async function createPaymentLinkFromUserPage(
  page: Page,
  planId: string,
  expiryDays = "7",
): Promise<string> {
  await page
    .getByRole("button", { name: "Create Payment Link" })
    .first()
    .click();
  await expect(
    page.getByRole("heading", { name: "Create Payment Link" }),
  ).toBeVisible();

  await page.getByTestId("payment-link-plan-select").selectOption(planId);
  await page.getByTestId("payment-link-expiry-days").fill(expiryDays);
  await page.getByTestId("payment-link-create-submit").click();

  const url = await page.getByTestId("payment-link-url-input").inputValue();
  await page.getByTestId("payment-link-done").click();
  return url;
}
