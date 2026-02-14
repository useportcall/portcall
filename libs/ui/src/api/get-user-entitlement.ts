import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { getPortcallClient } from "../lib/portcall";

export async function getUserEntitlement(featureId: string) {
  const cookieStore = await cookies();
  const userId = cookieStore.get("user_id")?.value;

  if (!userId) {
    redirect("/login");
  }

  try {
    const portcall = getPortcallClient();
    const entitlement = await portcall.entitlements.get(userId, featureId);
    return entitlement ?? null;
  } catch (error) {
    console.error(`Error fetching entitlement:`, error);
    throw new Error(`Error fetching entitlement`);
  }
}
