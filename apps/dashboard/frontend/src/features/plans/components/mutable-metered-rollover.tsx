import { Select, SelectContent, SelectItem } from "@/components/ui/select";
import { useUpdatePlanFeatureForPlanItem } from "@/hooks";
import { PlanFeature } from "@/models/plan-feature";
import { SelectTrigger } from "@radix-ui/react-select";
import { ChevronDown } from "lucide-react";
import { useMemo } from "react";

export default function MutableMeteredRollover(props: {
  planFeature: PlanFeature;
  planId: string;
}) {
  const { mutate } = useUpdatePlanFeatureForPlanItem(
    props.planFeature.id,
    props.planId
  );

  const title = useMemo(() => {
    return toRolloverTitle(props.planFeature.rollover);
  }, [props.planFeature.rollover]);

  return (
    <Select
      value={props.planFeature.rollover?.toString() || "0"}
      onValueChange={(val) => {
        mutate({ rollover: parseInt(val) });
      }}
    >
      <SelectTrigger asChild>
        <button
          type="button"
          data-testid="metered-rollover-button"
          aria-label="Metered rollover"
          className="cursor-pointer text-sm w-16 overflow-hidden text-ellipsis whitespace-nowrap rounded py-1 flex items-center justify-between"
        >
          {title}
          <ChevronDown className="size-3" />
        </button>
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="0">No</SelectItem>
        <SelectItem value="-1">Yes</SelectItem>
      </SelectContent>
    </Select>
  );
}

function toRolloverTitle(value: number | null) {
  if (value === null || value === 0) {
    return "no";
  }

  if (value === -1) {
    return "yes";
  }

  return (
    "Up to " +
    value.toLocaleString("en-US", {
      style: "decimal",
    })
  );
}
