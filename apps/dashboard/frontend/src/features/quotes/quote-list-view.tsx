import { Card } from "@/components/ui/card";
import { StatusFilter } from "@/components/table/status-filter";
import { TablePagination } from "@/components/table/table-pagination";
import {
  Table,
  TableBody,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useClientPagination } from "@/hooks/use-client-pagination";
import { useListQuotes } from "@/hooks/api/quotes";
import { CreateQuoteButton } from "./components/buttons";
import { useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router-dom";
import {
  QuoteNameTableCell,
  QuoteStatusTableCell,
  QuoteFixedChargeTableCell,
  QuoteActionsTableCell,
} from "./components/quote-table-cells";

export default function QuoteListView() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { data: quotes } = useListQuotes();
  const [status, setStatus] = useState("all");

  const filtered = useMemo(() => {
    if (status === "all") return quotes.data;
    return quotes.data.filter((q) => q.status === status);
  }, [quotes.data, status]);

  const pagination = useClientPagination(filtered, 20);

  return (
    <div className="w-full p-4 lg:p-10 flex flex-col gap-6">
      <div className="flex flex-row justify-between">
        <div className="flex flex-col space-y-8 w-full">
          <div className="flex justify-between w-full">
            <div className="flex flex-col justify-start space-y-2">
              <h1 className="text-lg md:text-xl font-semibold">{t("views.quotes.title")}</h1>
              <p className="text-sm text-muted-foreground">
                {t("views.quotes.description")}
              </p>
            </div>
            <div className="flex flex-row gap-4 items-end">
              <CreateQuoteButton />
            </div>
          </div>
        </div>
      </div>
      <Card className="p-2 rounded-xl animate-fade-in space-y-2">
        <div className="flex items-end gap-2">
          <StatusFilter
            label={t("views.quotes.table.status_filter_label")}
            value={status}
            options={[
              { value: "draft", label: "Draft" },
              { value: "sent", label: "Sent" },
              { value: "accepted", label: "Accepted" },
              { value: "rejected", label: "Rejected" },
              { value: "voided", label: "Voided" },
            ]}
            allLabel={t("views.quotes.table.filter_all")}
            onValueChange={(value) => {
              setStatus(value);
              pagination.resetPage();
            }}
          />
        </div>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-80">{t("views.quotes.table.name")}</TableHead>
              <TableHead className="w-40">{t("views.quotes.table.status")}</TableHead>
              <TableHead className="w-40">{t("views.quotes.table.fixed_charge")}</TableHead>
              <TableHead></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {pagination.pagedItems.map((quote) => (
              <TableRow
                key={quote.id}
                onClick={() => navigate(`/quotes/${quote.id}`)}
                className="rounded hover:bg-accent cursor-pointer"
              >
                <QuoteNameTableCell quote={quote} />
                <QuoteStatusTableCell quote={quote} />
                <QuoteFixedChargeTableCell quote={quote} />
                <QuoteActionsTableCell quote={quote} />
              </TableRow>
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
            previousLabel={t("views.quotes.table.previous")}
            nextLabel={t("views.quotes.table.next")}
            rowsPerPageLabel={t("views.quotes.table.rows_per_page")}
            pageLabel={t("views.quotes.table.page", {
              page: pagination.page,
              totalPages: pagination.totalPages,
              total: pagination.total,
            })}
          />
        )}
        {!filtered.length && (
          <div className="w-full flex flex-col justify-center items-center p-10 gap-4">
            <p className="text-sm text-muted-foreground">{t("views.quotes.empty")}</p>
          </div>
        )}
      </Card>
    </div>
  );
}
