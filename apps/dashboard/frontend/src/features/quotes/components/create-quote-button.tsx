import { Button } from "@/components/ui/button";
import { useCreateQuote } from "@/hooks/api/quotes";
import { Loader2, Plus } from "lucide-react";

export function CreateQuoteButton() {
  const { mutate, isPending } = useCreateQuote();

  return (
    <Button size="sm" variant="outline" onClick={() => mutate({})}>
      {isPending ? (
        <Loader2 className="size-4 animate-spin" />
      ) : (
        <Plus className="size-4" />
      )}
      Add quote
    </Button>
  );
}
