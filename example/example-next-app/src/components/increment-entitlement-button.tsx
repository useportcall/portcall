"use client";

import createMeterEvent from "@/app/actions/create-meter-event";
import { useTransition } from "react";
import { toast } from "sonner";
import { Button } from "./ui/button";

export function IncrementEntitlementButton({
  featureId,
}: {
  featureId: string;
}) {
  const [isPending, startTransition] = useTransition();
  const handleClick = () => {
    startTransition(async () => {
      const result = await createMeterEvent({ featureId });

      if (!result.ok) {
        toast.error(result.message);
        return;
      }
    });
  };

  return (
    <Button
      size="sm"
      className="w-fit mt-1 text-xs"
      onClick={handleClick}
      disabled={isPending}
    >
      Increment usage 🔼
    </Button>
  );
}
