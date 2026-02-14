import { useResetBillingMeter, useUpdateBillingMeter } from "@/hooks";
import { BillingMeter } from "@/models/billing-meter";
import { useState } from "react";

export function useBillingMeterControls(
  subscriptionId: string,
  meter: BillingMeter,
) {
  const [isEditing, setIsEditing] = useState(false);
  const [editValue, setEditValue] = useState(meter.usage.toString());
  const { mutateAsync: updateMeter, isPending: isUpdating } = useUpdateBillingMeter(
    subscriptionId,
    meter.feature_id,
  );
  const { mutateAsync: resetMeter, isPending: isResetting } = useResetBillingMeter(
    subscriptionId,
    meter.feature_id,
  );

  return {
    isEditing,
    setIsEditing,
    editValue,
    setEditValue,
    isUpdating,
    isResetting,
    increment: () => updateMeter({ increment: 1 }),
    decrement: () =>
      meter.usage > 0 ? updateMeter({ decrement: 1 }) : Promise.resolve(),
    setUsage: async () => {
      const usage = Number.parseInt(editValue, 10);
      if (!Number.isNaN(usage) && usage >= 0) {
        await updateMeter({ usage });
        setIsEditing(false);
      }
    },
    reset: () => resetMeter({}),
    cancelEdit: () => {
      setIsEditing(false);
      setEditValue(meter.usage.toString());
    },
  };
}
