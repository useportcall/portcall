"use client";

import switchSubscriptionPlan from "@/app/actions/switch-subscription-plan";
import { Button } from "@/components/ui/button";
import { Loader2, Repeat } from "lucide-react";
import { useState, useTransition } from "react";
import { toast } from "sonner";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogTitle,
  DialogTrigger,
} from "./ui/dialog";

export function SwitchSubscriptionPlanButton({
  subscriptionId,
  planId,
}: {
  subscriptionId: string;
  planId: string;
}) {
  const [open, setOpen] = useState(false);
  const [isPending, startTransition] = useTransition();

  const handleClick = () => {
    startTransition(async () => {
      const result = await switchSubscriptionPlan({ subscriptionId, planId });

      if (!result.ok) {
        toast.error(result.message);
        return;
      }

      toast.success("Plan switched successfully");

      setOpen(false);
    });
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button type="submit">
          Switch <Repeat />
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogTitle>Switch Plan</DialogTitle>
        <DialogDescription>
          Are you sure you want to switch plan?
        </DialogDescription>
        <DialogFooter>
          <Button onClick={handleClick} disabled={isPending}>
            {isPending && <Loader2 className="animate-spin h-4 w-4" />}
            Yes
          </Button>
          <DialogClose asChild>
            <Button variant="outline">No</Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
