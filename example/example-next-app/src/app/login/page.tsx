"use client";

import React from "react";
import { useAuth } from "@/hooks/use-auth";

export default function LoginPage() {
  const { login } = useAuth();
  const [email, setEmail] = React.useState("");
  const [loading, setLoading] = React.useState(false);

  const handle = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    await login(email || "demo@local");
    setLoading(false);
  };

  return (
    <div className="h-screen w-screen flex items-center justify-center p-4">
      <form onSubmit={handle} className="w-full max-w-sm">
        <label className="block mb-2">Email</label>
        <input
          className="w-full p-2 border rounded mb-4"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <button
          className="w-full p-2 bg-blue-600 text-white rounded"
          onClick={(e) => handle(e)}
          disabled={loading}
        >
          {loading ? "Logging in..." : "Login"}
        </button>
      </form>
    </div>
  );
}
