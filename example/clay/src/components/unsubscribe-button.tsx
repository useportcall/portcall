"use client";

import { Button } from "@/components/ui/button";
import { Loader2 } from "lucide-react";
import { useTransition } from "react";
import { toast } from "sonner";
import cancelSubscription from "../app/actions/cancel-subscription";

export function UnsubscribeButton({
  subscriptionId,
}: {
  subscriptionId: string;
}) {
  const [isPending, startTransition] = useTransition();

  const handleClick = () => {
    startTransition(async () => {
      const result = await cancelSubscription(subscriptionId);

      if (!result.ok) {
        toast.error(result.message);
        return;
      }
    });
  };

  return (
    <Button onClick={handleClick} disabled={isPending}>
      {isPending && <Loader2 className="animate-spin h-4 w-4" />}
      Yes
    </Button>
  );
}
