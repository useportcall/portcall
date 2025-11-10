import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";

import { useListPlanFeatures, useUpdatePlanFeature } from "@/hooks";
import { PlanFeature } from "@/models/plan-feature";
import { PlanItem } from "@/models/plan-item";
import { useState } from "react";

export default function MutableMeteredLimit({
  planItem,
}: {
  planItem: PlanItem;
}) {
  const { data: planFeatures } = useListPlanFeatures({
    planItemId: planItem.id,
  });

  if (!planFeatures?.data?.length) return null;

  return (
    <Popover>
      <PopoverTrigger asChild>
        <button className="cursor-pointer text-sm font-medium w-16">
          {planFeatures.data[0].quota === -1 && "no limit"}
          {!planFeatures.data[0].quota && "no limit"}
          {planFeatures.data[0].quota &&
            planFeatures.data[0].quota > 0 &&
            planFeatures.data[0].quota.toLocaleString("en-US", {
              style: "decimal",
            })}
        </button>
      </PopoverTrigger>
      <PopoverContent>
        {planFeatures.data.map((pf) => (
          <PlanFeatureInput key={pf.id} planFeature={pf} />
        ))}
      </PopoverContent>
    </Popover>
  );
}

function PlanFeatureInput(props: { planFeature: PlanFeature }) {
  const [value, setValue] = useState<string>("");

  const { mutate: updatePlanFeature } = useUpdatePlanFeature(props.planFeature);

  return (
    <input
      type="text"
      placeholder="100"
      value={value}
      onChange={(e) => {
        if (e.target.value === "-") {
          setValue("-");
          return;
        }

        if (e.target.value === "") {
          setValue("");
          return;
        }

        const numValue = Number(e.target.value);

        if (isNaN(numValue)) {
          return;
        }

        setValue(numValue.toFixed(0));
      }}
      className="outline-none"
      onBlur={(e) => {
        let numValue: null | number = parseInt(e.target.value);
        if (isNaN(numValue)) {
          numValue = -1;
        }

        updatePlanFeature({ quota: numValue });
      }}
      onKeyDown={(e) => {
        if (e.key === "Enter") {
          e.preventDefault();

          let numValue: null | number = parseInt(value);
          if (isNaN(numValue)) {
            numValue = -1;
          }

          setValue(numValue.toString());
          updatePlanFeature({ quota: numValue });
        }
      }}
    />
  );
}
