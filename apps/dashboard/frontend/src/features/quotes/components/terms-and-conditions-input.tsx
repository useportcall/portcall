import { Textarea } from "@/components/ui/textarea";
import { useUpdateQuote } from "@/hooks/api/quotes";
import { Quote } from "@/models/quote";
import { DEFAULT_TERMS } from "../constants";
import { useSyncedInput } from "./use-synced-input";

export function TermsAndConditionsInput({ quote }: { quote: Quote }) {
  const { mutateAsync } = useUpdateQuote(quote.id);
  const { currentValue, setCurrentValue } = useSyncedInput(
    quote.toc || DEFAULT_TERMS,
  );

  return (
    <div className="flex gap-1.5 px-2.5 py-1 items-start text-sm">
      <Textarea
        className="w-full min-h-[100px]"
        value={currentValue}
        onChange={(e) => setCurrentValue(e.target.value)}
        onBlur={async () => {
          if (currentValue !== quote.toc) {
            await mutateAsync({ toc: currentValue });
          }
        }}
      />
    </div>
  );
}
