import { cn } from "@/lib/utils";
import { PlanItem } from "@/models/plan-item";
import {
  useFirstFixedPlanItem,
  useRetrievePlan,
  useUpdatePlanItem,
} from "@/hooks";
import { useMemo, useState } from "react";
import { currencyCodetoSymbol } from "../utils";

export function PlanFixedUnitAmountInput({ id }: { id: string }) {
  const { data: plan } = useRetrievePlan(id);
  const { planItem } = useFirstFixedPlanItem(id);
  const currency = plan?.data?.currency;
  const discountPct = plan?.data?.discount_pct ?? 0;

  const label = useMemo(
    () => (currency ? currencyCodetoSymbol(currency) : ""),
    [currency]
  );
  const discountedAmount = useMemo(() => {
    if (!planItem || discountPct === 0) return null;
    return planItem.unit_amount - (planItem.unit_amount * discountPct) / 100;
  }, [planItem, discountPct]);

  if (!plan?.data || !planItem)
    return (
      <div className="flex items-center mx-2 h-14 animate-pulse bg-accent rounded" />
    );

  return (
    <div className="flex items-center px-2 relative">
      <span
        className={cn("text-5xl font-semibold", {
          "line-through text-black/30": !!plan.data.discount_pct,
        })}
      >
        {label}
      </span>
      <AmountInput
        planItem={planItem}
        planId={plan.data.id}
        strikethrough={!!plan.data.discount_pct}
      />
      {!!plan.data.discount_pct && (
        <div className="absolute right-2 top-2 flex flex-col items-end gap-1 animate-fade-in">
          <span className="text-xs bg-black text-white px-1.5 py-0.5 rounded font-medium">
            -{plan.data.discount_pct}%
          </span>
          <span className="text-xs font-medium text-black">
            starting at {label}
            {(discountedAmount! / 100).toFixed(2)}
          </span>
        </div>
      )}
    </div>
  );
}

function AmountInput({
  planItem,
  planId,
  strikethrough,
}: {
  planItem: PlanItem;
  planId: string;
  strikethrough?: boolean;
}) {
  const [value, setValue] = useState((planItem.unit_amount / 100).toFixed(2));
  const { mutate: updatePlanItem } = useUpdatePlanItem(planItem.id, planId);

  return (
    <input
      data-testid="plan-price-input"
      type="text"
      placeholder="10.00"
      value={value}
      onChange={(e) =>
        !isNaN(Number(e.target.value)) && setValue(e.target.value)
      }
      className={cn("text-5xl font-bold outline-none", {
        "line-through text-black/30": strikethrough,
      })}
      onBlur={() => {
        const n = Number.isNaN(parseFloat(value)) ? 0 : parseFloat(value);
        setValue(n.toFixed(2));
        updatePlanItem({ unit_amount: n * 100 });
      }}
    />
  );
}
