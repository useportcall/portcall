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
import { useGetQuote, useSendQuote } from "@/hooks/api/quotes";
import { useState } from "react";
import { getSendButtonTitle, isSendDisabled } from "../utils/quote-ui";

export function SendQuoteButton({ id }: { id: string }) {
  const { data: quote } = useGetQuote(id);
  const { mutateAsync } = useSendQuote(id);
  const [open, setOpen] = useState(false);
  const [sent, setSent] = useState(false);

  return (
    <Dialog
      open={open}
      onOpenChange={(value) => {
        if (!value) setTimeout(() => setSent(false), 300);
        setOpen(value);
      }}
    >
      <DialogTrigger asChild>
        <Button size="sm" disabled={isSendDisabled(quote.data.status)}>
          {getSendButtonTitle(quote.data.status)}
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{sent ? "Sent" : "Send Quote"}</DialogTitle>
          <DialogDescription>
            {sent
              ? "This quote has been successfully sent to the client."
              : "Are you sure you want to send this quote? This action cannot be reversed."}
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <DialogClose>
            <Button variant="outline">{sent ? "Close" : "Cancel"}</Button>
          </DialogClose>
          {!sent && (
            <Button
              onClick={async () => {
                await mutateAsync({});
                setSent(true);
              }}
            >
              Confirm
            </Button>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
