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
import {
  useCreatePlanGroup,
  useListPlanGroups,
  useRetrievePlan,
  useUpdatePlan,
} from "@/hooks";
import { PopoverPortal } from "@radix-ui/react-popover";
import { Tag } from "lucide-react";
import React, { useState } from "react";

export function PlanGroupSelect({ id }: { id: string }) {
  const [open, setOpen] = useState(false);

  const { data: planGroups } = useListPlanGroups();

  const { data: plan } = useRetrievePlan(id);
  const updatePlan = useUpdatePlan(id);

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
  const planId = plan?.data?.id;
  const planGroupName = plan?.data?.plan_group?.name || "No group";

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="ghost"
          size="sm"
          className="w-full text-left flex justify-start"
          disabled={!planId}
        >
          <Tag className="h-4 w-4" />
          {planGroupName}
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
                  if (!planId) return;
                  // Handle search
                  if (
                    !planGroups.data.find((f) => f.id === e.currentTarget.value)
                  ) {
                    await createPlanGroup.mutateAsync({
                      name: e.currentTarget.value,
                      plan_id: planId,
                    });
                    setOpen(false);
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
