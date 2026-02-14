import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { MeterUsageControls } from "@/features/billing-meters/components/meter-usage-controls";
import { ResetMeterDialog } from "@/features/billing-meters/components/reset-meter-dialog";
import { useBillingMeterControls } from "@/features/billing-meters/hooks/use-billing-meter-controls";
import { formatUsd } from "@/features/billing-meters/utils/format";
import { BillingMeter } from "@/models/billing-meter";
import { useTranslation } from "react-i18next";

export function UserBillingMeterCard({ meter, subscriptionId }: { meter: BillingMeter; subscriptionId: string }) {
  const { t } = useTranslation();
  const controls = useBillingMeterControls(subscriptionId, meter);

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm flex items-center justify-between">
          <Badge variant="outline">{meter.feature_name || meter.feature_id}</Badge>
          <span className="text-sm font-semibold text-primary">{formatUsd(meter.projected_cost)}</span>
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <MeterUsageControls
              usage={meter.usage}
              isEditing={controls.isEditing}
              editValue={controls.editValue}
              setEditValue={controls.setEditValue}
              onStartEdit={() => controls.setIsEditing(true)}
              onCancelEdit={controls.cancelEdit}
              onSetUsage={controls.setUsage}
              onIncrement={controls.increment}
              onDecrement={controls.decrement}
              isUpdating={controls.isUpdating}
            />
          </div>
          <ResetMeterDialog disabled={false} isResetting={controls.isResetting} onReset={controls.reset} />
        </div>
        <p className="text-xs text-muted-foreground">
          {meter.pricing_model === "unit"
            ? t("views.user.billing_meters.pricing_per_unit", { amount: formatUsd(meter.unit_amount) })
            : t(`views.user.billing_meters.pricing_${meter.pricing_model}`)}
        </p>
      </CardContent>
    </Card>
  );
}
