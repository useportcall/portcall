import EmptyTable from "@/components/empty-table";
import { TablePagination } from "@/components/table/table-pagination";
import { Card } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useClientPagination } from "@/hooks/use-client-pagination";
import { useListSubscriptions } from "@/hooks";
import { useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router-dom";
import { SubscriptionListRow } from "./subscription-list-row";
import { SubscriptionListToolbar } from "./subscription-list-toolbar";

export function SubscriptionListTable() {
  const { t, i18n } = useTranslation();
  const navigate = useNavigate();
  const { data: subscriptions } = useListSubscriptions();
  const [search, setSearch] = useState("");
  const [status, setStatus] = useState("all");

  const filtered = useMemo(() => {
    const query = search.trim().toLowerCase();
    return subscriptions.data.filter((subscription) => {
      const statusMatch = status === "all" || subscription.status.toLowerCase() === status;
      if (!statusMatch) return false;
      if (!query) return true;
      const email = subscription.user?.email?.toLowerCase() || "";
      const userId = subscription.user?.id?.toLowerCase() || "";
      const planName = subscription.plan?.name?.toLowerCase() || "";
      return [email, userId, planName, subscription.status.toLowerCase()].some((value) => value.includes(query));
    });
  }, [subscriptions.data, search, status]);
  const pagination = useClientPagination(filtered, 20);

  if (!subscriptions.data.length) {
    return <EmptyTable message={t("views.subscriptions.table.empty")} button={<></>} />;
  }

  return (
    <>
      <SubscriptionListToolbar
        search={search}
        status={status}
        onSearchChange={(value) => {
          setSearch(value);
          pagination.resetPage();
        }}
        onStatusChange={(value) => {
          setStatus(value);
          pagination.resetPage();
        }}
        searchPlaceholder={t("views.subscriptions.table.search")}
        statusLabel={t("views.subscriptions.table.status_filter_label")}
        allStatusesLabel={t("views.subscriptions.table.filter_all")}
      />
      <Card className="p-2 animate-fade-in">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>{t("views.subscriptions.table.user")}</TableHead>
              <TableHead>{t("views.subscriptions.table.plan")}</TableHead>
              <TableHead>{t("views.subscriptions.table.status")}</TableHead>
              <TableHead>{t("views.subscriptions.table.created")}</TableHead>
              <TableHead>{t("views.subscriptions.table.next_reset")}</TableHead>
              <TableHead>ID</TableHead>
              <TableHead />
            </TableRow>
          </TableHeader>
          <TableBody>
            {pagination.pagedItems.map((subscription) => (
              <SubscriptionListRow
                key={subscription.id}
                subscription={subscription}
                locale={i18n.language}
                onNavigate={() => navigate(`/subscriptions/${subscription.id}`)}
                actions={[
                  {
                    label: t("views.subscriptions.table.cancel"),
                    disabled: subscription.status !== "active",
                  },
                ]}
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
            previousLabel={t("views.subscriptions.table.previous")}
            nextLabel={t("views.subscriptions.table.next")}
            rowsPerPageLabel={t("views.subscriptions.table.rows_per_page")}
            pageLabel={t("views.subscriptions.table.page", {
              page: pagination.page,
              totalPages: pagination.totalPages,
              total: pagination.total,
            })}
          />
        )}
        {!filtered.length && <EmptyTable message={t("views.subscriptions.table.no_results")} button={<></>} />}
      </Card>
    </>
  );
}
