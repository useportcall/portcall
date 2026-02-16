import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { portcall } from "./portcall";

export async function getUserSubscription() {
  const cookieStore = await cookies();
  const userId = cookieStore.get("user_id")?.value;

  console.log("Fetching subscription for user_id:", userId);

  if (!userId) {
    redirect("/login");
  }

  try {
    // Use the generated SDK to get the active subscription
    const subscription = await portcall.subscriptions.getActive(userId);
    return subscription;
  } catch (error) {
    console.error("Error fetching subscription:", error);
    return null;
  }
}
