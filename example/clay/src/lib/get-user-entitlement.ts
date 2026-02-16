import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { portcall } from "./portcall";
import type { FeatureId } from "@/portcall/client";

export async function getUserEntitlement(featureId: string) {
  const cookieStore = await cookies();
  const userId = cookieStore.get("user_id")?.value;

  if (!userId) {
    redirect("/login");
  }

  try {
    // Use the generated SDK to get the entitlement
    const entitlement =
      await portcall.entitlements[
        featureId as keyof typeof portcall.entitlements
      ]?.get(userId);
    return entitlement ?? null;
  } catch (error) {
    console.error(`Error fetching entitlement for ${featureId}:`, error);
    return null;
  }
}
