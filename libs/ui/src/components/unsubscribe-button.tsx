"use client";

import { Loader2, X } from "lucide-react";
import { Button } from "./ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogTitle,
  DialogTrigger,
} from "./ui/dialog";
import { useTransition } from "react";
import { toast } from "sonner";
import cancelSubscription from "../api/cancel-subscription";

export function UnsubscribeButton({
  id,
  className,
}: {
  id: string;
  className?: string;
}) {
  const [isPending, startTransition] = useTransition();

  const handleClick = () => {
    startTransition(async () => {
      const result = await cancelSubscription(id);

      if (!result.ok) {
        toast.error(result.message);
        return;
      }
    });
  };

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button type="submit" className={className}>
          Unsubscribe
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogTitle>Cancel subscription</DialogTitle>
        <DialogDescription>
          Are you sure you want to unsubscribe?
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
