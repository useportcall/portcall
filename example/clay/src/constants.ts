// Import from generated Portcall client
import { PLANS, FEATURES } from "@/portcall/client";

// Re-export from generated Portcall client for backward compatibility
export { PLANS, FEATURES, PLANS as PLAN_IDS } from "@/portcall/client";

// Keep legacy exports for backward compatibility with existing code
export const PUBLIC_API =
  process.env.PC_API_BASE_URL || "http://localhost:9080";
export const API_KEY = process.env.PC_API_SECRET || "test";

export const HOSTED_DOMAIN =
  process.env.NEXT_PUBLIC_HOSTED_DOMAIN || "localhost:3000";

// Clay pricing plan IDs from generated client
export const CLAY_FREE_PLAN_ID = PLANS.CLAY_FREE;
export const CLAY_STARTER_PLAN_ID = PLANS.CLAY_STARTER;
export const CLAY_EXPLORER_PLAN_ID = PLANS.CLAY_EXPLORER;
export const CLAY_PRO_PLAN_ID = PLANS.CLAY_PRO;
