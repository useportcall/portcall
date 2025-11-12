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
import { useUpdatePlan } from "@/hooks";
import { cn } from "@/lib/utils";
import { Plan } from "@/models/plan";
import { PopoverPortal } from "@radix-ui/react-popover";
import { ChevronDown, Heart } from "lucide-react";
import { useState } from "react";

export function FreeTrialComboBox({ plan }: { plan: Plan }) {
  const [value, setValue] = useState(plan.trial_period_days);

  const { mutate: updatePlan } = useUpdatePlan(plan.id);

  const values = [0, 7, 14, 30];

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
          variant="ghost"
          size="sm"
          className="w-full text-left flex justify-between"
        >
          <span className="flex gap-2 items-center justify-start">
            <Heart className={cn("h-4 w-4", { "text-teal-800": !!value })} />
            {!!value && `${value} day free trial`}
            {!value && "No free trial"}
          </span>{" "}
          <ChevronDown />
        </Button>
      </PopoverTrigger>
      <PopoverPortal>
        <PopoverContent className="w-[200px]" align="start" side="left">
          <Command>
            <CommandInput
              placeholder="Search..."
              onBlur={(e) => {
                if (e.target.value === "") return;
                setValue(Number(e.target.value));
                updatePlan({ trial_period_days: Number(e.target.value) });
              }}
            />
            <CommandList>
              {values.map((val) => (
                <CommandItem
                  key={val}
                  onSelect={() => {
                    setValue(val);
                    updatePlan({ trial_period_days: val });
                  }}
                >
                  {val} days
                </CommandItem>
              ))}
            </CommandList>
          </Command>
        </PopoverContent>
      </PopoverPortal>
    </Popover>
  );
}
