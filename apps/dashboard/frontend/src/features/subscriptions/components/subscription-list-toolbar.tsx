import { Input } from "@/components/ui/input";
import { StatusFilter } from "@/components/table/status-filter";

export function SubscriptionListToolbar(props: {
  search: string;
  status: string;
  onSearchChange: (value: string) => void;
  onStatusChange: (value: string) => void;
  searchPlaceholder: string;
  statusLabel: string;
  allStatusesLabel: string;
}) {
  return (
    <div className="flex w-full items-end gap-2">
      <div className="relative flex items-center gap-2 w-full max-w-xs">
        <Input
          type="text"
          placeholder={props.searchPlaceholder}
          value={props.search}
          onChange={(event) => props.onSearchChange(event.target.value)}
        />
      </div>
      <StatusFilter
        label={props.statusLabel}
        value={props.status}
        options={[
          { value: "active", label: "Active" },
          { value: "canceled", label: "Canceled" },
          { value: "trialing", label: "Trialing" },
          { value: "past_due", label: "Past due" },
        ]}
        allLabel={props.allStatusesLabel}
        onValueChange={props.onStatusChange}
      />
    </div>
  );
}
