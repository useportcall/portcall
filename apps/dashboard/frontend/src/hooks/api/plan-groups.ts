import { CreatePlanGroupRequest, PlanGroup } from "@/models/plan-group";
import { useAppMutation, useAppQuery } from "./api";

export function useListPlanGroups() {
  const path = "/groups";
  return useAppQuery<PlanGroup[]>({
    path,
    queryKey: [path],
  });
}

export function useCreatePlanGroup() {
  return useAppMutation<CreatePlanGroupRequest, any>({
    method: "post",
    path: "/groups",
    invalidate: ["/groups"],
  });
}
