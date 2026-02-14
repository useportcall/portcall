import { DialogClose } from "@/components/ui/dialog";

export function PreviewDialogHeader({ onClose }: { onClose: () => void }) {
  return (
    <div className="flex items-center h-10 px-4 bg-muted border-b border-border relative">
      <div className="flex space-x-2 mr-4">
        <span className="w-3 h-3 rounded-full bg-[#ff5f56] border border-[#e0443e] inline-block" />
        <span className="w-3 h-3 rounded-full bg-[#ffbd2e] border border-[#dea123] inline-block" />
        <span className="w-3 h-3 rounded-full bg-[#27c93f] border border-[#13a10e] inline-block" />
      </div>
      <span className="font-medium text-foreground text-sm">Quote Preview</span>
      <DialogClose asChild>
        <button
          className="absolute right-4 top-1/2 -translate-y-1/2 px-2 py-1 text-muted-foreground hover:text-foreground rounded transition"
          aria-label="Close"
          onClick={onClose}
        >
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor" className="w-5 h-5">
            <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </DialogClose>
    </div>
  );
}
