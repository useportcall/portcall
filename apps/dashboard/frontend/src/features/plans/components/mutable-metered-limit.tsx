import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";

import { useUpdatePlanFeatureForPlanItem } from "@/hooks";
import { PlanFeature } from "@/models/plan-feature";
import { useState } from "react";

export default function MutableMeteredLimit({
  planFeature,
}: {
  planFeature: PlanFeature;
}) {
  const [open, setOpen] = useState(false);
  const [value, setValue] = useState<string>("");

  const { mutate: updatePlanFeature } = useUpdatePlanFeatureForPlanItem(
    planFeature.id
  );

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <button className="cursor-pointer text-sm font-medium w-16">
          {planFeature.quota === -1 && "no limit"}
          {!planFeature.quota && "no limit"}
          {planFeature.quota &&
            planFeature.quota > 0 &&
            planFeature.quota.toLocaleString("en-US", {
              style: "decimal",
            })}
        </button>
      </PopoverTrigger>
      <PopoverContent>
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
              setOpen(false);
            }
          }}
        />
      </PopoverContent>
    </Popover>
  );
}
