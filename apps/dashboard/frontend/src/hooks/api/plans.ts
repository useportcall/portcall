import { Plan } from "@/models/plan";
import { toast } from "sonner";
import { useMemo } from "react";
import { useNavigate } from "react-router-dom";
import { useAppMutation, useAppQuery } from "./api";

const PLANS_PATH = "/plans";

function usePlanNavigate() {
  const navigate = useNavigate();
  return (id: string) => navigate(`${PLANS_PATH}/${id}`);
}

export function useCreatePlan() {
  const goToPlan = usePlanNavigate();
  return useAppMutation<unknown, Plan>({
    method: "post",
    path: PLANS_PATH,
    invalidate: [PLANS_PATH],
    onSuccess: (result) => {
      toast("Plan created", { description: "The plan has been successfully created." });
      goToPlan(result.data.id);
    },
    onError: () => toast("Failed to create plan", { description: "Please try again later." }),
  });
}

export function useRetrievePlan(planId: string) {
  const path = `${PLANS_PATH}/${planId}`;
  return useAppQuery<Plan>({ path, queryKey: [path] });
}

export function useListPlans(group: string) {
  const query = useMemo(() => (group === "*" ? "" : `?group=${group}`), [group]);
  const path = PLANS_PATH + query;
  return useAppQuery<Plan[]>({ path, queryKey: [path] });
}

export function usePublishPlan(planId: string) {
  return useAppMutation({ method: "post", path: `${PLANS_PATH}/${planId}/publish`, invalidate: `${PLANS_PATH}/${planId}`, onError: () => toast("Failed to publish plan") });
}

export function useDuplicatePlan(planId: string) {
  const goToPlan = usePlanNavigate();
  return useAppMutation({
    method: "post",
    path: `${PLANS_PATH}/${planId}/duplicate`,
    invalidate: PLANS_PATH,
    onSuccess: (result: any) => {
      toast("Plan duplicated", { description: "The plan was successfully duplicated." });
      goToPlan(result.data.id);
    },
    onError: () => toast("Failed to duplicate plan"),
  });
}

export function useCopyPlan(planId: string) {
  return useAppMutation<{ target_app_id: string }, any>({
    method: "post",
    path: `${PLANS_PATH}/${planId}/copy`,
    invalidate: PLANS_PATH,
    onSuccess: () => toast("Plan copied", { description: "The plan was successfully copied to the target app." }),
    onError: () => toast("Failed to copy plan", { description: "Please try again later." }),
  });
}

export function useDeletePlan(planId: string) {
  return useAppMutation({ method: "delete", path: `${PLANS_PATH}/${planId}`, invalidate: PLANS_PATH, onSuccess: () => toast("Plan deleted"), onError: () => toast("Failed to delete plan") });
}

export function useUpdatePlan(planId: string) {
  return useAppMutation<Partial<Plan>, Plan>({ method: "post", path: `${PLANS_PATH}/${planId}`, invalidate: `${PLANS_PATH}/${planId}`, onError: () => toast("Failed to update plan") });
}
