"use server";

import { revalidatePath } from "next/cache";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { portcall } from "@/lib/portcall";
import type { PlanId } from "@/portcall/client";

export default async function switchSubscriptionPlan(props: {
  subscriptionId: string;
  planId: string;
}) {
  const cookieStore = await cookies();

  const userId = cookieStore.get("user_id")?.value;

  if (!userId) {
    redirect("/login");
  }

  try {
    // Use the SDK to update the subscription
    await portcall.subscriptions.update(props.subscriptionId, {
      plan_id: props.planId as PlanId,
    });

    await new Promise((resolve) => setTimeout(resolve, 1000));

    revalidatePath("/");

    return { ok: true };
  } catch (error) {
    console.error("Error switching subscription plan:", error);
    return { ok: false, message: "Failed to switch subscription plan" };
  }
}
