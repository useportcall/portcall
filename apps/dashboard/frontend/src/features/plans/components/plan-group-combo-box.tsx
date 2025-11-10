import { Button } from "@/components/ui/button";
import { Filter } from "lucide-react";

import {
  Command,
  CommandEmpty,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import { Popover, PopoverContent } from "@/components/ui/popover";
import { useListPlanGroups } from "@/hooks";
import { PopoverTrigger } from "@radix-ui/react-popover";
import { useMemo } from "react";

export function PlanGroupComboBox({
  value,
  onChange,
}: {
  value: string;
  onChange: (value: string) => void;
}) {
  const { data: planGroups } = useListPlanGroups();

  const selected = useMemo(() => {
    if (value === "*") return "All groups";
    return (
      planGroups?.data.find((group) => group.id === value)?.name ?? "error"
    );
  }, [planGroups, value]);

  if (!planGroups) return null;

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button size="sm" variant={"outline"}>
          <Filter className="size-4" />
          {selected}
        </Button>
      </PopoverTrigger>
      <PopoverContent>
        <Command>
          <CommandInput placeholder="Search plan groups..." />
          <CommandList>
            <CommandEmpty />
            {planGroups.data.map((group) => (
              <CommandItem key={group.id} onSelect={() => onChange(group.id)}>
                {group.name}
              </CommandItem>
            ))}
            <CommandItem onSelect={() => onChange("*")}>
              Filter by group
            </CommandItem>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
}
