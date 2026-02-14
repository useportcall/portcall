import { Button } from "@/components/ui/button";
import { Dialog, DialogContent } from "@/components/ui/dialog";
import { useState } from "react";
import { PreviewDialogHeader } from "./preview-dialog-header";

export function PreviewQuoteButton({ id, quoteURL }: { id: string; quoteURL?: string | null }) {
  const [open, setOpen] = useState(false);
  const quoteAppUrl = (window as { __ENV__?: { QUOTE_APP_URL?: string } }).__ENV__
    ?.QUOTE_APP_URL || "http://localhost:8010";
  const previewURL = quoteURL || `${quoteAppUrl}/quotes/${id}`;

  return (
    <>
      <Button variant="outline" size="sm" onClick={() => setOpen(true)}>
        Preview
      </Button>
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent
          showCloseButton={false}
          className="w-full min-w-[80%] h-[calc(100vh-2rem)] p-0 flex flex-col bg-card rounded-xl shadow-2xl border border-border overflow-hidden"
        >
          <PreviewDialogHeader onClose={() => setOpen(false)} />
          <div className="flex-1 min-h-0 relative bg-card">
            <iframe
              src={previewURL}
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
