import { CopyableIdCell } from "@/components/table/copyable-id-cell";
import { ActionsCell, ActionMenuItem } from "@/components/table/actions-cell";
import { StatusBadgeCell } from "@/components/table/status-badge-cell";
import { Badge } from "@/components/ui/badge";
import { TableCell, TableRow } from "@/components/ui/table";
import { useDeletePlan, useDuplicatePlan, useListApps } from "@/hooks";
import { useApp } from "@/hooks/use-app";
import { formatFixedCharge } from "@/lib/format";
import { App } from "@/models/app";
import { Plan } from "@/models/plan";
import { createPortal } from "react-dom";
import { Link, useNavigate } from "react-router-dom";
import { useMemo, useState } from "react";
import { CopyPlanDialog } from "./copy-plan-dialog";

export function PlanRow({ plan }: { plan: Plan }) {
  const navigate = useNavigate();
  return <TableRow onClick={() => navigate(`/plans/${plan.id}`)} className="rounded hover:bg-accent cursor-pointer"><PlanNameTableCell plan={plan} /><PlanGroupTableCell plan={plan} /><PlanFixedChargeTableCell plan={plan} /><StatusBadgeCell status={plan.status} /><PlanActionsTableCell plan={plan} /></TableRow>;
}

export function PlanGroupTableCell({ plan }: { plan: Plan }) { return <TableCell className="w-40">{!!plan.plan_group && <Badge variant="outline">{plan.plan_group.name}</Badge>}</TableCell>; }
export function PlanFixedChargeTableCell({ plan }: { plan: Plan }) { return <TableCell className="w-40">{useMemo(() => formatFixedCharge(plan), [plan])}</TableCell>; }
export function PlanNameTableCell({ plan }: { plan: Plan }) { return <CopyableIdCell title={plan.name} id={plan.id} />; }
export function PlanStatusTableCell({ plan }: { plan: Plan }) { return <StatusBadgeCell status={plan.status} />; }

export function PlanActionsTableCell({ plan }: { plan: Plan }) {
  const [copyDialogOpen, setCopyDialogOpen] = useState(false);
  const { data: apps } = useListApps();
  const { id } = useApp();
  const { mutateAsync: duplicatePlan, isPending: isDuplicating } = useDuplicatePlan(plan.id);
  const { mutate: deletePlan, isPending: isDeleting } = useDeletePlan(plan.id);
  const copyLabel = useMemo(() => { const app = apps.data.find((a: App) => a.id === id); if (!app) return "Copy to app"; return app.is_live ? "Copy to test" : "Copy to prod"; }, [id, apps.data.length]);

  const actions: ActionMenuItem[] = [
    { label: <Link to={`/plans/${plan.id}`}>Edit</Link>, asChild: true },
    { label: "Duplicate", onClick: (e) => { e.stopPropagation(); duplicatePlan({}); }, loading: isDuplicating, disabled: isDuplicating },
    { label: copyLabel, onClick: (e) => { e.stopPropagation(); setCopyDialogOpen(true); } },
    { label: "Delete", onClick: (e) => { e.stopPropagation(); deletePlan({}); }, loading: isDeleting, disabled: isDeleting || plan.status !== "draft", className: "text-red-600 focus:text-red-700" },
  ];

  return <>{<ActionsCell actions={actions} />}{copyDialogOpen && createPortal(<CopyPlanDialog plan={plan} open={copyDialogOpen} onOpenChange={setCopyDialogOpen} />, document.body)}</>;
}
