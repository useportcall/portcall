// Dashboard actions: add a mock/local payment provider.
import { Page, expect } from "@playwright/test";

/** Ensure a local/mock payment provider exists. */
export async function addMockProvider(page: Page) {
  // The base seed already creates a local connection, so check first
  const existing = page.getByText("Local Dev");
  if (await existing.isVisible().catch(() => false)) return;

  await page.getByTestId("add-provider-button").click();
  await expect(page.getByText("Add Payment Provider")).toBeVisible();
  await page.locator("#provider").click();
  await page.getByRole("option", { name: "Mock Provider" }).click();
  await page.getByTestId("submit-provider-button").click();
  await expect(page.getByTestId("add-provider-button")).toBeVisible({
    timeout: 5000,
  });
}
