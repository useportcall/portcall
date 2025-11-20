import { cookies } from "next/headers";
import { redirect } from "next/navigation";

export async function getUserSubscription() {
  const cookieStore = await cookies();
  const userId = cookieStore.get("user_id")?.value;

  console.log("Fetching subscription for user_id:", userId);

  if (!userId) {
    redirect("/login");
  }

  const url = new URL("/v1/subscriptions", process.env.PC_API_BASE_URL);
  url.searchParams.set("user_id", userId);
  url.searchParams.set("status", "active");
  url.searchParams.set("limit", "1");

  const headers = new Headers();
  headers.set("X-API-Key", `${process.env.PC_API_TOKEN}`);
  headers.set("Content-Type", "application/json");

  const res = await fetch(url.toString(), {
    method: "GET",
    headers: headers,
    cache: "no-store",
  });

  if (!res.ok) {
    throw new Error(`Error fetching subscriptions: ${res.statusText}`);
  }

  const { data }: { data: any[] } = await res.json();

  const subscription = data?.[0] ?? null;

  return subscription;
}
