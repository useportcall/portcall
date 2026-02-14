import { Input } from "@/components/ui/input";
import { useUpdateQuote } from "@/hooks/api/quotes";
import { Quote } from "@/models/quote";
import { useSyncedInput } from "./use-synced-input";

export function PreparedByEmailInput({ quote }: { quote: Quote }) {
  const { mutateAsync } = useUpdateQuote(quote.id);
  const { currentValue, setCurrentValue } = useSyncedInput(
    quote.prepared_by_email,
  );

  return (
    <div className="flex flex-col gap-2 px-2 text-sm">
      <p className="text-muted-foreground text-xs">
        Email shown in the "Prepared by" section. Leave empty to use your
        company email.
      </p>
      <Input
        type="email"
        className="w-full text-sm"
        value={currentValue}
        onChange={(e) => setCurrentValue(e.target.value)}
        onBlur={async () => {
          if (currentValue !== quote.prepared_by_email) {
            await mutateAsync({ prepared_by_email: currentValue });
          }
        }}
        placeholder="Enter email (optional)"
      />
    </div>
  );
}
