import { Card } from "@/components/ui/card";
import { TablePagination } from "@/components/table/table-pagination";
import { Table, TableBody, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { useClientPagination } from "@/hooks/use-client-pagination";
import { useListPlans } from "@/hooks";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router-dom";
import { CreatePlanButton } from "./components/create-plan-button";
import { PlanGroupComboBox } from "./components/plan-group-combo-box";
import { PlanActionsTableCell, PlanFixedChargeTableCell, PlanGroupTableCell, PlanNameTableCell, PlanStatusTableCell } from "./components/plan-table-cells";

export default function PlanListView() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [group, setGroup] = useState<string>("*");
  const { data: plans } = useListPlans(group);
  const pagination = useClientPagination(plans.data, 20);

  return (
    <div className="w-full p-4 lg:p-10 flex flex-col gap-6">
      <div className="flex flex-row justify-between">
        <div className="flex flex-col space-y-8 w-full">
          <div className="flex justify-between w-full">
            <div className="flex flex-col justify-start space-y-2">
              <h1 className="text-lg md:text-xl font-semibold">{t("views.plans.title")}</h1>
              <p className="text-sm text-muted-foreground">{t("views.plans.description")}</p>
            </div>
            <div className="flex flex-row gap-4 items-end">
              <PlanGroupComboBox value={group} onChange={setGroup} />
              <CreatePlanButton />
            </div>
          </div>
        </div>
      </div>
      <Card className="p-2 rounded-xl animate-fade-in">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-80">{t("views.plans.table.name")}</TableHead>
              <TableHead className="w-40">{t("views.plans.table.group")}</TableHead>
              <TableHead className="w-40">{t("views.plans.table.fixed_charge")}</TableHead>
              <TableHead className="w-40">{t("views.plans.table.status")}</TableHead>
              <TableHead></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {pagination.pagedItems.map((plan) => (
              <TableRow key={plan.id} onClick={() => navigate(`/plans/${plan.id}`)} className="rounded hover:bg-accent cursor-pointer">
                <PlanNameTableCell plan={plan} />
                <PlanGroupTableCell plan={plan} />
                <PlanFixedChargeTableCell plan={plan} />
                <PlanStatusTableCell plan={plan} />
                <PlanActionsTableCell plan={plan} />
              </TableRow>
            ))}
          </TableBody>
        </Table>
        {!!plans.data.length && (
          <TablePagination
            page={pagination.page}
            totalPages={pagination.totalPages}
            onPrevious={pagination.onPrevious}
            onNext={pagination.onNext}
            rowsPerPage={pagination.rowsPerPage}
            rowsPerPageOptions={[10, 20, 50]}
            onRowsPerPageChange={pagination.onRowsPerPageChange}
            previousLabel={t("views.plans.table.previous")}
            nextLabel={t("views.plans.table.next")}
            rowsPerPageLabel={t("views.plans.table.rows_per_page")}
            pageLabel={t("views.plans.table.page", {
              page: pagination.page,
              totalPages: pagination.totalPages,
              total: pagination.total,
            })}
          />
        )}
        {!plans.data.length && (
          <div className="w-full flex flex-col justify-center items-center p-10 gap-4">
            <p className="text-sm text-muted-foreground">{t("views.plans.empty")}</p>
          </div>
        )}
      </Card>
    </div>
  );
}
