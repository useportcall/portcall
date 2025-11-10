import { Button } from "@/components/ui/button";
import {
  Command,
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
import { useCreateFeature, useListFeatures } from "@/hooks";
import { useState } from "react";

type Props = {
  value: string;
  setValue: (value: string) => void;
};

export function FeatureSelect({ value, setValue }: Props) {
  const [text, setText] = useState("");

  const { data: features } = useListFeatures({ isMetered: true });

  const { mutateAsync: addFeature } = useCreateFeature({ isMetered: true });

  if (!features) return null;

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button className="bg-black cursor-pointer px-2 h-6 text-xs">
          {value || "Select a feature"}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-[200px] p-0" align="start">
        <Command>
          <CommandInput
            value={text}
            onValueChange={setText}
            placeholder="Find or create feature..."
            onKeyDown={async (e) => {
              if (e.key === "Enter") {
                await addFeature({
                  feature_id: text,
                  is_metered: true,
                });
                setValue?.(text);
              }
            }}
          />
          <CommandList>
            <CommandGroup>
              {features.data.map((feature) => (
                <CommandItem
                  key={feature.id}
                  onSelect={() => setValue?.(feature.id)}
                >
                  {feature.id}
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
}
