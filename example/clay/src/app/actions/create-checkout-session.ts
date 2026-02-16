"use server";

import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { portcall } from "@/lib/portcall";
import { PLANS, type PlanId } from "@/portcall/client";

// Map plan IDs to camelCase plan keys used by the SDK
const PLAN_ID_TO_KEY: Record<string, keyof typeof portcall.plans> = {
  [PLANS.CLAY_FREE]: "clayFree",
  [PLANS.CLAY_STARTER]: "clayStarter",
  [PLANS.CLAY_EXPLORER]: "clayExplorer",
  [PLANS.CLAY_PRO]: "clayPro",
};

export async function createCheckoutSession(planId: string) {
  const cookieStore = await cookies();
  const userId = cookieStore.get("user_id")?.value;
  const userEmail = cookieStore.get("user_email")?.value;

  if (!userId || !userEmail) {
    redirect("/login");
  }

  try {
    const planKey = PLAN_ID_TO_KEY[planId];
    if (!planKey) {
      return { ok: false, message: "Invalid plan" };
    }

    const session = await portcall.plans[planKey].subscribe(userEmail, {
      cancelUrl: "http://localhost:3000/pricing?status=canceled",
      redirectUrl: "http://localhost:3000/pricing?status=success",
    });

    const checkoutBaseUrl =
      process.env.NEXT_PUBLIC_CHECKOUT_URL || "http://localhost:8700";
    const fullCheckoutUrl = session.url.startsWith("http")
      ? session.url
      : `${checkoutBaseUrl}${session.url}`;

    return { ok: true, url: fullCheckoutUrl };
  } catch (error) {
    console.error("Error creating checkout session:", error);
    return { ok: false, message: "Failed to create checkout session" };
  }
}
