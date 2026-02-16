import { createPortcallClient } from "@/portcall/client";

/**
 * Server-side Portcall client
 *
 * This client is configured with your API key and is ready to use.
 * Import this in your server components and server actions.
 */
export const portcall = createPortcallClient({
  apiKey: process.env.PC_API_SECRET!,
  baseURL: process.env.PC_API_BASE_URL || "http://localhost:9080",
});

// Re-export types and constants for convenience
export { PLANS, FEATURES } from "@/portcall/client";
export type {
  PlanName,
  PlanId,
  FeatureId,
  User,
  Subscription,
  Entitlement,
  CheckoutSession,
  Plan,
} from "@/portcall/client";
