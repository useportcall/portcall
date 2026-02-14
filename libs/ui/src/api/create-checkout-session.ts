"use server";

import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { getPortcallClient } from "../lib/portcall";

export async function createCheckoutSession(planId: string) {
  const cookieStore = await cookies();

  const userId = cookieStore.get("user_id")?.value;

  if (!userId) {
    redirect("/login");
  }

  try {
    const portcall = getPortcallClient();
    const session = await portcall.checkoutSessions.create({
      plan_id: planId,
      user_id: userId,
      cancel_url: "http://localhost:3000?status=canceled", // TODO: fix
      redirect_url: "http://localhost:3000?status=success", // TODO: fix
    });

    return { ok: true, url: session.url };
  } catch (error) {
    console.error("Failed to create checkout session:", error);
    return { ok: false, message: "Failed to create checkout session" };
  }
}
