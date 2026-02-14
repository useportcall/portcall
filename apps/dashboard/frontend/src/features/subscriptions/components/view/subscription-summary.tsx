import { Badge } from "@/components/ui/badge";
import { Subscription } from "@/models/subscription";
import { Link } from "react-router-dom";
import { formatShortDate } from "@/features/billing-meters/utils/format";
import { useTranslation } from "react-i18next";

export function SubscriptionSummary({ subscription }: { subscription: Subscription }) {
  const { t } = useTranslation();
  return (
    <div className="space-y-4">
      <h2 className="text-xl font-semibold">{t("views.subscription.summary.title")}</h2>
      <SummaryRow label={t("views.subscription.summary.subscription_id")} value={subscription.id} href={`/subscriptions/${subscription.id}`} />
      <SummaryRow label={t("views.subscription.summary.related_plan")} value={subscription.plan.id} href={`/plans/${subscription.plan.id}`} />
      <div className="space-y-1">
        <p className="text-sm text-muted-foreground">{t("views.subscription.summary.status")}</p>
        <Badge variant={subscription.status === "active" ? "success" : "outline"}>
          {subscription.status}
        </Badge>
      </div>
      <SummaryRow label={t("views.subscription.summary.next_billing_date")} value={formatShortDate(subscription.next_reset_at)} />
    </div>
  );
}

function SummaryRow({ label, value, href }: { label: string; value: string; href?: string }) {
  return (
    <div className="space-y-1">
      <p className="text-sm text-muted-foreground">{label}</p>
      {href ? (
        <Link to={href} target="_blank" className="text-xs truncate whitespace-nowrap overflow-ellipsis hover:underline">{value}</Link>
      ) : (
        <p className="text-xs">{value}</p>
      )}
    </div>
  );
}
