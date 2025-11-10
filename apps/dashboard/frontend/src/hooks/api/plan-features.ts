import { PlanFeature } from "@/models/plan-feature";
import { useMemo } from "react";
import { useParams } from "react-router-dom";
import { useAppQuery, useAppMutation } from "./api";
import { toast } from "sonner";
import { PLAN_ITEMS_PATH } from "./plan-items";

export function useListPlanFeatures(props: {
  planItemId?: string;
  planId?: string;
}) {
  const path = useMemo(() => {
    let query = "?";
    if (props.planItemId) query += `plan_item_id=${props.planItemId}`;
    if (props.planId) {
      if (query.length > 1) query += "&";
      query += `plan_id=${props.planId}`;
    }
    return `/plan-features${query}`;
  }, [props.planItemId, props.planId]);
  return useAppQuery<PlanFeature[]>({ path, queryKey: [path] });
}

export function useCreatePlanFeature(plan_id: string) {
  const { id } = useParams();
  return useAppMutation<any, any>({
    method: "post",
    path: `/plan-features`,
    invalidate: [
      `/plan-features?plan_id=${plan_id}`,
      `/features?is_metered=false`,
      `${PLAN_ITEMS_PATH}?plan_id=${id}&is_metered=true`,
    ],
    onError: () => toast("Failed to update entitlements"),
  });
}

export function useUpdatePlanFeature(planFeature: PlanFeature) {
  const { id } = useParams();
  return useAppMutation<any, any>({
    method: "post",
    path: `/plan-features/${planFeature.id}`,
    invalidate: [
      `/plan-features?plan_item=${planFeature.plan_item}`,
      "/features?is_metered=true",
      `${PLAN_ITEMS_PATH}?plan_id=${id}&is_metered=true`,
    ],
    onError: () => toast("Failed to update metered feature"),
  });
}

export function useDeletePlanFeature(planFeature: PlanFeature) {
  const { id } = useParams();
  return useAppMutation<any, any>({
    method: "delete",
    path: `/plan-features/${planFeature.id}`,
    invalidate: [`/plan-features?plan_id=${id}`],
    onError: () => toast("Failed to delete feature"),
  });
}
