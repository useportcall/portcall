import { Button } from "@/components/ui/button";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Plan } from "@/models/plan";
import { Hourglass } from "lucide-react";
import { useMemo, useState } from "react";
import { toIntervalText } from "../utils";
import { useUpdatePlan } from "@/hooks";

export function PlanIntervalSelect({ plan }: { plan: Plan }) {
  const [open, setOpen] = useState(false);
  const [value, setValue] = useState<string>(plan.interval);
  const intervals = ["month", "year", "no_reset"];

  const { mutate } = useUpdatePlan(plan.id);

  const label = useMemo(() => toIntervalText(value), [value]);

  return (
    <div>
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <div>
            <Button
              variant="ghost"
              size="sm"
              className="w-full text-left flex justify-start"
            >
              <Hourglass className="h-4 w-4" />
              {label}
            </Button>
          </div>
        </PopoverTrigger>
        <PopoverContent className="w-[200px] p-0" align="start" side="left">
          <Command>
            <CommandList>
              {intervals.length === 0 && (
                <CommandEmpty>No results found.</CommandEmpty>
              )}
              <CommandGroup>
                {intervals.map((interval) => (
                  <CommandItem
                    key={interval}
                    onSelect={() => {
                      setValue(interval);
                      setOpen(false);
                      mutate({ interval });
                    }}
                  >
                    {interval}
                  </CommandItem>
                ))}
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>
    </div>
  );
}
