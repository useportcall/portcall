import { cn } from "@/lib/utils";
import { CheckIcon, ChevronDown, Globe } from "lucide-react";
import { countries } from "country-data-list";
import { forwardRef, useEffect, useMemo, useState } from "react";
import { CircleFlag } from "react-circle-flags";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "./ui/command";
import { Popover, PopoverContent, PopoverTrigger } from "./ui/popover";

export interface Country {
  alpha2: string;
  alpha3: string;
  countryCallingCodes: string[];
  currencies: string[];
  emoji?: string;
  ioc: string;
  languages: string[];
  name: string;
  status: string;
}

interface CountryDropdownProps {
  options?: Country[];
  onChange?: (country: Country) => void;
  defaultValue?: string;
  disabled?: boolean;
  placeholder?: string;
  slim?: boolean;
  classNames?: string;
}

const defaultOptions = countries.all.filter(
  (country: Country) => country.emoji && country.status !== "deleted" && country.ioc !== "PRK",
);

export const CountryDropdown = forwardRef<HTMLButtonElement, CountryDropdownProps>(function CountryDropdown(
  { options = defaultOptions, onChange, defaultValue, disabled, placeholder = "Select a country", slim = false, ...props },
  ref,
) {
  const [open, setOpen] = useState(false);
  const [selectedCountry, setSelectedCountry] = useState<Country | undefined>();

  useEffect(() => {
    setSelectedCountry(defaultValue ? options.find((country) => country.alpha2 === defaultValue) : undefined);
  }, [defaultValue, options]);

  const filteredOptions = useMemo(() => options.filter((option) => option.name), [options]);

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger
        id="country_code"
        ref={ref}
        className={cn("flex h-9 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm shadow-sm ring-offset-background focus:outline-none focus:ring-1 focus:ring-ring disabled:cursor-not-allowed disabled:opacity-50", slim && "w-20")}
        disabled={disabled}
        {...props}
      >
        {selectedCountry ? (
          <CountryLabel alpha2={selectedCountry.alpha2} name={selectedCountry.name} slim={slim} />
        ) : slim ? (
          <Globe size={20} />
        ) : (
          <span>{placeholder}</span>
        )}
        <ChevronDown size={16} />
      </PopoverTrigger>
      <PopoverContent collisionPadding={10} side="bottom" align="start" className="min-w-[--radix-popper-anchor-width] p-0">
        <Command className="w-full max-h-[200px] sm:max-h-[270px]"><CommandList><div className="sticky top-0 z-10 bg-popover"><CommandInput placeholder="Search country..." /></div><CommandEmpty>No country found.</CommandEmpty><CommandGroup>
          {filteredOptions.map((option) => (
            <CommandItem
              className="flex items-center w-full gap-2"
              key={option.alpha2}
              onSelect={() => {
                setSelectedCountry(option);
                onChange?.(option);
                setOpen(false);
              }}
            >
              <CountryLabel alpha2={option.alpha2} name={option.name} slim={false} />
              <CheckIcon className={cn("ml-auto h-4 w-4 shrink-0", option.name === selectedCountry?.name ? "opacity-100" : "opacity-0")} />
            </CommandItem>
          ))}
        </CommandGroup></CommandList></Command>
      </PopoverContent>
    </Popover>
  );
});

function CountryLabel({ alpha2, name, slim }: { alpha2: string; name: string; slim: boolean }) {
  return (
    <div className="flex items-center flex-grow w-0 gap-2 overflow-hidden"><div className="inline-flex items-center justify-center w-5 h-5 shrink-0 overflow-hidden rounded-full"><CircleFlag countryCode={alpha2.toLowerCase()} height={20} /></div>{!slim && <span className="overflow-hidden text-ellipsis whitespace-nowrap">{name}</span>}</div>
  );
}
