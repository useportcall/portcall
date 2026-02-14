import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { useGetQuote, useVoidQuote } from "@/hooks/api/quotes";
import { isVoidDisabled } from "../utils/quote-ui";

export function VoidQuoteButton({ id }: { id: string }) {
  const { data: quote } = useGetQuote(id);
  const { mutateAsync } = useVoidQuote(id);

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline" size="sm" disabled={isVoidDisabled(quote.data.status)}>
          Void
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Void Quote</DialogTitle>
          <DialogDescription>
            Are you sure you want to void this quote? This action cannot be reversed.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <DialogClose>
            <Button variant="outline">Cancel</Button>
          </DialogClose>
          <Button variant="destructive" onClick={() => mutateAsync({})}>
            Confirm
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
