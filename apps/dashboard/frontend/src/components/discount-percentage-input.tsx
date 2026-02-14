import { useUpdatePlan } from "@/hooks";
import { Plan } from "@/models/plan";
import { Percent } from "lucide-react";
import { useState } from "react";

export function DiscountPercentageInput(props: { plan: Plan }) {
  const [discountPct, setDiscountPct] = useState(props.plan.discount_pct);
  const { mutateAsync } = useUpdatePlan(props.plan.id);

  return (
    <div className="flex gap-1.5 px-2.5 py-1 items-center text-sm font-medium">
      <Percent className="size-4" />
      <input
        type="number"
        className="outline-none w-full"
        min={0}
        max={100}
        value={discountPct}
        onChange={(e) => setDiscountPct(Number(e.target.value))}
        onBlur={async () => {
          const value = Math.min(Math.max(discountPct, 0), 100);
          setDiscountPct(value);

          if (value !== props.plan.discount_pct) {
            await mutateAsync({ discount_pct: value });
          }
        }}
        placeholder="Enter discount percentage"
      />
    </div>
  );
}
