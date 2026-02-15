// Dashboard actions: create and publish a plan via the browser UI.
import { Page, expect } from "@playwright/test";

/** Click "Add plan" and wait for the plan editor to load. */
export async function createPlan(page: Page) {
  await page.getByTestId("create-plan-button").click();
  await expect(page.getByTestId("plan-name-input")).toBeVisible();
}

/** Set plan name (clear + type + blur to trigger save). */
export async function setPlanName(page: Page, name: string) {
  const input = page.getByTestId("plan-name-input");
  await input.click();
  await input.fill(name);
  await input.press("Tab"); // blur triggers save
  await page.waitForTimeout(300);
}

/** Set fixed unit price (clear + type + blur to trigger save). */
export async function setPlanPrice(page: Page, price: string) {
  const input = page.getByTestId("plan-price-input");
  await input.click();
  await input.fill(price);
  await input.press("Tab");
  await page.waitForTimeout(300);
}

/** Publish the plan: click "Save and publish" then "Confirm". */
export async function publishPlan(page: Page) {
  await page.getByTestId("publish-plan-button").click();
  await page.getByTestId("confirm-publish-button").click();
  // After publish, "Published" button and success dialog may appear
  await expect(page.getByText("Plan published!")).toBeVisible({
    timeout: 10000,
  });
  // Close the dialog if it auto-shows
  const closeBtn = page.getByRole("button", { name: "Close" }).first();
  if (await closeBtn.isVisible()) await closeBtn.click();
}
