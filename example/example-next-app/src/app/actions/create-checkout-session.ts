"use server";

import { cookies } from "next/headers";
import { redirect } from "next/navigation";

export async function createCheckoutSession(planId: string) {
  const cookieStore = await cookies();

  const userId = cookieStore.get("user_id")?.value;

  if (!userId) {
    redirect("/login");
  }

  const url = new URL("/v1/checkout-sessions", process.env.PC_API_BASE_URL);

  const headers = new Headers();
  headers.set("Authorization", `Bearer ${process.env.PC_API_TOKEN}`);
  headers.set("Content-Type", "application/json");

  const res = await fetch(url.toString(), {
    method: "POST",
    headers: headers,
    body: JSON.stringify({
      plan_id: planId,
      user_id: userId,
      cancel_url: "http://localhost:3000?status=canceled", // TODO: fix
      redirect_url: "http://localhost:3000?status=success", // TODO: fix
    }),
  });

  if (!res.ok) {
    return { ok: false, message: "Failed to create checkout session" };
  }

  const data: { data: any } = await res.json();

  return { ok: true, url: data.data.url };
}
