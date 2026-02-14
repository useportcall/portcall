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
import { useRetrievePlan, useUpdatePlan } from "@/hooks";
import { cn } from "@/lib/utils";
import { PopoverPortal } from "@radix-ui/react-popover";
import { ChevronDown, Heart } from "lucide-react";

export function FreeTrialComboBox({ id }: { id: string }) {
  const { data: plan } = useRetrievePlan(id);
  const { mutate: updatePlan } = useUpdatePlan(id);
  const trialDays = plan?.data?.trial_period_days || 0;

  const values = [0, 7, 14, 30];

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
          variant="ghost"
          size="sm"
          className="w-full text-left flex justify-between"
          disabled={!plan?.data}
        >
          <span className="flex gap-2 items-center justify-start">
            <Heart
              className={cn("h-4 w-4", {
                "text-teal-800": !!trialDays,
              })}
            />
            {!!trialDays && `${trialDays} day free trial`}
            {!trialDays && "No free trial"}
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
                updatePlan({ trial_period_days: Number(e.target.value) });
              }}
            />
            <CommandList>
              {values.map((val) => (
                <CommandItem
                  key={val}
                  onSelect={() => {
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
