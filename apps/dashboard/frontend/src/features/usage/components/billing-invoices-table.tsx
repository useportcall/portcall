import { Badge } from "@/components/ui/badge";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { BillingInvoice } from "@/hooks/api/quota";
import { Eye, FileText } from "lucide-react";
import { useTranslation } from "react-i18next";

interface BillingInvoicesTableProps {
    invoices: BillingInvoice[];
}

/**
 * BillingInvoicesTable displays the invoices from the dogfood billing system.
 */
export function BillingInvoicesTable({ invoices }: BillingInvoicesTableProps) {
    const { t, i18n } = useTranslation();
    if (!invoices || invoices.length === 0) {
        return (
            <div className="text-center py-8 text-muted-foreground">
                <FileText className="h-8 w-8 mx-auto mb-2 opacity-50" />
                <p>{t("views.usage.invoices.empty")}</p>
            </div>
        );
    }

    return (
        <Table>
            <TableHeader>
                <TableRow>
                    <TableHead>{t("views.usage.invoices.table.invoice")}</TableHead>
                    <TableHead>{t("views.usage.invoices.table.amount")}</TableHead>
                    <TableHead>{t("views.usage.invoices.table.status")}</TableHead>
                    <TableHead>{t("views.usage.invoices.table.date")}</TableHead>
                    <TableHead className="text-right">{t("views.usage.invoices.table.actions")}</TableHead>
                </TableRow>
            </TableHeader>
            <TableBody>
                {invoices.map((invoice) => (
                    <TableRow key={invoice.id}>
                        <TableCell className="font-medium">
                            {invoice.invoice_number}
                        </TableCell>
                        <TableCell>
                            {(invoice.total / 100).toLocaleString("en-US", {
                                style: "currency",
                                currency: invoice.currency,
                            })}
                        </TableCell>
                        <TableCell>
                            <Badge variant={invoice.status === "paid" ? "success" : "outline"}>
                                {invoice.status}
                            </Badge>
                        </TableCell>
                        <TableCell>
                            {new Date(invoice.created_at).toLocaleDateString(i18n.language || "en-US", {
                                day: "2-digit",
                                month: "2-digit",
                                year: "numeric",
                            })}
                        </TableCell>
                        <TableCell className="text-right">
                            {invoice.pdf_url && (
                                <a href={invoice.pdf_url} target="_blank" rel="noreferrer">
                                    <Badge
                                        variant="outline"
                                        className="cursor-pointer flex gap-1 w-fit ml-auto"
                                    >
                                        <Eye className="h-3 w-3" />
                                        {t("views.usage.invoices.table.view")}
                                    </Badge>
                                </a>
                            )}
                        </TableCell>
                    </TableRow>
                ))}
            </TableBody>
        </Table>
    );
}
