import { cookies } from "next/headers";
import { redirect } from "next/navigation";

export async function getUserEntitlement(featureId: string) {
  const cookieStore = await cookies();
  const userId = cookieStore.get("user_id")?.value;

  if (!userId) {
    redirect("/login");
  }

  const url = new URL(
    `/v1/users/${userId}/entitlements/${featureId}`,
    process.env.PC_API_BASE_URL
  );

  const headers = new Headers();
  headers.set("X-API-Key", `${process.env.PC_API_TOKEN}`);
  headers.set("Content-Type", "application/json");

  const res = await fetch(url.toString(), {
    method: "GET",
    headers: headers,
    cache: "no-store",
  });

  if (!res.ok) {
    throw new Error(`Error fetching entitlement: ${res.statusText}`);
  }

  const { data }: { data: any } = await res.json();

  return data ?? null;
}
