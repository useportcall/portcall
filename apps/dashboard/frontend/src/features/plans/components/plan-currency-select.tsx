import { Button } from "@/components/ui/button";
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
import { useUpdatePlan } from "@/hooks";
import { Plan } from "@/models/plan";
import { Coins } from "lucide-react";
import { useState } from "react";

export function PlanCurrencySelect({ plan }: { plan: Plan }) {
  const [open, setOpen] = useState(false);
  const currencies = ["USD", "EUR", "GBP", "AUD", "CAD", "JPY"];

  const { mutate: updatePlan } = useUpdatePlan(plan.id);

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <div>
          <Button
            variant="ghost"
            size="sm"
            className="w-full text-left flex justify-start"
          >
            <Coins className="h-4 w-4" />
            {plan.currency}
          </Button>
        </div>
      </PopoverTrigger>
      <PopoverContent className="w-[200px] p-0" align="start" side="left">
        <Command>
          <CommandInput placeholder="Search currency..." />
          <CommandList>
            {currencies.length === 0 && (
              <CommandEmpty>No results found.</CommandEmpty>
            )}
            <CommandGroup>
              {currencies.map((currency) => (
                <CommandItem
                  key={currency}
                  onSelect={() => {
                    setOpen(false);
                    updatePlan({ currency });
                  }}
                >
                  {currency}
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
}
