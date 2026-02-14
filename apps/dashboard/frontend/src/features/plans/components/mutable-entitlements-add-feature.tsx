import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
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
import {
  useCreateBasicFeature,
  useCreatePlanFeature,
  useListFeatures,
} from "@/hooks";
import { PlanFeature } from "@/models/plan-feature";
import { PopoverPortal } from "@radix-ui/react-popover";
import { Plus } from "lucide-react";
import { useMemo, useState } from "react";
import { toSnakeCase } from "../utils";
import { getUnlinkedFeatures } from "./mutable-entitlements-utils";

export function AddFeatureButton({
  planFeatures,
  planId,
}: {
  planFeatures: PlanFeature[];
  planId: string;
}) {
  const [open, setOpen] = useState(false);
  const [value, setValue] = useState("");
  const { mutateAsync: createFeature } = useCreateBasicFeature();
  const { mutate: addPlanFeature } = useCreatePlanFeature(planId);
  const { data: features } = useListFeatures({ isMetered: false });
  const filtered = useMemo(
    () => getUnlinkedFeatures(planFeatures, features?.data),
    [features?.data, planFeatures],
  );

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          type="button"
          data-testid="add-feature-button"
          aria-label="Add included feature"
          size="icon"
          variant="ghost"
          className="size-[22px] rounded-md hover:bg-accent"
        >
          <Plus />
        </Button>
      </PopoverTrigger>
      <PopoverPortal>
        <PopoverContent align="end" side="right">
          <Command>
            <CommandInput
              data-testid="feature-name-input"
              value={value}
              onValueChange={(next) => setValue(toSnakeCase(next))}
              placeholder="Find entitlements or create a new one"
              className="text-xs w-full outline-none px-2 py-1"
              onKeyDown={async (event) => {
                if (event.key !== "Enter") return;
                const featureID = toSnakeCase(value);
                if (!featureID) return;
                const existing = filtered.find(
                  (feature) => feature.id === featureID,
                );
                if (existing) {
                  addPlanFeature({
                    feature_id: featureID,
                    plan_id: planId,
                    interval: "no_reset",
                  });
                } else {
                  await createFeature({
                    feature_id: featureID,
                    plan_id: planId,
                    is_metered: false,
                  });
                }
                setValue("");
              }}
            />
            <CommandList>
              <CommandEmpty>No features found.</CommandEmpty>
              <CommandGroup>
                {filtered.map((feature) => (
                  <CommandItem
                    key={feature.id}
                    onSelect={() =>
                      addPlanFeature({
                        feature_id: feature.id,
                        plan_id: planId,
                        interval: "no_reset",
                      })
                    }
                  >
                    <Checkbox checked={false} />
                    {feature.id}
                  </CommandItem>
                ))}
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </PopoverPortal>
    </Popover>
  );
}
