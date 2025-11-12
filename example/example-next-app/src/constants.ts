export const PUBLIC_API = "http://localhost:8600";
export const API_KEY = process.env.NEXT_PUBLIC_API_KEY || "test";

export const HOSTED_DOMAIN =
  process.env.NEXT_PUBLIC_HOSTED_DOMAIN || "localhost:3000";

// pricing plan ids
export const CLAY_FREE_PLAN_ID = process.env.NEXT_PUBLIC_CLAY_FREE_PLAN_ID!;
export const CLAY_STARTER_PLAN_ID =
  process.env.NEXT_PUBLIC_CLAY_STARTER_PLAN_ID!;
export const CLAY_EXPLORER_PLAN_ID =
  process.env.NEXT_PUBLIC_CLAY_EXPLORER_PLAN_ID!;
export const CLAY_PRO_PLAN_ID = process.env.NEXT_PUBLIC_CLAY_PRO_PLAN_ID!;
