import { SelectItem } from "@/components/ui/select";
import { getProviderIcon } from "./provider-meta";

/** Render all available provider options inside a Select. */
export function ProviderSelectItems() {
  return (
    <>
      <SelectItem value="stripe">
        <div className="flex items-center gap-3 py-1">
          {getProviderIcon("stripe", "w-6 h-6")}
          <span className="font-medium">Stripe</span>
        </div>
      </SelectItem>
      <SelectItem value="braintree">
        <div className="flex items-center gap-3 py-1">
          {getProviderIcon("braintree", "w-6 h-6")}
          <span className="font-medium">Braintree</span>
        </div>
      </SelectItem>
      <SelectItem value="local">
        <div className="flex items-center gap-3 py-1">
          {getProviderIcon("local", "w-6 h-6")}
          <span className="font-medium">Mock Provider</span>
        </div>
      </SelectItem>
    </>
  );
}
