import { Label } from "@/components/ui/label";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";

interface StatusFilterProps {
    label: string;
    value: string;
    options: { value: string; label: string }[];
    allLabel: string;
    onValueChange: (value: string) => void;
}

/**
 * Reusable status filter dropdown.
 * Used across Subscriptions, Invoices, and Quotes tables.
 */
export function StatusFilter(props: StatusFilterProps) {
    return (
        <div className="space-y-1">
            <Label className="text-xs text-muted-foreground">{props.label}</Label>
            <Select value={props.value} onValueChange={props.onValueChange}>
                <SelectTrigger className="w-[170px] h-9">
                    <SelectValue />
                </SelectTrigger>
                <SelectContent>
                    <SelectItem value="all">{props.allLabel}</SelectItem>
                    {props.options.map((option) => (
                        <SelectItem key={option.value} value={option.value}>
                            {option.label}
                        </SelectItem>
                    ))}
                </SelectContent>
            </Select>
        </div>
    );
}
