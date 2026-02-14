import { ReactNode } from "react";
import { Button } from "@/components/ui/button";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { TableCell } from "@/components/ui/table";
import { Loader2, MoreHorizontal } from "lucide-react";

export interface ActionMenuItem {
    label: string | ReactNode;
    onClick?: (e: React.MouseEvent) => void;
    disabled?: boolean;
    loading?: boolean;
    className?: string;
    href?: string;
    asChild?: boolean;
    children?: ReactNode;
}

interface ActionsCellProps {
    actions: ActionMenuItem[];
    align?: "left" | "center" | "right";
    className?: string;
}

/**
 * Reusable table cell component for action dropdown menus.
 * Used across Plans, Quotes, and Users tables.
 */
export function ActionsCell({
    actions,
    align = "right",
    className,
}: ActionsCellProps) {
    return (
        <TableCell
            className={className || "w-40"}
            align={align}
            onClick={(e) => e.stopPropagation()}
        >
            <DropdownMenu modal={false}>
                <DropdownMenuTrigger asChild>
                    <Button size="icon" variant="ghost" className="h-8 w-8 p-0">
                        <MoreHorizontal className="w-5 h-5" />
                    </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                    {actions.map((action, index) => {
                        if (action.children) {
                            return action.children;
                        }

                        return (
                            <DropdownMenuItem
                                key={index}
                                disabled={action.disabled || action.loading}
                                onClick={action.onClick}
                                className={action.className}
                                asChild={action.asChild}
                            >
                                {action.loading ? (
                                    <span className="flex items-center gap-2">
                                        <Loader2 className="animate-spin w-4 h-4" />
                                        {typeof action.label === "string"
                                            ? `${action.label}...`
                                            : action.label}
                                    </span>
                                ) : (
                                    action.label
                                )}
                            </DropdownMenuItem>
                        );
                    })}
                </DropdownMenuContent>
            </DropdownMenu>
        </TableCell>
    );
}
