import { Subscription } from "@/models/subscription";
import { useAppMutation, useAppQuery } from "./api";

export function useListSubscriptions() {
  return useAppQuery<Subscription[]>({
    path: `/subscriptions`,
    queryKey: ["/subscriptions"],
  });
}

export function useGetSubscription(subscriptionId: string) {
  return useAppQuery<Subscription>({
    path: `/subscriptions/${subscriptionId}`,
    queryKey: [`/subscriptions/${subscriptionId}`],
  });
}

export function useGetUserSubscription(userId: string) {
  return useAppQuery<Subscription>({
    path: `/users/${userId}/subscription`,
    queryKey: [`/users/${userId}/subscription`],
  });
}

export function useCancelSubscription(subscriptionId: string) {
  return useAppMutation<void, Subscription>({
    method: "post",
    path: `/subscriptions/${subscriptionId}/cancel`,
    invalidate: [`/subscriptions/${subscriptionId}`],
  });
}
