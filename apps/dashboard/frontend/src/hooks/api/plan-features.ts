import { PlanFeature } from "@/models/plan-feature";
import { useParams } from "react-router-dom";
import { toast } from "sonner";
import { useAppMutation, useAppQuery } from "./api";
import { PLAN_ITEMS_PATH } from "./plan-items";

export function useListBasicPlanFeatures(planId: string) {
  return useAppQuery<PlanFeature[]>({
    path: `/plan-features?plan_id=${planId}`,
    queryKey: [`/plan-features?plan_id=${planId}`],
  });
}

export function useCreatePlanFeature() {
  const { id } = useParams();
  return useAppMutation<any, any>({
    method: "post",
    path: `/plan-features`,
    invalidate: `/plan-features?plan_id=${id}`,
    onError: () => toast("Failed to update entitlements"),
  });
}

export function useUpdatePlanFeature(planFeatureId?: string) {
  const { id } = useParams();
  return useAppMutation<any, any>({
    method: "post",
    path: `/plan-features/${planFeatureId}`,
    invalidate: [
      "/features?is_metered=true",
      `${PLAN_ITEMS_PATH}?plan_id=${id}&is_metered=true`,
    ],
    onError: () => toast("Failed to update metered feature"),
  });
}

export function useUpdatePlanFeatureForPlanItem(planFeatureId: string) {
  const { id } = useParams();
  return useAppMutation<any, any>({
    method: "post",
    path: `/plan-features/${planFeatureId}`,
    invalidate: `${PLAN_ITEMS_PATH}?plan_id=${id}&is_metered=true`,
    onError: () => toast("Failed to update metered feature"),
    onSuccess: () => {
      window.dispatchEvent(new Event("saved"));
    },
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
