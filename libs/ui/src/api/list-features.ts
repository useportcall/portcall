"use server";

import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { getPortcallClient } from "../lib/portcall";

export async function listFeatures(isMetered?: boolean) {
  const cookieStore = await cookies();
  const userId = cookieStore.get("user_id")?.value;

  if (!userId) {
    redirect("/login");
  }

  try {
    const portcall = getPortcallClient();
    const features = await portcall.features.list({
      is_metered: isMetered,
    });

    return { data: features };
  } catch (error) {
    console.error(`Error fetching features:`, error);
    throw new Error(`Error fetching features`);
  }
}
