"use server";

import { cookies } from "next/headers";

export default async function login(props: { email: string }) {
  console.log("Logging in with email:", props.email);

  const url = new URL(
    `/v1/users?email=${encodeURIComponent(props.email)}`,
    process.env.PC_API_BASE_URL
  );

  const headers = new Headers();
  headers.set("X-API-Key", `${process.env.PC_API_TOKEN}`);
  headers.set("Content-Type", "application/json");

  const res = await fetch(url.toString(), {
    method: "GET",
    headers: headers,
  });

  if (!res.ok) {
    return { ok: false, message: "Failed to login" };
  }

  await new Promise((resolve) => setTimeout(resolve, 1000));

  const result = await res.json();

  const user = result?.data?.[0] ?? null;

  if (!user) {
    return { ok: false, message: "Failed to login" };
  }

  const cookieStore = await cookies();

  cookieStore.set("user_id", user.id, {
    httpOnly: true,
    expires: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000), // 30 days
    sameSite: "lax",
    path: "/",
    secure: process.env.NODE_ENV === "production",
  });

  console.log(
    "cookie set successfully, redirecting to home page.",
    cookieStore.get("user_id")
  );

  return { ok: true };
}
