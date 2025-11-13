"use client";

import switchSubscriptionPlan from "@/app/actions/switch-subscription-plan";
import { Button } from "@/components/ui/button";
import { Loader2 } from "lucide-react";
import { useRouter } from "next/navigation";
import { useTransition } from "react";
import { toast } from "sonner";

export function SwitchSubscriptionPlanButton({
  subscriptionId,
  planId,
}: {
  subscriptionId: string;
  planId: string;
}) {
  const router = useRouter();

  const [isPending, startTransition] = useTransition();

  const handleClick = () => {
    startTransition(async () => {
      const result = await switchSubscriptionPlan({ subscriptionId, planId });

      if (!result.ok) {
        toast.error(result.message);
        return;
      }

      // reload page to reflect new plan
      router.refresh();
    });
  };

  return (
    <Button onClick={handleClick} disabled={isPending}>
      {isPending && <Loader2 className="animate-spin h-4 w-4" />}
      Yes
    </Button>
  );
}
