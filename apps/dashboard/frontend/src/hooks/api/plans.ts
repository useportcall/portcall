import { Plan } from "@/models/plan";
import { useMemo } from "react";
import { useNavigate } from "react-router-dom";
import { useAppQuery, useAppMutation } from "./api";
import { toast } from "sonner";

const PLANS_PATH = "/plans";

export function useCreatePlan() {
  const navigate = useNavigate();

  return useAppMutation<unknown, Plan>({
    method: "post",
    path: PLANS_PATH,
    invalidate: [PLANS_PATH],
    onSuccess: (result) => {
      toast("Plan created", {
        description: "The plan has been successfully created.",
      });

      navigate(`${PLANS_PATH}/${result.data.id}`);
    },
    onError: () => {
      toast("Failed to create plan", {
        description: "Please try again later.",
      });
    },
  });
}

export function useRetrievePlan(planId: string) {
  const path = `${PLANS_PATH}/${planId}`;
  return useAppQuery<Plan>({ path, queryKey: [path] });
}

export function useListPlans(group: string) {
  const query = useMemo(
    () => (group === "*" ? "" : `?group=${group}`),
    [group]
  );
  const path = PLANS_PATH + query;
  return useAppQuery<Plan[]>({ path, queryKey: [path] });
}

export function usePublishPlan(planId: string) {
  return useAppMutation<any, any>({
    method: "post",
    path: `${PLANS_PATH}/${planId}/publish`,
    invalidate: `${PLANS_PATH}/${planId}`,
    onError: () => toast("Failed to publish plan"),
  });
}

export function useDuplicatePlan(planId: string) {
  return useAppMutation<any, any>({
    method: "post",
    path: `${PLANS_PATH}/${planId}/duplicate`,
    invalidate: PLANS_PATH,
    onSuccess: () => toast("Plan duplicated"),
    onError: () => toast("Failed to duplicate plan"),
  });
}

export function useDeletePlan(planId: string) {
  return useAppMutation<any, any>({
    method: "delete",
    path: `${PLANS_PATH}/${planId}`,
    invalidate: PLANS_PATH,
    onSuccess: () => toast("Plan deleted"),
    onError: () => toast("Failed to delete plan"),
  });
}

export function useUpdatePlan(planId: string) {
  return useAppMutation<any, any>({
    method: "post",
    path: `${PLANS_PATH}/${planId}`,
    invalidate: `${PLANS_PATH}/${planId}`,
    onError: () => toast("Failed to update plan"),
  });
}
