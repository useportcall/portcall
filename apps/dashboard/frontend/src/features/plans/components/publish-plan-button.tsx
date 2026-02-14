import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { usePublishPlan } from "@/hooks";
import { Plan } from "@/models/plan";
import { useState } from "react";

export function PublishPlanButton({ plan }: { plan: Plan }) {
  const { mutateAsync: publish } = usePublishPlan(plan.id);
  const [open, setOpen] = useState(false);

  if (plan.status === "published") {
    return (
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogTrigger asChild>
          <Button size="sm" disabled>
            Published
          </Button>
        </DialogTrigger>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Plan published!</DialogTitle>
            <DialogDescription>
              This plan is now live and available to customers.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <DialogClose>
              <Button variant="outline">Close</Button>
            </DialogClose>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    );
  }

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button size="sm" data-testid="publish-plan-button">
          Save and publish
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Publish Plan</DialogTitle>
          <DialogDescription>
            Are you sure you want to publish this plan? This action cannot be
            reversed.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <DialogClose>
            <Button variant="outline">Cancel</Button>
          </DialogClose>
          <Button
            data-testid="confirm-publish-button"
            onClick={async () => {
              await publish({});
              setOpen(true);
            }}
          >
            Confirm
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
