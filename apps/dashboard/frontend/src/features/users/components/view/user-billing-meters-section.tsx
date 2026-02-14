import { Card, CardContent } from "@/components/ui/card";
import { useGetUserSubscription, useListBillingMetersByUser } from "@/hooks";
import { TrendingUp } from "lucide-react";
import { useParams } from "react-router-dom";
import { formatUsd } from "@/features/billing-meters/utils/format";
import { UserBillingMeterCard } from "./user-billing-meter-card";
import { useTranslation } from "react-i18next";

export function UserBillingMetersSection() {
  const { t } = useTranslation();
  const { id } = useParams();
  const { data: subscription } = useGetUserSubscription(id!);
  const { data: meters } = useListBillingMetersByUser(id!);

  if (!meters || !subscription?.data) return null;

  const totalProjectedCost = meters.data.reduce((acc, meter) => acc + meter.projected_cost, 0);

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-sm text-muted-foreground">{t("views.user.billing_meters.title")}</h2>
          <p className="text-xs text-muted-foreground">{t("views.user.billing_meters.description")}</p>
        </div>
        {totalProjectedCost > 0 && (
          <div className="text-right">
            <p className="text-xs text-muted-foreground">{t("views.user.billing_meters.projected_cost")}</p>
            <p className="text-lg font-semibold text-primary">{formatUsd(totalProjectedCost)}</p>
          </div>
        )}
      </div>
      {!meters.data.length && (
        <Card className="border-dashed">
          <CardContent className="flex flex-col items-center justify-center py-6 text-center">
            <TrendingUp className="h-6 w-6 text-muted-foreground mb-2" />
            <p className="text-sm text-muted-foreground">{t("views.user.billing_meters.empty")}</p>
          </CardContent>
        </Card>
      )}
      <div className="grid gap-3 md:grid-cols-2 lg:grid-cols-3">
        {meters.data.map((meter) => (
          <UserBillingMeterCard key={meter.feature_id} meter={meter} subscriptionId={subscription.data.id} />
        ))}
      </div>
    </div>
  );
}
