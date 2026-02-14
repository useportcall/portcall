import EmptyTable from "@/components/empty-table";
import { Badge } from "@/components/ui/badge";
import { useListSubscriptionInvoices } from "@/hooks";
import { Subscription } from "@/models/subscription";
import { formatUsd } from "@/features/billing-meters/utils/format";
import { useTranslation } from "react-i18next";

export function SubscriptionInvoiceList({ subscription }: { subscription: Subscription }) {
  const { t } = useTranslation();
  const { data: invoices } = useListSubscriptionInvoices(subscription.id);
  if (!invoices) return null;

  return (
    <div className="space-y-4">
      <h2 className="text-xl font-semibold">{t("views.subscription.invoices.title")}</h2>
      <div className="space-y-4">
        {!invoices.data.length && (
          <EmptyTable message={t("views.subscription.invoices.empty")} button="" />
        )}
        {invoices.data.map((invoice: any) => (
          <div key={invoice.id} className="rounded border p-2">
            <div className="flex justify-between items-start">
              <div className="flex flex-col items-start">
                <h3 className="font-semibold text-sm">{invoice.invoice_number}</h3>
                <span className="text-sm">{formatUsd(invoice.total)}</span>
              </div>
              <Badge variant="outline">{invoice.status}</Badge>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
