import { cookies } from "next/headers";
import Link from "next/link";
import { redirect } from "next/navigation";
import React from "react";
import { LogoutButton } from "./logout-button";

async function checkAuth() {
  const cookieStore = await cookies();

  const userId = cookieStore.get("user_id")?.value;

  if (!userId) {
    redirect("/login");
  }

  return true;
}

export async function AuthenticatedLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  await checkAuth();

  return (
    <div className="min-h-screen flex flex-col items-center w-screen">
      <nav className="py-4 px-10 flex justify-between items-center w-full fixed">
        <div className="flex gap-10 items-center">
          <Link href="/">Home</Link>
          <Link href="/pricing" className="ml-4">
            Pricing
          </Link>
        </div>
        <LogoutButton />
      </nav>

      <main className="max-w-5xl w-full">{children}</main>
    </div>
  );
}
