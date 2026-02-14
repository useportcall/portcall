import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogFooter, DialogTrigger } from "@/components/ui/dialog";
import { useListPlans } from "@/hooks";
import { useListConnections } from "@/hooks/api/connections";
import { useCreatePaymentLink } from "@/hooks/api/payment-links";
import { Link2, Loader2 } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { PaymentLinkDialogHeader } from "./payment-link-dialog-header";
import { PaymentLinkFormFields } from "./payment-link-form-fields";
import { PaymentLinkURLField } from "./payment-link-url-field";

type CreatePaymentLinkDialogProps = { userId: string; userEmail?: string };

export function CreatePaymentLinkDialog(props: CreatePaymentLinkDialogProps) {
  const { t } = useTranslation();
  const { userId, userEmail } = props;
  const [open, setOpen] = useState(false);
  const [selectedPlanID, setSelectedPlanID] = useState("");
  const [expiresInDays, setExpiresInDays] = useState("7");
  const [linkURL, setLinkURL] = useState("");
  const { data: plans, isLoading: isPlansLoading } = useListPlans("*");
  const { data: connections, isLoading: isConnectionsLoading } = useListConnections();
  const { mutateAsync: createPaymentLink, isPending } = useCreatePaymentLink();
  const publishedPlans = (plans?.data || []).filter((plan) => plan.status === "published");
  const hasPaymentConnection = (connections?.data?.length || 0) > 0;

  const onCreate = async () => {
    if (!selectedPlanID || !hasPaymentConnection) return;
    const days = Math.max(1, Math.min(90, Number(expiresInDays) || 7));
    const result = await createPaymentLink({
      plan_id: selectedPlanID,
      user_id: userId,
      user_email: userEmail,
      expires_at: new Date(Date.now() + days * 24 * 60 * 60 * 1000).toISOString(),
    });
    setLinkURL(result.data.url);
  };

  const closeDialog = () => {
    setOpen(false);
    setSelectedPlanID("");
    setExpiresInDays("7");
    setLinkURL("");
  };

  return (
    <Dialog open={open} onOpenChange={(isOpen) => (isOpen ? setOpen(true) : closeDialog())}>
      <DialogTrigger asChild>
        <Button variant="outline" className="gap-2" data-testid="payment-link-open">
          <Link2 className="size-4" />
          {t("views.user.payment_link.open_button")}
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-lg">
        <PaymentLinkDialogHeader />
        {linkURL ? (
          <PaymentLinkURLField linkURL={linkURL} />
        ) : hasPaymentConnection ? (
          <PaymentLinkFormFields
            selectedPlanID={selectedPlanID}
            expiresInDays={expiresInDays}
            plans={publishedPlans}
            onPlanChange={setSelectedPlanID}
            onDaysChange={setExpiresInDays}
          />
        ) : (
          <p className="rounded-md border border-amber-300 bg-amber-50 p-3 text-sm text-amber-900">
            Add a payment connection in Integrations before creating payment links.
          </p>
        )}
        <DialogFooter>
          {!linkURL && <Button variant="outline" onClick={closeDialog}>{t("common.cancel")}</Button>}
          <Button
            data-testid={linkURL ? "payment-link-done" : "payment-link-create-submit"}
            disabled={isPending || isPlansLoading || isConnectionsLoading || (!linkURL && (!selectedPlanID || !hasPaymentConnection))}
            onClick={linkURL ? closeDialog : onCreate}
          >
            {isPending && <Loader2 className="size-4 animate-spin mr-2" />}
            {linkURL ? t("views.user.payment_link.done_button") : t("views.user.payment_link.create_button")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
