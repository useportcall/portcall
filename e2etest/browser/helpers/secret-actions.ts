// Dashboard actions: create API secret via the browser UI.
import { Page, expect } from "@playwright/test";

/** Create a new API secret and return the key value. */
export async function createSecret(page: Page): Promise<string> {
  await page.getByTestId("create-secret-button").click();
  await expect(page.getByText("API secret created")).toBeVisible();

  const textarea = page.getByTestId("secret-key-value");
  await expect(textarea).toBeVisible();
  const key = await textarea.inputValue();

  await page.getByTestId("close-secret-dialog").click();
  return key;
}
