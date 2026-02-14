"use server";

import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { getPortcallClient } from "../lib/portcall";

export async function getUserSubscription() {
  const cookieStore = await cookies();
  const userId = cookieStore.get("user_id")?.value;

  if (!userId) {
    redirect("/login");
  }

  try {
    const portcall = getPortcallClient();
    const subscriptions = await portcall.subscriptions.list({
      user_id: userId,
      status: "active",
    });

    const subscription = subscriptions?.[0] ?? null;

    return { data: subscription, authenticated: true };
  } catch (error) {
    console.error(`Error fetching subscriptions:`, error);
    throw new Error(`Error fetching subscriptions`);
  }
}
