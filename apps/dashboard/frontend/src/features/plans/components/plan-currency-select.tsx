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
import { useRetrievePlan, useUpdatePlan } from "@/hooks";
import { Coins } from "lucide-react";
import { useState } from "react";

export function PlanCurrencySelect({ id }: { id: string }) {
  const [open, setOpen] = useState(false);
  const currencies = ["USD", "EUR", "GBP", "AUD", "CAD", "JPY"];

  const { data: plan } = useRetrievePlan(id);
  const { mutate: updatePlan } = useUpdatePlan(id);
  const currency = plan?.data?.currency ?? "Currency";

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <div>
          <Button
            variant="ghost"
            size="sm"
            className="w-full text-left flex justify-start"
            disabled={!plan?.data}
          >
            <Coins className="h-4 w-4" />
            {currency}
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
