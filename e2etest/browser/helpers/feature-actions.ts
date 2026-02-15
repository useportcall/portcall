// Dashboard actions: add included features to a plan via the browser UI.
import { Page, expect } from "@playwright/test";

/** Add a basic (boolean) feature by name via the combobox. */
export async function addIncludedFeature(page: Page, name: string) {
  await page.getByTestId("add-feature-button").click();
  const input = page.getByTestId("feature-name-input");
  await expect(input).toBeVisible();
  await input.fill(name);
  await input.press("Enter");
  await page.waitForTimeout(500);
  // Close the popover by pressing Escape
  await page.keyboard.press("Escape");
  await expect(page.getByTestId("feature-badges")).toContainText(name);
}
