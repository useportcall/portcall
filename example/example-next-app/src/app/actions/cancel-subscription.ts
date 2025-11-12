"use server";

import { cookies } from "next/headers";
import { redirect } from "next/navigation";

export default async function cancelSubscription(subscriptionId: string) {
  const cookieStore = await cookies();

  const userId = cookieStore.get("user_id")?.value;

  if (!userId) {
    redirect("/login");
  }

  const url = new URL(
    `/v1/subscriptions/${subscriptionId}/cancel`,
    process.env.PC_API_BASE_URL
  );

  const headers = new Headers();
  headers.set("Authorization", `Bearer ${process.env.PC_API_TOKEN}`);
  headers.set("Content-Type", "application/json");

  const res = await fetch(url.toString(), {
    method: "POST",
    headers: headers,
    body: JSON.stringify({}),
  });

  if (!res.ok) {
    return { ok: false, message: "Failed to cancel subscription" };
  }

  const data: { data: any } = await res.json();

  return { ok: true, url: data.data.url };
}
