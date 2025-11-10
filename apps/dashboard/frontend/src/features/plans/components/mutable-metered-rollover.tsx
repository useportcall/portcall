import { Select, SelectContent, SelectItem } from "@/components/ui/select";
import { useUpdatePlanFeature } from "@/hooks";
import { PlanFeature } from "@/models/plan-feature";
import { SelectTrigger } from "@radix-ui/react-select";
import { ChevronDown } from "lucide-react";
import { useMemo } from "react";

export default function MutableMeteredRollover(props: {
  planFeature: PlanFeature;
}) {
  const { mutate: updatePlanFeature } = useUpdatePlanFeature(props.planFeature);

  const title = useMemo(() => {
    return toRolloverTitle(props.planFeature.rollover);
  }, [props.planFeature.rollover]);

  return (
    <Select
      value={props.planFeature.rollover?.toString() || "0"}
      onValueChange={(val) => {
        updatePlanFeature({ rollover: parseInt(val) });
      }}
    >
      <SelectTrigger asChild>
        <button className="cursor-pointer font-medium text-sm self-center hover:bg-accent w-16 overflow-hidden text-ellipsis whitespace-nowrap rounded px-1 py-0.5 flex items-center justify-between">
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
    return "No";
  }

  if (value === -1) {
    return "Yes";
  }

  return (
    "Up to " +
    value.toLocaleString("en-US", {
      style: "decimal",
    })
  );
}
