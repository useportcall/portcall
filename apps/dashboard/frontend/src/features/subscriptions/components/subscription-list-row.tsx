import { Badge } from "@/components/ui/badge";
import { TableCell, TableRow } from "@/components/ui/table";
import { CopyableIdCell } from "@/components/table/copyable-id-cell";
import { ActionsCell, ActionMenuItem } from "@/components/table/actions-cell";
import { Subscription } from "@/models/subscription";

export function SubscriptionListRow(props: {
  subscription: Subscription;
  locale: string;
  onNavigate: () => void;
  actions: ActionMenuItem[];
}) {
  const { subscription } = props;
  return (
    <TableRow
      onClick={props.onNavigate}
      className="cursor-pointer hover:bg-accent"
    >
      <TableCell className="w-20 py-1">
        <div>
          <h4 className="font-semibold overflow-ellipsis">{subscription.user?.email}</h4>
          <span className="text-xs text-muted-foreground italic">{subscription.user?.id}</span>
        </div>
      </TableCell>
      <TableCell className="max-w-20 truncate whitespace-nowrap overflow-ellipsis">
        <Badge className="w-fit" variant="outline">{subscription.plan?.name}</Badge>
      </TableCell>
      <TableCell className="max-w-20 truncate whitespace-nowrap overflow-ellipsis">
        <Badge variant={subscription.status === "active" ? "success" : "outline"}>
          {subscription.status}
        </Badge>
      </TableCell>
      <TableCell className="max-w-20 truncate whitespace-nowrap overflow-ellipsis">
        {formatDate(subscription.created_at, props.locale)}
      </TableCell>
      <TableCell className="max-w-20 truncate whitespace-nowrap overflow-ellipsis">
        {subscription.next_reset_at ? formatDate(subscription.next_reset_at, props.locale) : ""}
      </TableCell>
      <CopyableIdCell title="" id={subscription.id} className="max-w-20" />
      <ActionsCell actions={props.actions} />
    </TableRow>
  );
}

function formatDate(value: string, locale = "en-US") {
  return new Date(value).toLocaleDateString(locale, {
    year: "2-digit",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  });
}
