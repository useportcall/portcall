import { TableCell } from "@/components/ui/table";
import {
    Tooltip,
    TooltipContent,
    TooltipTrigger,
} from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";
import { useCopyToClipboard } from "@/hooks/use-copy-to-clipboard";

interface CopyableIdCellProps {
    title: string;
    id: string;
    className?: string;
}

/**
 * Reusable table cell component for displaying a title and copyable ID.
 * Used across Plans, Quotes, and Users tables.
 */
export function CopyableIdCell({ title, id, className }: CopyableIdCellProps) {
    const { open, setOpen, copied, copyToClipboard, ref } = useCopyToClipboard();

    return (
        <TableCell className={cn("w-80", className)}>
            <div>
                <h4 className="font-semibold overflow-ellipsis">{title}</h4>
                <Tooltip
                    delayDuration={100}
                    open={open || copied}
                    onOpenChange={(o) => {
                        if (ref.current) return;
                        setOpen(o);
                    }}
                >
                    <TooltipTrigger asChild>
                        <span
                            className={cn(
                                "text-xs text-muted-foreground italic transition-colors hover:bg-accent cursor-pointer"
                            )}
                            onClick={(e) => {
                                e.stopPropagation();
                                copyToClipboard(id);
                            }}
                        >
                            {id}
                        </span>
                    </TooltipTrigger>
                    <TooltipContent className="p-1">
                        {copied ? "copied!" : "copy"}
                    </TooltipContent>
                </Tooltip>
            </div>
        </TableCell>
    );
}
