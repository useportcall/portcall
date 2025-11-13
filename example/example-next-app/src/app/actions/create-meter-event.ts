"use server";

import { revalidatePath } from "next/cache";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";

export default async function createMeterEvent(props: { featureId: string }) {
  const cookieStore = await cookies();

  const userId = cookieStore.get("user_id")?.value;

  if (!userId) {
    redirect("/login");
  }

  const url = new URL(`/v1/meter-events`, process.env.PC_API_BASE_URL);

  const headers = new Headers();
  headers.set("X-API-Key", `${process.env.PC_API_TOKEN}`);
  headers.set("Content-Type", "application/json");

  const res = await fetch(url.toString(), {
    method: "POST",
    headers: headers,
    body: JSON.stringify({
      feature_id: props.featureId,
      user_id: userId,
      usage: 10,
    }),
  });

  if (!res.ok) {
    return { ok: false, message: "Failed to cancel subscription" };
  }

  await new Promise((resolve) => setTimeout(resolve, 1000));

  revalidatePath(`/v1/users/${userId}/entitlements/${props.featureId}`);

  return { ok: true };
}
