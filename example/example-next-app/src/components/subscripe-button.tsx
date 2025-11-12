"use client";

import { Button } from "@/components/ui/button";
import { Loader2 } from "lucide-react";
import { useRouter } from "next/navigation";
import { useTransition } from "react";
import { toast } from "sonner";
import { createCheckoutSession } from "../app/actions/create-checkout-session";

export function SubscribeButton({ planId }: { planId: string }) {
  const router = useRouter();
  const [isPending, startTransition] = useTransition();

  const handleClick = () => {
    startTransition(async () => {
      const result = await createCheckoutSession(planId);

      if (!result.ok) {
        toast.error(result.message);
        return;
      }

      // ✅ redirect based on returned URL
      router.push(result.url!);
    });
  };

  return (
    <Button onClick={handleClick} disabled={isPending}>
      {isPending && <Loader2 className="animate-spin h-4 w-4" />}
      Yes
    </Button>
  );
}
