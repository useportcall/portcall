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
import { Select, SelectContent, SelectItem } from "@/components/ui/select";
import { useUpdatePlanFeature } from "@/hooks";
import { PlanFeature } from "@/models/plan-feature";
import { SelectTrigger } from "@radix-ui/react-select";
import { Calendar, ChevronDown } from "lucide-react";
import { useState } from "react";

export function MutableMeteredFeatureInterval({
  planFeature,
}: {
  planFeature: PlanFeature;
}) {
  const [value, setValue] = useState<string>(planFeature.interval);
  const intervals = ["per_use", "week", "month", "year", "no_reset"];

  const { mutate: update } = useUpdatePlanFeature(planFeature.id);

  return (
    <Select
      onValueChange={(val) => {
        setValue(val);
        update({ interval: val });
      }}
      value={value}
    >
      <SelectTrigger asChild>
        <button className="cursor-pointer font-medium text-sm self-center hover:bg-accent w-20 overflow-hidden text-ellipsis whitespace-nowrap rounded px-1 py-0.5 flex items-center justify-between">
          {intervalToTitle(value)}
          <ChevronDown className="size-3" />
        </button>
      </SelectTrigger>
      <SelectContent>
        {intervals.map((interval) => (
          <SelectItem
            key={interval}
            value={interval}
            onSelect={() => {
              setValue(interval);
              update({ interval });
            }}
          >
            {intervalToTitle(interval)}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}

export function ResetIntervalSelect({
  value,
  onChange,
}: {
  value: string;
  onChange: (value: string) => void;
}) {
  const [open, setOpen] = useState(false);
  const intervals = ["per_use", "week", "month", "year", "no_reset"];
  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <div>
          <Button
            type="button"
            variant="ghost"
            size="sm"
            className="w-full text-left flex justify-start border h-6 text-xs"
          >
            <Calendar className="size-3" />
            {intervalToTitle(value)}
          </Button>
        </div>
      </PopoverTrigger>
      <PopoverContent className="w-[200px] p-0" align="start" side="left">
        <Command>
          <CommandInput placeholder="Find interval..." />
          <CommandList>
            <CommandEmpty>No results found.</CommandEmpty>
            <CommandGroup>
              {intervals.map((interval) => (
                <CommandItem
                  key={interval}
                  onSelect={() => {
                    onChange(interval);
                    setOpen(false);
                  }}
                >
                  {intervalToTitle(interval)}
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
}

function intervalToTitle(interval: string) {
  switch (interval) {
    case "per_use":
      return "per use";
    case "week":
      return "weekly";
    case "month":
      return "monthly";
    case "year":
      return "yearly";
    case "no_reset":
      return "never";
    default:
      return "unknown";
  }
}
