// Reads the harness config written by global-setup.
import * as fs from "fs";
import * as path from "path";

export interface HarnessConfig {
  dashboard_url: string;
  api_url: string;
  checkout_url: string;
  quote_url: string;
  file_url: string;
  control_url: string;
  admin_api_key: string;
  app_id: number;
  app_public_id: string;
}

const CONFIG_PATH = path.join(__dirname, "../setup/.harness-config.json");

export function loadConfig(): HarnessConfig {
  const raw = readConfigWithRetry();
  return JSON.parse(raw) as HarnessConfig;
}

function readConfigWithRetry(): string {
  const deadline = Date.now() + 5_000;
  while (Date.now() < deadline) {
    try {
      return fs.readFileSync(CONFIG_PATH, "utf8");
    } catch {
      Atomics.wait(new Int32Array(new SharedArrayBuffer(4)), 0, 0, 50);
    }
  }
  throw new Error(`Harness config not found at ${CONFIG_PATH}`);
}
