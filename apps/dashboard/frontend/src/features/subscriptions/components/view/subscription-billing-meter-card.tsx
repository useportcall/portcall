import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { MeterUsageControls } from "@/features/billing-meters/components/meter-usage-controls";
import { ResetMeterDialog } from "@/features/billing-meters/components/reset-meter-dialog";
import { useBillingMeterControls } from "@/features/billing-meters/hooks/use-billing-meter-controls";
import { formatShortDate, formatUsd, getPricingLabel } from "@/features/billing-meters/utils/format";
import { BillingMeter } from "@/models/billing-meter";

export function SubscriptionBillingMeterCard({ meter, subscriptionId }: { meter: BillingMeter; subscriptionId: string }) {
  const controls = useBillingMeterControls(subscriptionId, meter);

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Badge variant="outline">{meter.feature_name || meter.feature_id}</Badge>
            <span className="text-xs font-normal text-muted-foreground">{getPricingLabel(meter.pricing_model)}</span>
          </div>
          <Tooltip>
            <TooltipTrigger asChild><span className="text-base font-semibold text-primary">{formatUsd(meter.projected_cost)}</span></TooltipTrigger>
            <TooltipContent><p>Projected cost at current usage</p></TooltipContent>
          </Tooltip>
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
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
        {meter.pricing_model === "unit" && <p className="text-xs text-muted-foreground">{formatUsd(meter.unit_amount)} per unit</p>}
        {meter.next_reset_at && <p className="text-xs text-muted-foreground">Resets: {formatShortDate(meter.next_reset_at)}</p>}
      </CardContent>
    </Card>
  );
}
