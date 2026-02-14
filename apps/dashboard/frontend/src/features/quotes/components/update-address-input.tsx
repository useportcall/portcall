import { Input } from "@/components/ui/input";
import { useUpdateAddress } from "@/hooks";
import { Address } from "@/models/address";
import { useId } from "react";
import { useSyncedInput } from "./use-synced-input";

export function UpdateAddressInput({
  k,
  address,
  placeholder,
}: {
  k: keyof Address;
  address: Address;
  placeholder?: string;
}) {
  const id = useId();
  const { mutateAsync } = useUpdateAddress(address.id);
  const { currentValue, setCurrentValue } = useSyncedInput(address[k]);

  return (
    <div className="flex flex-col mt-2 gap-2 items-start text-sm w-full">
      <Input
        id={id}
        type="text"
        className="w-full text-sm"
        value={currentValue}
        onChange={(e) => setCurrentValue(e.target.value)}
        onBlur={async () => {
          if (currentValue !== address[k]) {
            await mutateAsync({ [k]: currentValue });
          }
        }}
        placeholder={placeholder || `Enter ${k}`}
      />
    </div>
  );
}
