import { Button } from "@/components/ui/button";
import { useCreateQuote } from "@/hooks/api/quotes";
import { Loader2, Plus } from "lucide-react";

export default function QuoteListView() {
  return (
    <div className="w-full p-4 lg:p-10 flex flex-col gap-6">
      <div className="flex flex-row justify-between">
        <div className="flex flex-col space-y-8 w-full">
          <div className="flex justify-between w-full">
            <div className="flex flex-col justify-start space-y-2">
              <h1 className="text-lg md:text-xl font-semibold">Quotes</h1>
              <p className="text-sm text-slate-400">
                Manage existing quotes or add new quotes here.
              </p>
            </div>
            <div className="flex flex-row gap-4 items-end"></div>
          </div>
        </div>
      </div>
    </div>
  );
}

export function CreateQuoteButton() {
  const { mutate, isPending } = useCreateQuote();

  return (
    <Button
      size="sm"
      variant="outline"
      onClick={() => {
        mutate({});
      }}
    >
      {isPending && <Loader2 className="size-4 animate-spin" />}
      {!isPending && <Plus className="size-4" />}
      Add quote
    </Button>
  );
}
