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
import { useCancelSubscription } from "@/hooks";
import { Subscription } from "@/models/subscription";
import { Loader2 } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";

export function CancelUserSubscriptionButton({ subscription }: { subscription: Subscription }) {
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);
  const { mutateAsync, isPending } = useCancelSubscription(subscription.id, subscription.user_id);

  if (subscription.status !== "active") return null;

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild><Button variant="destructive">{t("views.user.subscription.cancel_action")}</Button></DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t("views.user.subscription.cancel_title")}</DialogTitle>
          <DialogDescription>{t("views.user.subscription.cancel_description")}</DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <DialogClose><Button variant="outline">{t("views.user.subscription.cancel_no")}</Button></DialogClose>
          <Button
            disabled={isPending}
            variant="destructive"
            onClick={async () => {
              await mutateAsync({});
              setOpen(false);
            }}
          >
            {t("views.user.subscription.cancel_yes")}
            {isPending && <Loader2 className="animate-spin size-3" />}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
