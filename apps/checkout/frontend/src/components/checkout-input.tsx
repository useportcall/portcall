import { FormControl, FormItem, FormLabel } from "@/components/ui/form";
import { CardElement } from "@stripe/react-stripe-js";

export function CheckoutInput() {
  return (
    <div>
      <FormItem className="w-full flex-2 relative">
        <FormLabel>Card defails</FormLabel>
        <FormControl className="min-h-9">
          <CardElement
            options={{
              disableLink: true,
              hidePostalCode: true,
              style: { base: { fontSize: "14px" } },
            }}
            className="border bg-white py-2 px-3 rounded-md shadow-xs transition-[color,box-shadow] border-input"
          />
        </FormControl>
      </FormItem>
    </div>
  );
}
