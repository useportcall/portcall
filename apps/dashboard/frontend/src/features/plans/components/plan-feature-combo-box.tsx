import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Feature } from "@/models/feature";
import { ChevronDown } from "lucide-react";
import React, { useState } from "react";
import {
  useCreateFeature,
  useListFeatures,
  useUpdatePlanFeature,
} from "@/hooks";
import { useParams } from "react-router-dom";
import { PlanFeature } from "@/models/plan-feature";
import { toSnakeCase } from "../utils";

export function PlanFeatureComboBox({
  planFeature,
}: {
  planFeature?: PlanFeature;
}) {
  const [open, setOpen] = React.useState(false);
  const [value, setValue] = useState("");

  const { data: features } = useListFeatures({ isMetered: true });

  const { mutateAsync: addPlanFeature } = useCreateFeature({ isMetered: true });

  const { id } = useParams();

  if (!features?.data) {
    return null;
  }

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <button className="cursor-pointer font-medium text-sm self-center hover:bg-accent w-20 rounded px-1 py-0.5 flex items-center justify-between">
          <span className="overflow-hidden text-ellipsis whitespace-nowrap">
            {planFeature?.feature.id || "Select feature"}
          </span>
          <ChevronDown className="size-3 shrink-0" />
        </button>
      </PopoverTrigger>
      <PopoverContent className="w-[200px] p-0" align="start">
        <Command>
          <CommandInput
            value={value}
            onValueChange={(v) => setValue(toSnakeCase(v))}
            placeholder="Search features..."
            onKeyDown={async (e) => {
              if (e.key === "Enter") {
                if (
                  !features.data.find((f) => f.id === e.currentTarget.value)
                ) {
                  await addPlanFeature({
                    plan_id: id,
                    feature_id: toSnakeCase(value),
                    is_metered: true,
                    plan_feature_id: planFeature?.id,
                  });

                  setOpen(false);
                }
              }
            }}
          />
          <CommandList>
            <CommandEmpty>Press "Enter" to add</CommandEmpty>
            <CommandGroup>
              {features.data.map((f) => (
                <PlanFeatureCommandItem
                  key={f.id}
                  planFeature={planFeature}
                  feature={f}
                />
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
}

function PlanFeatureCommandItem({
  planFeature,
  feature,
}: {
  planFeature?: PlanFeature;
  feature: Feature;
}) {
  const { mutate } = useUpdatePlanFeature(planFeature?.id);

  return (
    <CommandItem
      key={feature.id}
      onSelect={() => mutate({ feature_id: feature.id })}
    >
      {feature.id}
    </CommandItem>
  );
}
