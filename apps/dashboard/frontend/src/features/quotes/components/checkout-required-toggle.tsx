import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { useUpdateQuote } from "@/hooks/api/quotes";
import { Quote } from "@/models/quote";
import { CreditCard } from "lucide-react";
import { useEffect, useState } from "react";

export function CheckoutRequiredToggle({ quote }: { quote: Quote }) {
  const { mutate } = useUpdateQuote(quote.id);
  const [checked, setChecked] = useState(quote.direct_checkout_enabled);

  useEffect(() => {
    setChecked(quote.direct_checkout_enabled);
  }, [quote.direct_checkout_enabled]);

  return (
    <div className="px-2 flex items-center justify-start space-x-2">
      <CreditCard className="w-4 h-4" />
      <Label className="text-sm font-normal" htmlFor="checkout-required">
        Include checkout
      </Label>
      <Switch
        id="checkout-required"
        checked={checked}
        onCheckedChange={(newChecked) => {
          setChecked(newChecked);
          mutate({ direct_checkout_enabled: newChecked });
        }}
      />
    </div>
  );
}
