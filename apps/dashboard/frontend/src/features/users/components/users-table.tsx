import EmptyTable from "@/components/empty-table";
import { useListUsersPaginated } from "@/hooks";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { CreateUser } from "./create-user";
import { UsersFilters, type UserFilterOption } from "./users-filters";
import { UsersTableCard } from "./users-table-card";

export function UsersTable() {
  const { t } = useTranslation();
  const [search, setSearch] = useState("");
  const [subscribed, setSubscribed] = useState<UserFilterOption>("all");
  const [paymentMethod, setPaymentMethod] = useState<UserFilterOption>("all");
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(20);

  const { data, isFetching } = useListUsersPaginated({
    search,
    subscribed,
    paymentMethodAdded: paymentMethod,
    page,
    limit,
  });
  const users = data?.data.users ?? [];
  const currentPage = data?.data.page ?? page;
  const totalPages = data?.data.total_pages ?? 1;
  const total = data?.data.total ?? 0;

  return (
    <>
      <div className="flex flex-wrap items-start gap-2 justify-between">
        <UsersFilters
          search={search}
          subscription={subscribed}
          paymentMethod={paymentMethod}
          onSearchChange={(value) => {
            setSearch(value);
            setPage(1);
          }}
          onSubscriptionChange={(value) => {
            setSubscribed(value);
            setPage(1);
          }}
          onPaymentMethodChange={(value) => {
            setPaymentMethod(value);
            setPage(1);
          }}
          onReset={() => {
            setSearch("");
            setSubscribed("all");
            setPaymentMethod("all");
            setLimit(20);
            setPage(1);
          }}
          t={t}
        />
        <CreateUser />
      </div>

      {!!users.length && (
        <UsersTableCard
          users={users}
          page={currentPage}
          totalPages={totalPages}
          total={total}
          limit={limit}
          onPrevious={() => setPage((p) => Math.max(1, p - 1))}
          onNext={() => setPage((p) => p + 1)}
          onLimitChange={(value) => {
            setLimit(value);
            setPage(1);
          }}
          t={t}
        />
      )}
      {isFetching && <p className="text-xs text-muted-foreground px-1">{t("common.loading")}</p>}
      {!users.length && !isFetching && (
        <EmptyTable message={t("views.users.table.empty")} button={""} />
      )}
    </>
  );
}
