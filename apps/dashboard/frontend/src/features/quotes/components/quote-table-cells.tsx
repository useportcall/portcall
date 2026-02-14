import { TableCell } from "@/components/ui/table";
import { Quote } from "@/models/quote";
import { useMemo } from "react";
import { Link } from "react-router-dom";
import { useSendQuote, useVoidQuote } from "@/hooks/api/quotes";
import { formatQuoteFixedCharge } from "@/lib/format";
import { CopyableIdCell } from "@/components/table/copyable-id-cell";
import { StatusBadgeCell } from "@/components/table/status-badge-cell";
import { ActionsCell, ActionMenuItem } from "@/components/table/actions-cell";

export function QuoteFixedChargeTableCell(props: { quote: Quote }) {
  const fixedCharge = useMemo(
    () => formatQuoteFixedCharge(props.quote),
    [props.quote]
  );

  return <TableCell className="w-40">{fixedCharge}</TableCell>;
}

export function QuoteNameTableCell(props: { quote: Quote }) {
  return (
    <CopyableIdCell
      title={props.quote.plan?.name || "Untitled Quote"}
      id={props.quote.id}
    />
  );
}

export function QuoteStatusTableCell(props: { quote: Quote }) {
  return <StatusBadgeCell status={props.quote.status} />;
}

export function QuoteActionsTableCell({ quote }: { quote: Quote }) {
  const { mutateAsync: sendQuote, isPending: isSending } = useSendQuote(
    quote.id
  );
  const { mutate: voidQuote, isPending: isVoiding } = useVoidQuote(quote.id);

  const actions: ActionMenuItem[] = [
    {
      label: <Link to={`/quotes/${quote.id}`}>Edit</Link>,
      asChild: true,
    },
    {
      label: "Send",
      onClick: (e) => {
        e.stopPropagation();
        sendQuote({});
      },
      loading: isSending,
      disabled: isSending || quote.status !== "draft",
    },
    {
      label: "Void",
      onClick: (e) => {
        e.stopPropagation();
        voidQuote({});
      },
      loading: isVoiding,
      disabled: isVoiding || (quote.status !== "draft" && quote.status !== "sent"),
      className: "text-red-600 focus:text-red-700",
    },
  ];

  return <ActionsCell actions={actions} />;
}
