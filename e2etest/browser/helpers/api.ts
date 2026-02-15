import { HarnessConfig } from "./config";

type JSONValue = Record<string, unknown>;

function dashHeaders(cfg: HarnessConfig) {
  return {
    "Content-Type": "application/json",
    "X-Admin-API-Key": cfg.admin_api_key,
    "X-Target-App-ID": String(cfg.app_id),
  };
}

async function parseJSON(res: Response): Promise<JSONValue> {
  return (await res.json()) as JSONValue;
}

export async function dashGet(cfg: HarnessConfig, path: string): Promise<JSONValue> {
  const res = await fetch(`${cfg.dashboard_url}/api/apps/${cfg.app_public_id}/${path}`, { headers: dashHeaders(cfg) });
  const json = await parseJSON(res);
  if (!res.ok) throw new Error(`GET ${path}: ${res.status}`);
  return (json as any).data ?? json;
}

export async function dashPost(cfg: HarnessConfig, path: string, body?: JSONValue): Promise<JSONValue> {
  const res = await fetch(`${cfg.dashboard_url}/api/apps/${cfg.app_public_id}/${path}`, {
    method: "POST", headers: dashHeaders(cfg), body: body ? JSON.stringify(body) : undefined,
  });
  const json = await parseJSON(res);
  if (!res.ok) throw new Error(`POST ${path}: ${res.status}`);
  return (json as any).data ?? json;
}

export async function dashPostRoot(cfg: HarnessConfig, path: string, body?: JSONValue): Promise<JSONValue> {
  const res = await fetch(`${cfg.dashboard_url}/api/${path}`, {
    method: "POST", headers: dashHeaders(cfg), body: body ? JSON.stringify(body) : undefined,
  });
  const json = await parseJSON(res);
  if (!res.ok) throw new Error(`POST /api/${path}: ${res.status}`);
  return (json as any).data ?? json;
}

export async function apiPost(cfg: HarnessConfig, apiKey: string, path: string, body?: JSONValue): Promise<JSONValue> {
  const res = await fetch(`${cfg.api_url}${path}`, {
    method: "POST",
    headers: { "Content-Type": "application/json", "x-api-key": apiKey },
    body: body ? JSON.stringify(body) : undefined,
  });
  const json = await parseJSON(res);
  if (!res.ok) throw new Error(`POST ${path}: ${res.status}`);
  return (json as any).data ?? json;
}

export async function controlPost(cfg: HarnessConfig, path: string, body: JSONValue): Promise<JSONValue> {
  const res = await fetch(`${cfg.control_url}${path}`, {
    method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(body),
  });
  const json = await parseJSON(res);
  if (!res.ok) throw new Error(`CTRL ${path}: ${res.status}`);
  return json;
}

export async function controlGet(cfg: HarnessConfig, path: string): Promise<JSONValue> {
  const res = await fetch(`${cfg.control_url}${path}`);
  const json = await parseJSON(res);
  if (!res.ok) throw new Error(`CTRL ${path}: ${res.status}`);
  return json;
}

export async function dashGetBlob(cfg: HarnessConfig, path: string): Promise<Blob> {
  const res = await fetch(`${cfg.dashboard_url}/api/apps/${cfg.app_public_id}/${path}`, { headers: dashHeaders(cfg) });
  if (!res.ok) throw new Error(`GET ${path}: ${res.status}`);
  return res.blob();
}
