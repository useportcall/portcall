import { CheckoutSession } from "@/types/api";

export type CheckoutConsentMode = "save" | "charge";

const recurringIntervals = new Set(["day", "week", "month", "year"]);

function stripeIntentID(session: CheckoutSession): string {
  if (session.external_session_id) return session.external_session_id;
  if (!session.external_client_secret) return "";
  return session.external_client_secret.split("_secret_")[0] || "";
}

export function resolveCheckoutConsentMode(
  session: CheckoutSession,
): CheckoutConsentMode {
  if (session.external_provider === "stripe") {
    const id = stripeIntentID(session);
    if (id.startsWith("pi_")) return "charge";
    if (id.startsWith("seti_")) return "save";
  }

  if (session.plan && recurringIntervals.has(session.plan.interval)) {
    return "save";
  }

  return "save";
}
