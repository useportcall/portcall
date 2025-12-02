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
import { useGetQuote, useSendQuote, useVoidQuote } from "@/hooks/api/quotes";
import { useMemo, useState } from "react";

export function SendQuoteButton(props: { id: string }) {
  const { data: quote } = useGetQuote(props.id);
  const { mutateAsync: sendQuote } = useSendQuote(props.id);

  const [open, setOpen] = useState(false);
  const [sent, setSent] = useState(false);

  const buttonTitle = useMemo(() => {
    if (quote.data.status === "sent") {
      return "Resend";
    }
    return "Send";
  }, [quote]);

  const isDisabled = useMemo(() => {
    return ["accepted", "voided", "declined"].includes(quote.data.status);
  }, [quote]);

  return (
    <Dialog
      open={open}
      onOpenChange={(value) => {
        if (!value) {
          setTimeout(() => {
            setSent(false);
          }, 300);
        }
        setOpen(value);
      }}
    >
      <DialogTrigger asChild>
        <Button size="sm" disabled={isDisabled}>
          {buttonTitle}
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
        {sent && (
          <DialogFooter>
            <DialogClose>
              <Button variant="outline">Close</Button>
            </DialogClose>
          </DialogFooter>
        )}
        {!sent && (
          <DialogFooter>
            <DialogClose>
              <Button variant="outline">Cancel</Button>
            </DialogClose>
            <Button
              onClick={async () => {
                await sendQuote({});
                setSent(true);
              }}
            >
              Confirm
            </Button>
          </DialogFooter>
        )}
      </DialogContent>
    </Dialog>
  );
}

export function VoidQuoteButton(props: { id: string }) {
  const { data: quote } = useGetQuote(props.id);
  const { mutateAsync: voidQuote } = useVoidQuote(props.id);

  const isDisabled = useMemo(() => {
    return ["accepted", "voided", "declined", "draft"].includes(
      quote.data.status
    );
  }, [quote]);

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline" size="sm" disabled={isDisabled}>
          Void
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Void Quote</DialogTitle>
          <DialogDescription>
            Are you sure you want to void this quote? This action cannot be
            reversed.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <DialogClose>
            <Button variant="outline">Cancel</Button>
          </DialogClose>
          <Button
            variant="destructive"
            onClick={async () => {
              await voidQuote({});
            }}
          >
            Confirm
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

export function PreviewQuoteButton(props: { id: string }) {
  const [open, setOpen] = useState(false);
  return (
    <>
      <Button variant="outline" size="sm" onClick={() => setOpen(true)}>
        Preview
      </Button>
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent
          showCloseButton={false}
          className="w-full min-w-[80%] h-[calc(100vh-2rem)] p-0 flex flex-col bg-white rounded-xl shadow-2xl border border-gray-300 overflow-hidden"
        >
          {/* Mac window header */}
          <div className="flex items-center h-10 px-4 bg-[#f5f5f7] border-b border-gray-200 relative">
            {/* Traffic lights */}
            <div className="flex space-x-2 mr-4">
              <span className="w-3 h-3 rounded-full bg-[#ff5f56] border border-[#e0443e] inline-block"></span>
              <span className="w-3 h-3 rounded-full bg-[#ffbd2e] border border-[#dea123] inline-block"></span>
              <span className="w-3 h-3 rounded-full bg-[#27c93f] border border-[#13a10e] inline-block"></span>
            </div>
            <span className="font-medium text-gray-700 text-sm">
              Quote Preview
            </span>
            <DialogClose asChild>
              <button
                className="absolute right-4 top-1/2 -translate-y-1/2 px-2 py-1 text-gray-500 hover:text-gray-700 rounded transition"
                aria-label="Close"
                onClick={() => setOpen(false)}
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth={2}
                  stroke="currentColor"
                  className="w-5 h-5"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
            </DialogClose>
          </div>
          {/* Content area */}
          <div className="flex-1 min-h-0 relative bg-white">
            <iframe
              src={`http://localhost:8095/quotes/${props.id}`}
              title="Quote Preview"
              className="w-full h-full border-0 min-h-0 rounded-b-xl"
              sandbox="allow-same-origin allow-scripts allow-popups allow-forms"
            />
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}
