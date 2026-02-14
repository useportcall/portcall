import EmptyTable from "@/components/empty-table";
import { TablePagination } from "@/components/table/table-pagination";
import { Card } from "@/components/ui/card";
import { Table, TableBody, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { useClientPagination } from "@/hooks/use-client-pagination";
import { useListInvoices } from "@/hooks/api/invoices";
import { useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { InvoiceListFilters } from "./invoice-list-filters";
import { matchesInvoice } from "./invoice-list-match";
import { InvoiceListRow } from "./invoice-list-row";

export function InvoiceListTable() {
  const { t, i18n } = useTranslation();
  const [status, setStatus] = useState("all");
  const { data: invoices } = useListInvoices();
  const filtered = useMemo(
    () => invoices.data.filter((invoice) => matchesInvoice(invoice, status)),
    [invoices.data, status],
  );
  const pagination = useClientPagination(filtered, 20);

  if (!invoices.data.length) return <EmptyTable message={t("views.invoices.table.empty")} />;

  return (
    <Card data-testid="invoice-table" className="p-2 animate-fade-in space-y-2">
      <InvoiceListFilters
        status={status}
        onStatusChange={(value) => (setStatus(value), pagination.resetPage())}
        t={t}
      />
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>{t("views.invoices.table.invoice_number")}</TableHead>
            <TableHead>{t("views.invoices.table.total")}</TableHead>
            <TableHead>{t("views.invoices.table.customer")}</TableHead>
            <TableHead>{t("views.invoices.table.status")}</TableHead>
            <TableHead>{t("views.invoices.table.issued")}</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {pagination.pagedItems.map((invoice) => (
            <InvoiceListRow
              key={invoice.id}
              invoice={invoice}
              locale={i18n.language}
              viewInvoiceLabel={t("views.invoices.table.view_invoice")}
              pdfLabel={t("views.invoices.table.pdf")}
            />
          ))}
        </TableBody>
      </Table>
      {!!filtered.length && (
        <TablePagination
          page={pagination.page}
          totalPages={pagination.totalPages}
          onPrevious={pagination.onPrevious}
          onNext={pagination.onNext}
          rowsPerPage={pagination.rowsPerPage}
          rowsPerPageOptions={[10, 20, 50]}
          onRowsPerPageChange={pagination.onRowsPerPageChange}
          previousLabel={t("views.invoices.table.previous")}
          nextLabel={t("views.invoices.table.next")}
          rowsPerPageLabel={t("views.invoices.table.rows_per_page")}
          pageLabel={t("views.invoices.table.page", {
            page: pagination.page,
            totalPages: pagination.totalPages,
            total: pagination.total,
          })}
        />
      )}
      {!filtered.length && <EmptyTable message={t("views.invoices.table.no_results")} />}
    </Card>
  );
}
