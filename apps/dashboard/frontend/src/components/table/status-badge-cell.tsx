import { Badge } from "@/components/ui/badge";
import { TableCell } from "@/components/ui/table";

interface StatusBadgeCellProps {
    status: string;
    className?: string;
    variant?: "default" | "outline" | "secondary" | "destructive";
}

/**
 * Reusable table cell component for displaying status badges.
 * Used across Plans, Quotes, and Users tables.
 */
export function StatusBadgeCell({
    status,
    className,
    variant = "outline",
}: StatusBadgeCellProps) {
    return (
        <TableCell className={className || "w-40"}>
            <Badge variant={variant}>{status}</Badge>
        </TableCell>
    );
}
