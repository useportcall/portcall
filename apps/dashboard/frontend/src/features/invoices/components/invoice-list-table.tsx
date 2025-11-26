import EmptyTable from "@/components/empty-table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useListInvoices } from "@/hooks/api/invoices";
import { cn } from "@/lib/utils";
import { Eye } from "lucide-react";
import { Link } from "react-router-dom";

export function InvoiceListTable() {
  const { data: invoices } = useListInvoices();

  if (!invoices.data.length) {
    return <EmptyTable message="No invoice issued yet." />;
  }

  return (
    <Card className="p-2 animate-fade-in">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Invoice number</TableHead>
            <TableHead>Total</TableHead>
            <TableHead>Customer</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Issued</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {invoices.data.map((invoice) => (
            <TableRow key={invoice.id}>
              <TableCell className="truncate whitespace-nowrap overflow-ellipsis">
                <div>
                  <h4 className="font-semibold overflow-ellipsis">
                    {invoice.invoice_number}
                  </h4>
                </div>
              </TableCell>
              <TableCell className="truncate whitespace-nowrap overflow-ellipsis">
                {((invoice.total ?? 0) / 100).toLocaleString("en-US", {
                  currency: invoice.currency,
                  style: "currency",
                })}
              </TableCell>
              <TableCell className="truncate whitespace-nowrap overflow-ellipsis">
                <Link
                  to={`/users/${invoice.recipient_id}`}
                  target="_blank"
                  className="flex flex-col"
                >
                  <span> {invoice.recipient_email}</span>
                  <span className="text-xs text-slate-500 italic">
                    {invoice.recipient_id}
                  </span>
                </Link>
              </TableCell>

              <TableCell className="truncate whitespace-nowrap overflow-ellipsis">
                <Badge
                  className={cn(
                    invoice.status == "paid"
                      ? "bg-emerald-400/50 text-emerald-800"
                      : "bg-none"
                  )}
                >
                  {invoice.status}
                </Badge>
              </TableCell>
              <TableCell className="truncate whitespace-nowrap overflow-ellipsis">
                {/* should be dd/mm/yyyy hh:mm */}
                {new Date(invoice.created_at).toLocaleString("en-GB", {
                  day: "2-digit",
                  month: "2-digit",
                  year: "numeric",
                  hour: "2-digit",
                  minute: "2-digit",
                })}
              </TableCell>

              <TableCell align="right">
                <div className="flex gap-2 justify-end">
                  {/* <InvoiceEmail invoice={invoice} /> */}
                  {/* <a href={invoice.email_url} target="_blank">
                    <Badge
                      className={cn(
                        "bg-teal-800 cursor-pointer flex gap-2 w-fit",
                        {
                          "bg-teal-800/20": !invoice.email_url,
                        }
                      )}
                    >
                      <span>email</span> <Mail className="size-3" />
                    </Badge>
                  </a> */}
                  {!!invoice.pdf_url && (
                    <a href={invoice.pdf_url} target="_blank">
                      <Badge
                        variant={"outline"}
                        className="cursor-pointer flex gap-2 w-fit"
                      >
                        <span>View invoice</span> <Eye className="size-3" />
                      </Badge>
                    </a>
                  )}
                  {!invoice.pdf_url && <Button disabled>PDF</Button>}
                </div>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </Card>
  );
}
