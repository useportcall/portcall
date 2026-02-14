import { Card, CardContent } from "@/components/ui/card";
import { useListBillingMetersBySubscription } from "@/hooks";
import { Subscription } from "@/models/subscription";
import { TrendingUp } from "lucide-react";
import { formatUsd } from "@/features/billing-meters/utils/format";
import { SubscriptionBillingMeterCard } from "./subscription-billing-meter-card";
import { useTranslation } from "react-i18next";

export function SubscriptionBillingMetersSection({ subscription }: { subscription: Subscription }) {
  const { t } = useTranslation();
  const { data: meters } = useListBillingMetersBySubscription(subscription.id);
  if (!meters) return null;

  const totalProjectedCost = meters.data.reduce((acc, meter) => acc + meter.projected_cost, 0);

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-semibold">{t("views.subscription.billing_meters.title")}</h2>
          <p className="text-sm text-muted-foreground">{t("views.subscription.billing_meters.description")}</p>
        </div>
        {totalProjectedCost > 0 && (
          <div className="text-right">
            <p className="text-xs text-muted-foreground">{t("views.subscription.billing_meters.projected_total")}</p>
            <p className="text-lg font-semibold text-primary">{formatUsd(totalProjectedCost)}</p>
          </div>
        )}
      </div>
      {!meters.data.length && (
        <Card className="border-dashed">
          <CardContent className="flex flex-col items-center justify-center py-8 text-center">
            <TrendingUp className="h-8 w-8 text-muted-foreground mb-2" />
            <p className="text-sm text-muted-foreground">{t("views.subscription.billing_meters.empty")}</p>
          </CardContent>
        </Card>
      )}
      <div className="grid gap-4 md:grid-cols-2">
        {meters.data.map((meter) => (
          <SubscriptionBillingMeterCard key={meter.feature_id} meter={meter} subscriptionId={subscription.id} />
        ))}
      </div>
    </div>
  );
}
