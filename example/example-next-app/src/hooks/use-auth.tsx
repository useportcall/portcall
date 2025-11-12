"use client";

import { useRouter } from "next/navigation";
import React from "react";

type User = {
  id: string;
  email: string;
  name?: string;
};

const STORAGE_KEY = "demo_user";

function readUser(): User | null {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    return raw ? JSON.parse(raw) : null;
  } catch (e) {
    return null;
  }
}

export function useAuth() {
  const router = useRouter();
  const [user, setUser] = React.useState<User | null>(() => {
    if (typeof window === "undefined") return null;
    return readUser();
  });

  React.useEffect(() => {
    // sync on mount
    setUser(readUser());
  }, []);

  const login = async (email: string) => {
    const res = await fetch("/api/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email }),
    });

    const data = await res.json();
    const u: User = data?.user;
    if (u) {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(u));
      document.cookie = `user_id=${u.id}; path=/; secure; SameSite=Lax`;
      setUser(u);
      router.replace("/");
    }
    return u;
  };

  const logout = () => {
    localStorage.removeItem(STORAGE_KEY);
    document.cookie = `user_id=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT`;
    setUser(null);
    router.replace("/login");
  };

  console.log("useAuth user:", user);

  return { user, login, logout } as const;
}
