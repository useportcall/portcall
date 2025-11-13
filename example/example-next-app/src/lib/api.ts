import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { NextResponse } from "next/server";

export async function callPortcallApi(
  method: "GET" | "POST",
  path: string,
  body?: Object
) {
  const url = new URL(path, process.env.PC_API_BASE_URL);

  const headers = new Headers();
  headers.set("Authorization", `Bearer ${process.env.PC_API_TOKEN}`);
  headers.set("Content-Type", "application/json");

  try {
    const res = await fetch(url.toString(), {
      method: method,
      headers: headers,
      body: body ? JSON.stringify(body) : undefined,
    });

    const data: { data: unknown } = await res.json();

    return NextResponse.json(data.data, {
      status: res.status,
      headers: {
        "Cache-Control": res.headers.get("Cache-Control") ?? "no-cache",
      },
    });
  } catch (err) {
    console.error("portcall fetch failed:", url.toString(), err);
    return NextResponse.json(
      { error: "Internal server error" },
      { status: 502 }
    );
  }
}

export async function getUserSubscription() {
  const cookieStore = await cookies();
  const userId = cookieStore.get("user_id")?.value;

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

  console.log("User subscriptions data:", data);

  const subscription = data?.[0] ?? null;

  return subscription;
}

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

  const entitlement = data ?? null;

  return entitlement;
}
