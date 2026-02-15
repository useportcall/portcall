// Dashboard page-object helpers for Playwright tests.
import { Page, expect } from "@playwright/test";

/** Navigate to plans page via sidebar. */
export async function goToPlans(page: Page) {
  const plansHeading = page.getByRole("heading", { name: "Plans" });
  const visible = await plansHeading.isVisible({ timeout: 2_000 }).catch(() => false);
  if (!visible) {
    const origin = new URL(page.url()).origin;
    await page.goto(origin + "/plans");
  }
  await expect(plansHeading).toBeVisible();
}

/** Navigate to users page via sidebar. */
export async function goToUsers(page: Page) {
  await page.getByRole("link", { name: "Users" }).click();
  await expect(page.getByTestId("add-user-button")).toBeVisible();
}

/** Navigate to developer page via sidebar. */
export async function goToDeveloper(page: Page) {
  await page.getByRole("link", { name: "Developer" }).click();
  await expect(page.getByRole("heading", { name: "Developer" })).toBeVisible();
}

/** Navigate to integrations page via sidebar. */
export async function goToIntegrations(page: Page) {
  await page.getByRole("link", { name: "Payment integrations" }).click();
  await expect(
    page.getByText("Payment Integrations", { exact: true }),
  ).toBeVisible();
}

/** Navigate to subscriptions page via sidebar. */
export async function goToSubscriptions(page: Page) {
  await page.getByRole("link", { name: "Subscriptions" }).click();
  await expect(page.getByText("Subscriptions")).toBeVisible();
}

/** Navigate to invoices page via sidebar. */
export async function goToInvoices(page: Page) {
  await page.getByRole("link", { name: "Invoices" }).click();
  await expect(page.getByText("Invoices")).toBeVisible();
}
