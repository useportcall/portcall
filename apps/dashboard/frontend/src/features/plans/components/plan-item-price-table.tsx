import {
  Table,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useUpdatePlanItem } from "@/hooks";
import { PlanItem, Tier } from "@/models/plan-item";
import { Plus, Trash2 } from "lucide-react";
import { useEffect, useState } from "react";

export function MutablePlanItemPrice({ planItem }: { planItem: PlanItem }) {
  const { mutateAsync } = useUpdatePlanItem(planItem.id);

  return (
    <Table className="text-sm">
      <TableHeader className="border-b">
        <TableHead className="text-sm text-slate-800 font-medium border-r">
          From
        </TableHead>
        <TableHead className="text-sm text-slate-800 border-r">To</TableHead>
        <TableHead className="text-sm text-slate-800">Price</TableHead>
        <TableHead className="text-sm text-slate-800 flex justify-end">
          <button
            disabled={planItem.pricing_model === "unit"}
            className="p-0 disabled:text-accent-foreground/50"
            onClick={async () => {
              const tiers = planItem.tiers || [];

              if (tiers.length > 0) {
                // Adjust the previous tier's end to make space for new tier
                const lastTier = tiers[tiers.length - 1];
                if (lastTier.end === -1) {
                  lastTier.end = lastTier.start + 1000;
                }
              }

              tiers.push({
                start: tiers.length > 0 ? tiers[tiers.length - 1].end + 1 : 0,
                end: -1,
                unit_amount: 0,
              });
              await mutateAsync({ tiers });
            }}
          >
            <Plus className="size-4" />
          </button>
        </TableHead>
      </TableHeader>
      <MutableAmount planItem={planItem} />
    </Table>
  );
}

function MutableAmount({ planItem }: { planItem: PlanItem }) {
  if (["unit", "fixed"].includes(planItem.pricing_model)) {
    return (
      <TierRow
        index={0}
        tier={{ start: 0, end: -1, unit_amount: planItem.unit_amount }}
        planItem={planItem}
      />
    );
  } else {
    return (
      planItem.tiers?.map((tier, index) => (
        <TierRow tier={tier} key={index} planItem={planItem} index={index} />
      )) || null
    );
  }
}

function TierRow({
  planItem,
  tier,
  index,
}: {
  planItem: PlanItem;
  tier: Tier;
  index: number;
}) {
  const { mutateAsync } = useUpdatePlanItem(planItem.id);

  const [endValue, setEndValue] = useState(tier.end);

  useEffect(() => {
    setEndValue(tier.end);
  }, [tier.end]);

  const [value, setValue] = useState((tier.unit_amount / 100).toFixed(2));

  const isLast = index === (planItem.tiers?.length || 1) - 1;

  return (
    <TableRow className="border-b-0">
      <TableCell className="text-sm border-r">
        <input
          type="text"
          value={tier.start}
          readOnly
          className="w-12 outline-none"
          onBlur={(e) => {
            let start = Number(e.target.value);

            if (index === 0) {
              start = 0;
            }

            const tiers = planItem.tiers?.map((t, i) => {
              if (i === index) {
                return { ...t, start };
              }
              return t;
            });
            mutateAsync({ tiers });
          }}
        />
      </TableCell>
      <TableCell className="text-sm border-r">
        <input
          type="text"
          value={tier.end === -1 ? "∞" : endValue}
          readOnly={tier.end === -1}
          className="w-12 outline-none"
          onChange={(e) => {
            const newValue = parseInt(e.target.value);
            if (!isNaN(newValue)) {
              setEndValue(newValue);
            }
          }}
          onBlur={(e) => {
            let end = Number(e.target.value);

            const tiers = planItem.tiers?.map((t, i) => {
              if (i === index) {
                if (isLast) {
                  return { ...t, end: -1 };
                }

                return { ...t, end };
              }

              if (i === index + 1) {
                return { ...t, start: end + 1 };
              }

              return t;
            });

            mutateAsync({ tiers });
          }}
        />
      </TableCell>
      <TableCell className="text-sm">
        $
        <input
          className="w-16 outline-none"
          type="text"
          value={value}
          onChange={(e) => {
            const newValue = Number(e.target.value);
            if (!isNaN(newValue)) {
              setValue(e.target.value);
            }
          }}
          onBlur={async (e) => {
            setValue(Number(e.target.value).toFixed(2));

            const tiers =
              planItem.tiers?.map((t, i) => {
                if (i === index) {
                  return {
                    ...t,
                    unit_amount: parseInt(e.target.value) * 100,
                  };
                }
                return t;
              }) ?? [];

            await mutateAsync({
              unit_amount: parseInt(e.target.value) * 100,
              tiers,
            });
          }}
        />
      </TableCell>
      <TableCell className="flex justify-end">
        <button
          className="p-0 disabled:text-accent-foreground/50 text-slate-800 hover:text-red-400"
          disabled={index === 0}
          onClick={async () => {
            const tiers = planItem.tiers?.filter((_, i) => i !== index) ?? [];

            if (tiers.length > 0) {
              // Adjust the previous tier's end to -1 if last tier is deleted
              if (index === tiers.length) {
                tiers[tiers.length - 1].end = -1;
              }
            }

            await mutateAsync({ tiers });
          }}
        >
          <Trash2 className="w-4 h-4  cursor-pointer" />
        </button>
      </TableCell>
    </TableRow>
  );
}
