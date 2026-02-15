// Dashboard actions: create a user via the browser UI.
import { Page, expect } from "@playwright/test";

/** Open the "Add user" dialog and create a user by email. */
export async function createUser(page: Page, email: string) {
  await page.getByTestId("add-user-button").click();
  await expect(page.getByText("Add new user")).toBeVisible();

  await page.getByTestId("create-user-email-input").fill(email);
  await page.getByTestId("create-user-submit").click();

  // Form navigates to /users/:id — wait for URL change
  await page.waitForURL(/\/users\/[a-z0-9_]+$/, { timeout: 10000 });
}
