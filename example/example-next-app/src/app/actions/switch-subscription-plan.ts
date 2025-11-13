"use server";

import { revalidatePath } from "next/cache";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";

export default async function switchSubscriptionPlan(props: {
  subscriptionId: string;
  planId: string;
}) {
  const cookieStore = await cookies();

  const userId = cookieStore.get("user_id")?.value;

  if (!userId) {
    redirect("/login");
  }

  const url = new URL(
    `/v1/subscriptions/${props.subscriptionId}`,
    process.env.PC_API_BASE_URL
  );

  const headers = new Headers();
  headers.set("X-API-Key", `${process.env.PC_API_TOKEN}`);
  headers.set("Content-Type", "application/json");

  const res = await fetch(url.toString(), {
    method: "POST",
    headers: headers,
    body: JSON.stringify({ plan_id: props.planId }),
  });

  if (!res.ok) {
    return { ok: false, message: "Failed to cancel subscription" };
  }

  await new Promise((resolve) => setTimeout(resolve, 1000));

  revalidatePath("/");

  return { ok: true };
}
