"use client";

import React, { useTransition } from "react";
import { toast } from "sonner";
import { useRouter } from "next/navigation";
import login from "./api/login";

export default function LoginView() {
  const [email, setEmail] = React.useState("");
  const [isPending, startTransition] = useTransition();
  const router = useRouter();

  const handleClick = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    startTransition(async () => {
      const result = await login({ email });

      if (!result.ok) {
        toast.error(result.message);
        return;
      }

      router.replace("/");
    });
  };

  return (
    <div className="h-screen w-screen flex items-center justify-center p-4">
      <form onSubmit={handleClick} className="w-full max-w-sm">
        <label className="block mb-2">Email</label>
        <input
          className="w-full p-2 border rounded mb-4"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <button
          className="w-full p-2 bg-blue-600 text-white rounded"
          disabled={isPending}
        >
          {isPending ? "Logging in..." : "Login"}
        </button>
      </form>
    </div>
  );
}
