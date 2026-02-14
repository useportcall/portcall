"use client";

import { Button } from "../components/ui/button";
import { Loader2 } from "lucide-react";
import { useRouter } from "next/navigation";
import React, { useTransition } from "react";
import { toast } from "sonner";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogTitle,
  DialogTrigger,
} from "../components/ui/dialog";
import { ArrowRight } from "lucide-react";
import { createCheckoutSession } from "../api/create-checkout-session";

export function SubscribeButton({
  planId,
  className,
}: {
  planId: string;
  className?: string;
}) {
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
    <Dialog>
      <DialogTrigger asChild>
        <Button type="submit" className={className}>
          Subscribe <ArrowRight />
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogTitle>Subscribe</DialogTitle>
        <DialogDescription>
          You will be redirected to the checkout page.
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
