// Full-page screenshot helper. Saves locally as a test artifact and
// optionally uploads to a Discord webhook when running in live mode.
import { Page, test } from "@playwright/test";
import { HarnessConfig } from "./config";
import * as fs from "fs";

interface SnapshotConfig {
  mode: string;
  webhook_url: string;
}

let cachedConfig: SnapshotConfig | null = null;

async function getSnapshotConfig(cfg: HarnessConfig): Promise<SnapshotConfig> {
  if (cachedConfig) return cachedConfig;
  const res = await fetch(`${cfg.control_url}/e2e/snapshot/config`);
  cachedConfig = (await res.json()) as SnapshotConfig;
  return cachedConfig;
}

/** Take a full-page screenshot, save it locally, and send to Discord if configured. */
export async function takeSnapshot(
  page: Page,
  cfg: HarnessConfig,
  name: string,
) {
  const filePath = test.info().outputPath(`${name}.png`);
  await page.screenshot({ path: filePath, fullPage: true });

  const snapCfg = await getSnapshotConfig(cfg);
  if (snapCfg.mode !== "live" || !snapCfg.webhook_url) return;

  const bytes = fs.readFileSync(filePath);
  const blob = new Blob([bytes], { type: "image/png" });
  const form = new FormData();
  form.append("file", blob, `${name}.png`);
  form.append("payload_json", JSON.stringify({ content: name }));
  await fetch(snapCfg.webhook_url, { method: "POST", body: form });
}
