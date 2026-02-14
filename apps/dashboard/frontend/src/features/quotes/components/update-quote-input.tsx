import { Input } from "@/components/ui/input";
import { useUpdateQuote } from "@/hooks/api/quotes";
import { Quote } from "@/models/quote";
import { useId } from "react";
import { useSyncedInput } from "./use-synced-input";

type QuoteInputKey =
  | "recipient_email"
  | "recipient_name"
  | "recipient_title"
  | "company_name"
  | "prepared_by_email";

export function UpdateQuoteInput({
  k,
  quote,
  placeholder,
}: {
  k: QuoteInputKey;
  quote: Quote;
  placeholder?: string;
}) {
  const id = useId();
  const { mutateAsync } = useUpdateQuote(quote.id);
  const { currentValue, setCurrentValue } = useSyncedInput(quote[k]);

  return (
    <div className="flex flex-col mt-2 gap-2 items-start text-sm w-full">
      <Input
        id={id}
        type="text"
        className="w-full text-sm"
        value={currentValue}
        onChange={(e) => setCurrentValue(e.target.value)}
        onBlur={async () => {
          if (currentValue !== quote[k]) {
            await mutateAsync({ [k]: currentValue });
          }
        }}
        placeholder={placeholder || `Enter ${k.replace(/_/g, " ")}`}
      />
    </div>
  );
}
