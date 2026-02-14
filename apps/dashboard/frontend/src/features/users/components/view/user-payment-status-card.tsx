import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { Invoice } from "@/models/invoice";
import { cn } from "@/lib/utils";
import { useTranslation } from "react-i18next";

type Props = {
  subscriptionStatus: string;
  latestInvoice?: Invoice;
};

type Health = { label: string; hint: string; className: string };

export function UserPaymentStatusCard({ subscriptionStatus, latestInvoice }: Props) {
  const { t } = useTranslation();
  const health = deriveHealth(t, subscriptionStatus, latestInvoice?.status);

  return (
    <Card data-testid="payment-health-card">
      <CardContent className="pt-5 space-y-3">
        <div className="flex items-center justify-between">
          <p className="text-sm text-muted-foreground">{t("views.user.payment_health.title")}</p>
          <Badge data-testid="payment-health-badge" className={cn("border-none", health.className)}>{health.label}</Badge>
        </div>
        <p className="text-xs text-muted-foreground">{health.hint}</p>
        <div className="space-y-1">
          <p className="text-xs">
            <strong>{t("views.user.payment_health.subscription_label")}:</strong> {subscriptionStatus}
          </p>
          <p className="text-xs">
            <strong>{t("views.user.payment_health.latest_invoice_label")}:</strong>{" "}
            {latestInvoice
              ? `${latestInvoice.invoice_number} (${latestInvoice.status})`
              : t("views.user.payment_health.invoice_none")}
          </p>
        </div>
      </CardContent>
    </Card>
  );
}

function deriveHealth(t: (key: string) => string, subscriptionStatus: string, invoiceStatus?: string): Health {
  if (!invoiceStatus) {
    return {
      label: t("views.user.payment_health.no_invoices_label"),
      hint: t("views.user.payment_health.no_invoices_hint"),
      className: "bg-slate-100 text-slate-700",
    };
  }
  if (invoiceStatus === "uncollectible" || subscriptionStatus === "canceled") {
    return {
      label: t("views.user.payment_health.collection_failed_label"),
      hint: t("views.user.payment_health.collection_failed_hint"),
      className: "bg-red-100 text-red-800",
    };
  }
  if (invoiceStatus === "past_due") {
    return {
      label: t("views.user.payment_health.action_required_label"),
      hint: t("views.user.payment_health.action_required_hint"),
      className: "bg-amber-100 text-amber-800",
    };
  }
  if (invoiceStatus === "issued" || invoiceStatus === "pending") {
    return {
      label: t("views.user.payment_health.upcoming_charge_label"),
      hint: t("views.user.payment_health.upcoming_charge_hint"),
      className: "bg-slate-100 text-slate-700",
    };
  }
  return {
    label: t("views.user.payment_health.healthy_label"),
    hint: t("views.user.payment_health.healthy_hint"),
    className: "bg-emerald-400/20 text-emerald-700 dark:bg-emerald-400/15 dark:text-emerald-300",
  };
}
