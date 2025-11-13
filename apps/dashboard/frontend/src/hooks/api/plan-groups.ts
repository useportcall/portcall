import { CreatePlanGroupRequest, PlanGroup } from "@/models/plan-group";
import { useAppMutation, useAppQuery } from "./api";
import { useParams } from "react-router-dom";

export function useListPlanGroups() {
  const path = "/groups";
  return useAppQuery<PlanGroup[]>({
    path,
    queryKey: [path],
  });
}

export function useCreatePlanGroup() {
  const { id } = useParams();
  return useAppMutation<CreatePlanGroupRequest, any>({
    method: "post",
    path: "/groups",
    invalidate: ["/groups", "/plans/" + id],
  });
}
