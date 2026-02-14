"use server";

import { revalidatePath } from "next/cache";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { getPortcallClient } from "../lib/portcall";

export default async function cancelSubscription(subscriptionId: string) {
  const cookieStore = await cookies();

  const userId = cookieStore.get("user_id")?.value;

  if (!userId) {
    redirect("/login");
  }

  try {
    const portcall = getPortcallClient();
    const subscription = await portcall.subscriptions.cancel(subscriptionId);

    await new Promise((resolve) => setTimeout(resolve, 1000));

    revalidatePath("/");

    return { ok: true, subscription };
  } catch (error) {
    console.error("Failed to cancel subscription:", error);
    return { ok: false, message: "Failed to cancel subscription" };
  }
}
