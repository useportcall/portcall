"use server";

import { cookies } from "next/headers";
import { redirect } from "next/navigation";

export async function logout() {
  const { delete: deleteCookie } = await cookies();

  // Delete the cookie
  deleteCookie("user_id");

  // Optionally redirect
  redirect("/login");
}
