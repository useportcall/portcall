import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { TableCell, TableRow } from "@/components/ui/table";
import { Invoice } from "@/models/invoice";
import { Eye } from "lucide-react";
import { Link } from "react-router-dom";

export function InvoiceListRow(props: {
  invoice: Invoice;
  locale: string;
  viewInvoiceLabel: string;
  pdfLabel: string;
}) {
  const { invoice } = props;
  return (
    <TableRow>
      <TableCell className="truncate whitespace-nowrap overflow-ellipsis"><h4 className="font-semibold overflow-ellipsis">{invoice.invoice_number}</h4></TableCell>
      <TableCell className="truncate whitespace-nowrap overflow-ellipsis">{formatCurrency(invoice.total ?? 0, invoice.currency, props.locale)}</TableCell>
      <TableCell className="truncate whitespace-nowrap overflow-ellipsis">
        <Link to={`/users/${invoice.recipient_id}`} target="_blank" className="flex flex-col">
          <span>{invoice.recipient_email}</span>
          <span className="text-xs text-muted-foreground italic">{invoice.recipient_id}</span>
        </Link>
      </TableCell>
      <TableCell className="truncate whitespace-nowrap overflow-ellipsis">
        <Badge variant={invoice.status === "paid" ? "success" : "outline"}>
          {invoice.status}
        </Badge>
      </TableCell>
      <TableCell className="truncate whitespace-nowrap overflow-ellipsis">{formatIssuedDate(invoice.created_at, props.locale)}</TableCell>
      <TableCell align="right">
        <div className="flex gap-2 justify-end">
          {invoice.pdf_url ? (
            <a href={invoice.pdf_url} target="_blank">
              <Badge variant="outline" className="cursor-pointer flex gap-2 w-fit">
                <span>{props.viewInvoiceLabel}</span><Eye className="size-3" />
              </Badge>
            </a>
          ) : (
            <Button disabled>{props.pdfLabel}</Button>
          )}
        </div>
      </TableCell>
    </TableRow>
  );
}

function formatCurrency(value: number, currency: string, locale = "en-US") {
  return (value / 100).toLocaleString(locale, { currency, style: "currency" });
}

function formatIssuedDate(value: string, locale = "en-US") {
  return new Date(value).toLocaleString(locale, {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}
