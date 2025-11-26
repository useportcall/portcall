import { Card } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useListPlans } from "@/hooks";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { CreatePlanButton } from "./components/create-plan-button";
import { PlanGroupComboBox } from "./components/plan-group-combo-box";
import {
  PlanActionsTableCell,
  PlanFixedChargeTableCell,
  PlanGroupTableCell,
  PlanNameTableCell,
  PlanStatusTableCell,
} from "./components/plan-table-cells";

export default function PlanListView() {
  const navigate = useNavigate();
  const [group, setGroup] = useState<string>("*");
  const { data: plans } = useListPlans(group);

  return (
    <div className="w-full p-4 lg:p-10 flex flex-col gap-6">
      <div className="flex flex-row justify-between">
        <div className="flex flex-col space-y-8 w-full">
          <div className="flex justify-between w-full">
            <div className="flex flex-col justify-start space-y-2">
              <h1 className="text-lg md:text-xl font-semibold">Plans</h1>
              <p className="text-sm text-slate-400">
                Manage existing plans or add new plans here.
              </p>
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
              <TableHead className="w-80">Name</TableHead>
              <TableHead className="w-40">Group</TableHead>
              <TableHead className="w-40">Fixed charge</TableHead>
              <TableHead className="w-40">Status</TableHead>
              <TableHead></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {plans.data.map((plan) => (
              <TableRow
                key={plan.id}
                onClick={() => navigate(`/plans/${plan.id}`)}
                className="rounded hover:bg-slate-50 cursor-pointer"
              >
                <PlanNameTableCell plan={plan} />
                <PlanGroupTableCell plan={plan} />
                <PlanFixedChargeTableCell plan={plan} />
                <PlanStatusTableCell plan={plan} />
                <PlanActionsTableCell plan={plan} />
              </TableRow>
            ))}
          </TableBody>
        </Table>
        {!plans.data.length && (
          <div className="w-full flex flex-col justify-center items-center p-10 gap-4">
            <p className="text-sm text-slate-400">No plans found.</p>
          </div>
        )}
      </Card>
    </div>
  );
}
