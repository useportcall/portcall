import { Button } from "@/components/ui/button";
import {
  Command,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { useCreatePlanGroup, useListPlanGroups, useUpdatePlan } from "@/hooks";
import { Plan } from "@/models/plan";
import { PopoverPortal } from "@radix-ui/react-popover";
import { Tag } from "lucide-react";
import React, { useState } from "react";

export function PlanGroupSelect({ plan }: { plan: Plan }) {
  const [open, setOpen] = useState(false);

  const { data: planGroups } = useListPlanGroups();

  const updatePlan = useUpdatePlan(plan.id);

  const createPlanGroup = useCreatePlanGroup();

  const [searchValue, setSearchValue] = React.useState("");

  // Filter commands based on search
  const filtered = React.useMemo(() => {
    if (!planGroups) return [];

    if (!searchValue) return planGroups.data;
    return planGroups.data.filter((group) =>
      group.name.toLowerCase().includes(searchValue.toLowerCase())
    );
  }, [planGroups, searchValue]);

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="ghost"
          size="sm"
          className="w-full text-left flex justify-start"
        >
          <Tag className="h-4 w-4" />
          {plan.plan_group?.name || "No group"}
        </Button>
      </PopoverTrigger>
      <PopoverPortal>
        <PopoverContent className="w-[200px] p-0" align="start" side="left">
          <Command>
            <CommandInput
              placeholder="Search groups..."
              value={open ? searchValue : ""}
              onValueChange={setSearchValue}
              onKeyDown={async (e) => {
                if (!planGroups || filtered.length !== 0) return;

                if (e.key === "Enter") {
                  // Handle search
                  if (
                    !planGroups.data.find((f) => f.id === e.currentTarget.value)
                  ) {
                    await createPlanGroup.mutateAsync({
                      name: e.currentTarget.value,
                    });
                  }
                }
              }}
            />
            <CommandList>
              {filtered.map((group) => (
                <CommandItem
                  key={group.id}
                  onSelect={async () => {
                    await updatePlan.mutateAsync({ plan_group_id: group.id });
                    setOpen(false);
                  }}
                >
                  {group.name}
                </CommandItem>
              ))}
            </CommandList>
          </Command>
        </PopoverContent>
      </PopoverPortal>
    </Popover>
  );
}
