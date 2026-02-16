"use server";

import { revalidatePath } from "next/cache";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { portcall } from "@/lib/portcall";

export default async function cancelSubscription(subscriptionId: string) {
  const cookieStore = await cookies();

  const userId = cookieStore.get("user_id")?.value;

  if (!userId) {
    redirect("/login");
  }

  try {
    // Use the SDK to cancel the subscription
    await portcall.subscriptions.cancel(subscriptionId);

    await new Promise((resolve) => setTimeout(resolve, 1000));

    revalidatePath("/");

    return { ok: true };
  } catch (error) {
    console.error("Error canceling subscription:", error);
    return { ok: false, message: "Failed to cancel subscription" };
  }
}
