import * as fs from "fs";
import { HarnessConfig } from "./config";

interface SnapshotConfig {
  mode: string;
  webhook_url: string;
}

async function getSnapshotConfig(cfg: HarnessConfig): Promise<SnapshotConfig> {
  const res = await fetch(`${cfg.control_url}/e2e/snapshot/config`);
  return (await res.json()) as SnapshotConfig;
}

export async function uploadArtifactIfLive(
  cfg: HarnessConfig,
  filePath: string,
  fileName: string,
  contentType: string,
  message: string,
): Promise<boolean> {
  const snapCfg = await getSnapshotConfig(cfg);
  if (snapCfg.mode !== "live" || !snapCfg.webhook_url) return false;

  try {
    const bytes = fs.readFileSync(filePath);
    const blob = new Blob([bytes], { type: contentType });
    const form = new FormData();
    form.append("file", blob, fileName);
    form.append("payload_json", JSON.stringify({ content: message }));
    await fetch(snapCfg.webhook_url, { method: "POST", body: form });
    return true;
  } catch {
    return false;
  }
}
