// Checkout page actions: fill mock card details and submit.
import { Page, expect } from "@playwright/test";

/** Fill the mock checkout form and submit. */
export async function fillAndSubmitCheckout(page: Page) {
  // Card details
  await page.getByTestId("checkout-card-input").fill("4242 4242 4242 4242");
  await page.getByTestId("checkout-expiry-input").fill("12/30");
  await page.getByTestId("checkout-cvc-input").fill("123");

  // Billing details — country (button-based combobox)
  await page.getByRole("button", { name: "Country/territory" }).click();
  await page.getByPlaceholder("Search country").fill("United States");
  await page.getByText("United States").first().click();

  // Address
  await page.getByPlaceholder("1 Oxford Street").fill("1 Test Blvd");
  await page.getByPlaceholder("WC1A 1NU").fill("94105");
  await page.getByPlaceholder("London").fill("San Francisco");

  // Submit
  await page.getByTestId("checkout-submit-button").click();
}

export async function autofillAndSubmitMockCheckout(page: Page) {
  const autofillButton = page.getByTestId("checkout-autofill-button");
  await expect(autofillButton).toBeVisible({ timeout: 10_000 });
  await autofillButton.click();
  await expect(page.getByTestId("checkout-card-input")).toHaveValue(
    "4242 4242 4242 4242",
  );
  await expect(page.getByTestId("checkout-expiry-input")).toHaveValue("12/30");
  await expect(page.getByTestId("checkout-cvc-input")).toHaveValue("123");
  await page.getByTestId("checkout-submit-button").click();
}

/** Wait for checkout payment success state. */
export async function waitForCheckoutSuccess(page: Page) {
  await expect(page.getByText("Payment successful!")).toBeVisible({
    timeout: 10_000,
  });
}
