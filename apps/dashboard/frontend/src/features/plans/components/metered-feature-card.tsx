import { Card, CardContent, CardHeader } from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useDeletePlanItem, useUpdatePlanItem } from "@/hooks";
import { MeteredPlanItem } from "@/models/plan-item";
import { MoreVertical } from "lucide-react";
import { useState } from "react";
import MutableMeteredLimit from "./mutable-metered-limit";
import MutableMeteredRollover from "./mutable-metered-rollover";
import { PricingModelSelect } from "./mutable-price";
import { MutableMeteredFeatureInterval } from "./mutable-reset";
import { PlanFeatureComboBox } from "./plan-feature-combo-box";
import { MutablePlanItemPrice } from "./plan-item-price-table";

export function PlanItemCard({ planItem }: { planItem: MeteredPlanItem }) {
  const updatePlanItem = useUpdatePlanItem(planItem.id);
  const deletePlanItem = useDeletePlanItem(planItem.id);

  if (!planItem?.features?.length) return null;

  return (
    <Card>
      <CardHeader className="flex flex-col lg:flex-row justify-between items-center gap-4">
        <div className="self-start flex-col justify-evenly gap-4 max-w-[200px]">
          <p className="text-xs self-start  text-slate-600">Name</p>
          <input
            type="text"
            className="font-semibold w-full outline-none"
            placeholder="API Calls"
            defaultValue={planItem.public_title}
            onBlur={(e) => {
              if (e.target.value === planItem.public_title) return;
              return updatePlanItem.mutateAsync({
                public_title: e.target.value,
              });
            }}
          />
          <MutableMultiline
            className="text-xs text-slate-600 w-full resize-none outline-none"
            placeholder="No description"
            value={planItem.public_description}
            saveFn={(value) => {
              return updatePlanItem.mutateAsync({ public_description: value });
            }}
          />
        </div>
        <div className="flex flex-wrap lg:flex-nowrap gap-6">
          <div className="flex flex-col gap-1">
            <p className="text-xs self-center text-slate-600">Resets</p>
            <MutableMeteredFeatureInterval planFeature={planItem.feature} />
          </div>
          <div className="flex flex-col gap-1">
            <p className="text-xs self-center text-slate-600">
              Metered feature
            </p>
            <PlanFeatureComboBox planFeature={planItem.feature} />
          </div>
          <div className="flex flex-col gap-1">
            <p className="text-xs self-center  text-slate-600">Pricing</p>
            <PricingModelSelect
              value={planItem.pricing_model}
              onChange={(pricing_model) => {
                updatePlanItem.mutateAsync({ pricing_model });
              }}
            />
          </div>
          <div className="flex flex-col gap-1">
            <p className="text-xs self-center  text-slate-600">Rollover</p>
            <MutableMeteredRollover planFeature={planItem.feature} />
          </div>
          <div className="flex flex-col gap-1">
            <p className="text-xs self-center text-slate-600">Limit</p>
            <MutableMeteredLimit planFeature={planItem.feature} />
          </div>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <button>
                <MoreVertical className="w-6 h-6 p-1 hover:bg-slate-100 rounded-sm" />
              </button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem
                onClick={async () => {
                  // Handle delete action
                  await deletePlanItem.mutateAsync({});
                }}
              >
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </CardHeader>
      <CardContent>
        <MutablePlanItemPrice planItem={planItem} />
      </CardContent>
    </Card>
  );
}

function MutableMultiline({
  value,
  placeholder,
  saveFn,
  className,
}: {
  value: string;
  placeholder?: string;
  saveFn: (value: string) => Promise<any>;
  className?: string;
}) {
  const [state, setState] = useState(value);

  return (
    <textarea
      className={className}
      placeholder={placeholder}
      value={state}
      onChange={(e) => setState(e.target.value)}
      onBlur={(e) => {
        if (e.target.value !== value) {
          saveFn(e.target.value);
        }
      }}
    />
  );
}
