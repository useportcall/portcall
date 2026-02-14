import { PlanFeature } from "@/models/plan-feature";
import { toast } from "sonner";
import { useAppMutation, useAppQuery } from "./api";
import { PLAN_ITEMS_PATH } from "./plan-items";

export function useListBasicPlanFeatures(planId: string) {
  return useAppQuery<PlanFeature[]>({
    path: `/plan-features?plan_id=${planId}`,
    queryKey: [`/plan-features?plan_id=${planId}`],
  });
}

export function useCreatePlanFeature(planId: string) {
  return useAppMutation<any, any>({
    method: "post",
    path: `/plan-features`,
    invalidate: `/plan-features?plan_id=${planId}`,
    onError: () => toast("Failed to update entitlements"),
  });
}

export function useUpdatePlanFeature(planFeatureId: string, planId: string) {
  return useAppMutation<any, any>({
    method: "post",
    path: `/plan-features/${planFeatureId}`,
    invalidate: [
      "/features?is_metered=true",
      `${PLAN_ITEMS_PATH}?plan_id=${planId}&is_metered=true`,
    ],
    onError: () => toast("Failed to update metered feature"),
  });
}

export function useUpdatePlanFeatureForPlanItem(
  planFeatureId: string,
  planId: string
) {
  return useAppMutation<any, any>({
    method: "post",
    path: `/plan-features/${planFeatureId}`,
    invalidate: `${PLAN_ITEMS_PATH}?plan_id=${planId}&is_metered=true`,
    onError: () => toast("Failed to update metered feature"),
    onSuccess: () => {
      window.dispatchEvent(new Event("saved"));
    },
  });
}

export function useDeletePlanFeature(planFeature: PlanFeature, planId: string) {
  return useAppMutation<any, any>({
    method: "delete",
    path: `/plan-features/${planFeature.id}`,
    invalidate: [`/plan-features?plan_id=${planId}`],
    onError: () => toast("Failed to delete feature"),
  });
}
