import { BillingMeter } from "@/models/billing-meter";
import { useAppMutation, useAppQuery } from "./api";
import { toast } from "sonner";

export function useListBillingMetersBySubscription(subscriptionId: string) {
  const path = `/billing-meters?subscription_id=${subscriptionId}`;

  return useAppQuery<BillingMeter[]>({
    path,
    queryKey: [path],
  });
}

export function useListBillingMetersByUser(userId: string) {
  const path = `/billing-meters?user_id=${userId}`;

  return useAppQuery<BillingMeter[]>({
    path,
    queryKey: [path],
  });
}

export function useGetBillingMeter(subscriptionId: string, featureId: string) {
  const path = `/billing-meters/${subscriptionId}/${featureId}`;

  return useAppQuery<BillingMeter>({
    path,
    queryKey: [path],
  });
}

export interface UpdateBillingMeterRequest {
  usage?: number;
  increment?: number;
  decrement?: number;
}

export function useUpdateBillingMeter(
  subscriptionId: string,
  featureId: string,
) {
  const path = `/billing-meters/${subscriptionId}/${featureId}`;

  return useAppMutation<UpdateBillingMeterRequest, BillingMeter>({
    method: "post",
    path,
    invalidate: [`/billing-meters?subscription_id=${subscriptionId}`, path],
    onSuccess: () => {
      toast("Billing meter updated");
    },
    onError: () => {
      toast("Failed to update billing meter");
    },
  });
}

export function useResetBillingMeter(
  subscriptionId: string,
  featureId: string,
) {
  const path = `/billing-meters/${subscriptionId}/${featureId}/reset`;

  return useAppMutation<Record<string, never>, BillingMeter>({
    method: "post",
    path,
    invalidate: [
      `/billing-meters?subscription_id=${subscriptionId}`,
      `/billing-meters/${subscriptionId}/${featureId}`,
    ],
    onSuccess: () => {
      toast("Billing meter reset to zero");
    },
    onError: () => {
      toast("Failed to reset billing meter");
    },
  });
}
