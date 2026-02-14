import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { useUpdateQuote } from "@/hooks/api/quotes";
import { Quote } from "@/models/quote";
import { format } from "date-fns";
import { CalendarOff } from "lucide-react";
import { useState } from "react";

export function QuoteExpiryPicker({ quote }: { quote: Quote }) {
  const [date, setDate] = useState<Date | null>(
    quote.expires_at ? new Date(quote.expires_at) : null,
  );
  const { mutateAsync } = useUpdateQuote(quote.id);

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button variant="ghost" size="sm" className="w-full justify-start text-left font-normal">
          <CalendarOff className="size-4" />
          {date ? (
            <span>Quote expires {format(date, "PPP")}</span>
          ) : (
            <span className="text-muted-foreground">Add expiry date</span>
          )}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-auto p-0">
        <Calendar
          mode="single"
          captionLayout="dropdown"
          selected={date ?? undefined}
          onSelect={async (selectedDate) => {
            const nextDate = selectedDate ?? null;
            setDate(nextDate);
            await mutateAsync({ expires_at: nextDate?.toISOString() ?? null });
          }}
        />
      </PopoverContent>
    </Popover>
  );
}
