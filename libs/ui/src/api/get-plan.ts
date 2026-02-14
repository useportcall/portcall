"use server";

import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { getPortcallClient } from "../lib/portcall";

export async function getPlan(planId: string) {
  const cookieStore = await cookies();
  const userId = cookieStore.get("user_id")?.value;

  if (!userId) {
    redirect("/login");
  }

  try {
    const portcall = getPortcallClient();
    const plan = await portcall.plans.get(planId);
    return { data: plan };
  } catch (error) {
    console.error(`Error fetching plan:`, error);
    throw new Error(`Error fetching plan`);
  }
}
