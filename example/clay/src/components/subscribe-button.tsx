"use client";

import { Button } from "@/components/ui/button";
import { Loader2 } from "lucide-react";
import { useRouter } from "next/navigation";
import React, { useTransition } from "react";
import { toast } from "sonner";
import { createCheckoutSession } from "../app/actions/create-checkout-session";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { ArrowRight } from "lucide-react";

export function SubscribeButton({
  planId,
  children,
}: {
  planId: string;
  children?: React.ReactNode;
}) {
  const router = useRouter();
  const [isPending, startTransition] = useTransition();

  const handleClick = () => {
    startTransition(async () => {
      console.log("Subscribe button clicked for plan:", planId);
      const result = await createCheckoutSession(planId);
      console.log("Checkout session result:", result);

      if (!result.ok) {
        console.error("Failed to create checkout session:", result.message);
        toast.error(result.message);
        return;
      }

      console.log("Redirecting to:", result.url);
      // ✅ redirect based on returned URL
      router.push(result.url!);
    });
  };

  return (
    <Dialog>
      <DialogTrigger asChild>
        {children || (
          <Button type="submit">
            Subscribe <ArrowRight />
          </Button>
        )}
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
