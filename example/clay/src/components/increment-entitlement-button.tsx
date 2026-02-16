"use client";

import createMeterEvent from "@/app/actions/create-meter-event";
import { useTransition } from "react";
import { toast } from "sonner";
import { Button } from "./ui/button";

export function IncrementEntitlementButton({
  usage,
  featureId,
}: {
  usage: number;
  featureId: string;
}) {
  const [isPending, startTransition] = useTransition();
  const handleClick = () => {
    startTransition(async () => {
      const result = await createMeterEvent({ featureId, usage });

      if (!result.ok) {
        toast.error(result.message);
        return;
      }
    });
  };

  return (
    <Button
      size="sm"
      className="w-fit text-xs"
      onClick={handleClick}
      disabled={isPending}
    >
      Submit 🔼
    </Button>
  );
}
