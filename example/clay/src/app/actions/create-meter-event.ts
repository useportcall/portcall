"use server";

import { revalidatePath } from "next/cache";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { portcall } from "@/lib/portcall";
import type { FeatureId } from "@/portcall/client";

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
    // Use the SDK to record the meter event
    await portcall.meterEvents.record(
      userId,
      props.featureId as FeatureId,
      props.usage,
    );

    await new Promise((resolve) => setTimeout(resolve, 1000));

    revalidatePath(`/`);

    return { ok: true };
  } catch (error) {
    console.error("Error creating meter event:", error);
    return { ok: false, message: "Failed to increment entitlement" };
  }
}
