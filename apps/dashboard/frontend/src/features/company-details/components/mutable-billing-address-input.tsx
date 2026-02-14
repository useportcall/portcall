import { LabeledInput } from "@/components/base-view";
import { InputSaveIndicator, useInputSaveIndicator } from "@/components/input-save-indicator";
import { useUpdateAddress } from "@/hooks/api/addresses";
import { cn } from "@/lib/utils";
import type { Address } from "@/models/address";

export function MutableBillingAddressInput({ prop, address, placeholder }: { prop: keyof Address; address: Address; placeholder: string }) {
  const { mutateAsync } = useUpdateAddress(address.id);
  const { visible, showSaved } = useInputSaveIndicator();
  return <div className="flex flex-row gap-2"><LabeledInput id={String(prop)} helperText="" label={placeholder} defaultValue={address[prop]} placeholder={placeholder} className={cn("text-lg outline-none w-96")} onBlur={async (e) => { if (e.target.value === address[prop]) return; await mutateAsync({ [prop]: e.target.value }); showSaved(); }} /><InputSaveIndicator visible={visible} /></div>;
}
