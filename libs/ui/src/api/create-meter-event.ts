"use server";

import { revalidatePath } from "next/cache";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { getPortcallClient } from "../lib/portcall";

export default async function createMeterEvent(props: {
  featureId: string;
  usage: number;
}) {
  const cookieStore = await cookies();

  const userId = cookieStore.get("user_id")?.value;

  if (!userId) {
    redirect("/login");
  }

  try {
    const portcall = getPortcallClient();
    await portcall.meterEvents.create({
      feature_id: props.featureId,
      user_id: userId,
      usage: props.usage,
    });

    await new Promise((resolve) => setTimeout(resolve, 1000));

    revalidatePath(`/v1/users/${userId}/entitlements/${props.featureId}`);

    return { ok: true };
  } catch (error) {
    console.error("Failed to create meter event:", error);
    return { ok: false, message: "Failed to increment entitlement" };
  }
}
