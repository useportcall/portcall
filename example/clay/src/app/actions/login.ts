"use server";

import { cookies } from "next/headers";
import { portcall } from "@/lib/portcall";

export default async function login(props: { email: string }) {
  try {
    const user = await portcall.users.findOrCreate(props.email);
    if (!user) {
      return { ok: false, message: "Failed to login" };
    }

    const cookieStore = await cookies();
    const cookieOptions = {
      httpOnly: true,
      expires: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000),
      sameSite: "lax" as const,
      path: "/",
      secure: process.env.NODE_ENV === "production",
    };

    cookieStore.set("user_id", user.id, cookieOptions);
    cookieStore.set("user_email", user.email, cookieOptions);

    return { ok: true };
  } catch (error) {
    console.error("Login error:", error);
    return { ok: false, message: "Failed to login" };
  }
}
