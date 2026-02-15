// API-based seed helpers for creating data without browser interaction.
// Use these for fast setup in snapshot and smoke tests.
import { HarnessConfig } from "./config";
import { dashPost, apiPost, controlPost } from "./api";

/** Create and publish a plan via the dashboard API. Returns the plan public_id. */
export async function seedPlan(
  cfg: HarnessConfig,
  name: string,
  cents: number,
): Promise<string> {
  const plan = (await dashPost(cfg, "plans")) as any;
  await dashPost(cfg, `plans/${plan.id}`, { name, currency: "usd" });
  await dashPost(cfg, `plan-items/${plan.items[0].id}`, {
    unit_amount: cents,
  });
  await dashPost(cfg, `plans/${plan.id}/publish`);
  return plan.id as string;
}

/** Create an API secret via the dashboard API. Returns the secret key. */
export async function seedSecret(cfg: HarnessConfig): Promise<string> {
  const secret = (await dashPost(cfg, "secrets")) as any;
  return (secret.key ?? secret.secret_key) as string;
}

/** Create a user with billing address via the public API. Returns the user public_id. */
export async function seedUser(
  cfg: HarnessConfig,
  key: string,
  email: string,
): Promise<string> {
  const user = (await apiPost(cfg, key, "/v1/users", {
    name: "Snapshot User",
    email,
  })) as any;
  await apiPost(cfg, key, `/v1/users/${user.id}/billing-address`, {
    line1: "1 Test Blvd",
    city: "San Francisco",
    postal_code: "94105",
    country: "US",
  });
  return user.id as string;
}

/** Create a checkout session. Returns { url, external_session_id }. */
export async function seedCheckoutSession(
  cfg: HarnessConfig,
  key: string,
  planId: string,
  userId: string,
): Promise<{ url: string; externalSessionId: string }> {
  const session = (await apiPost(cfg, key, "/v1/checkout-sessions", {
    plan_id: planId,
    user_id: userId,
    redirect_url: cfg.dashboard_url + "/subscriptions",
    cancel_url: cfg.dashboard_url + "/plans",
  })) as any;
  return {
    url: normalizeCheckoutURL(session.url, cfg.checkout_url),
    externalSessionId: session.external_session_id,
  };
}

function normalizeCheckoutURL(raw: string, base: string): string {
  try {
    const parsed = new URL(raw);
    const id = parsed.searchParams.get("id");
    const token = parsed.searchParams.get("st");
    if (!id || !token) return raw;
    const out = new URL(base);
    out.searchParams.set("id", id);
    out.searchParams.set("st", token);
    return out.toString();
  } catch {
    return raw;
  }
}

/** Seed a test invoice via the control server. Returns the invoice public_id. */
export async function seedInvoice(cfg: HarnessConfig): Promise<string> {
  const res = (await controlPost(cfg, "/e2e/seed-invoice", {})) as any;
  return res.invoice_id as string;
}
