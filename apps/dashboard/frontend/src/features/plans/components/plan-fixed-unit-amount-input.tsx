import { useFirstFixedPlanItem, useUpdatePlanItem } from "@/hooks";
import { Plan } from "@/models/plan";
import { useMemo, useState } from "react";
import { currencyCodetoSymbol } from "../utils";
import { PlanItem } from "@/models/plan-item";

type Props = { plan: Plan };

export function PlanFixedUnitAmountInput({ plan }: Props) {
  const label = useMemo(
    () => currencyCodetoSymbol(plan.currency),
    [plan.currency]
  );

  const { planItem } = useFirstFixedPlanItem(plan.id);

  if (!planItem) {
    return (
      <div className="flex items-center mx-2 h-14 animate-pulse bg-accent rounded" />
    );
  }

  return (
    <div className="flex items-center px-2">
      <span className="text-5xl font-semibold">{label}</span>
      {!!planItem && <Input planItem={planItem} />}
    </div>
  );
}

function Input(props: { planItem: PlanItem }) {
  const [value, setValue] = useState<string>(
    (props.planItem.unit_amount / 100).toFixed(2)
  );

  const { mutate: updatePlanItem } = useUpdatePlanItem(props.planItem.id);

  return (
    <input
      type="text"
      placeholder="10.00"
      value={value}
      onChange={(e) => {
        if (isNaN(Number(e.target.value))) return;
        setValue(e.target.value);
      }}
      className="text-5xl font-bold outline-none"
      onBlur={() => {
        let numValue = parseFloat(value);

        if (isNaN(numValue)) {
          numValue = 0;
        }

        // ensure it's two decimals
        setValue(numValue.toFixed(2));
        updatePlanItem({ unit_amount: numValue * 100 });
      }}
    />
  );
}
