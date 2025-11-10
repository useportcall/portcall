import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Loader2, MoreHorizontal } from "lucide-react";
import { Link, useNavigate } from "react-router-dom";

import { TableCell, TableRow } from "@/components/ui/table";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { cn } from "@/lib/utils";
import { Plan } from "@/models/plan";
import { useEffect, useMemo, useRef, useState } from "react";
import { useDeletePlan, useDuplicatePlan } from "@/hooks";

function formatFixedCharge(plan: Plan) {
  const fixedPlanItem = plan.items.find(
    (item) => item.pricing_model === "fixed"
  );

  if (!fixedPlanItem) return "";

  const currency = plan.currency.toUpperCase();

  const fee = Number(fixedPlanItem.unit_amount / 100).toFixed(2);

  let interval = "";
  if (plan.interval === "month") interval = "mo";
  else if (plan.interval === "quarter") interval = "qtr";
  else if (plan.interval === "year") interval = "yr";
  else if (plan.interval === "week") interval = "wk";
  return `${currency} ${fee} / ${interval}`;
}

export function PlanRow({ plan }: { plan: Plan }) {
  const navigate = useNavigate();

  return (
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
  );
}

export function PlanGroupTableCell(props: { plan: Plan }) {
  return (
    <TableCell className="w-40">
      {!!props.plan.plan_group && (
        <Badge variant={"outline"}>{props.plan.plan_group.name}</Badge>
      )}
    </TableCell>
  );
}

export function PlanStatusTableCell(props: { plan: Plan }) {
  return (
    <TableCell className="w-40">
      <Badge variant={"outline"}>{props.plan.status}</Badge>
    </TableCell>
  );
}

export function PlanFixedChargeTableCell(props: { plan: Plan }) {
  const fixedCharge = useMemo(
    () => formatFixedCharge(props.plan),
    [props.plan]
  );

  return <TableCell className="w-40">{fixedCharge}</TableCell>;
}

export function PlanNameTableCell(props: { plan: Plan }) {
  const [open, setOpen] = useState(false);
  const [copied, setCopied] = useState(false);
  const copyRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const ref = useRef<boolean>(false);

  useEffect(() => {
    return () => {
      if (copyRef.current) {
        clearTimeout(copyRef.current);
      }
    };
  }, []);

  function handleCopied() {
    // Clear any existing timer to prevent multiple timers running
    if (copyRef.current) {
      clearTimeout(copyRef.current);
    }

    setCopied(true);

    // Store the timer ID in a ref
    copyRef.current = setTimeout(() => {
      ref.current = false;
      setOpen(false);
      setCopied(false);
    }, 1000);
  }

  return (
    <TableCell className="w-80">
      <div>
        <h4 className="font-semibold overflow-ellipsis">{props.plan.name}</h4>
        <Tooltip
          delayDuration={100}
          open={open || copied}
          onOpenChange={(o) => {
            if (ref.current) return;
            setOpen(o);
          }}
        >
          <TooltipTrigger asChild>
            <span
              className={cn(
                "text-xs text-slate-500 italic transition-colors hover:bg-slate-200"
              )}
              onClick={(e) => {
                // don't propagate
                e.stopPropagation();

                // set ref
                ref.current = true;

                // copy to clipboard
                navigator.clipboard.writeText(props.plan.id);

                handleCopied();
              }}
            >
              {props.plan.id}
            </span>
          </TooltipTrigger>
          <TooltipContent className="p-1">
            {copied ? "copied!" : "copy"}
          </TooltipContent>
        </Tooltip>
      </div>
    </TableCell>
  );
}

export function PlanActionsTableCell({ plan }: { plan: Plan }) {
  return (
    <TableCell className="w-40" align="right">
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button size="icon" variant="ghost" className="h-8 w-8 p-0">
            <MoreHorizontal className="w-5 h-5" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <EditPlanMenuItem planId={plan.id} />
          <DuplicatePlanMenuItem plan={plan} />
          <DeletePlanMenuItem plan={plan} />
        </DropdownMenuContent>
      </DropdownMenu>
    </TableCell>
  );
}

function EditPlanMenuItem({ planId }: { planId: string }) {
  return (
    <DropdownMenuItem asChild>
      <Link to={`/plans/${planId}`}>Edit</Link>
    </DropdownMenuItem>
  );
}

function DuplicatePlanMenuItem({ plan }: { plan: Plan }) {
  const { mutate, isPending } = useDuplicatePlan(plan.id);

  return (
    <DropdownMenuItem disabled={isPending} onClick={() => mutate({})}>
      {isPending ? (
        <span className="flex items-center gap-2">
          <Loader2 className="animate-spin w-4 h-4" /> Duplicating...
        </span>
      ) : (
        "Duplicate"
      )}
    </DropdownMenuItem>
  );
}

function DeletePlanMenuItem({ plan }: { plan: Plan }) {
  const { mutate, isPending } = useDeletePlan(plan.id);

  return (
    <DropdownMenuItem
      disabled={isPending || plan.status !== "draft"}
      onClick={(e) => {
        e.stopPropagation();
        mutate({});
      }}
      className="text-red-600 focus:text-red-700"
    >
      {isPending ? (
        <span className="flex items-center gap-2">
          <Loader2 className="animate-spin w-4 h-4" /> Deleting...
        </span>
      ) : (
        "Delete"
      )}
    </DropdownMenuItem>
  );
}
