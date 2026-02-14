import { Input } from "@/components/ui/input";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";

import { useUpdatePlanFeatureForPlanItem } from "@/hooks";
import { PlanFeature } from "@/models/plan-feature";
import { useState } from "react";
import {
  getQuotaTitle,
  parseQuotaInput,
  sanitizeQuotaInput,
} from "./mutable-metered-limit-utils";

export default function MutableMeteredLimit({
  planFeature,
  planId,
}: {
  planFeature: PlanFeature;
  planId: string;
}) {
  const [open, setOpen] = useState(false);
  const [value, setValue] = useState<string>("");

  const { mutate: updatePlanFeature } = useUpdatePlanFeatureForPlanItem(
    planFeature.id,
    planId
  );

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <button
          type="button"
          data-testid="metered-limit-button"
          aria-label="Metered limit"
          className="cursor-pointer text-sm w-16 py-1 text-start"
          onClick={() => setValue(planFeature.quota > 0 ? String(planFeature.quota) : "")}
        >
          {getQuotaTitle(planFeature.quota)}
        </button>
      </PopoverTrigger>
      <PopoverContent className="p-0">
        <Input
          data-testid="metered-limit-input"
          aria-label="Metered limit input"
          type="text"
          placeholder="100"
          value={value}
          onChange={(e) => {
            const next = sanitizeQuotaInput(e.target.value);
            if (next === null) return;
            setValue(next);
          }}
          onBlur={(e) => {
            const numValue = parseQuotaInput(e.target.value);
            updatePlanFeature({ quota: numValue });
          }}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              e.preventDefault();

              const numValue = parseQuotaInput(value);
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
