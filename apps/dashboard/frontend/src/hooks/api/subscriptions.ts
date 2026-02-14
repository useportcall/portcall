import { Subscription } from "@/models/subscription";
import { useAppMutation, useAppQuery } from "./api";
import { toast } from "sonner";

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

export interface CreateSubscriptionRequest {
  user_id: string;
  plan_id: string;
}

export function useCreateSubscription() {
  return useAppMutation<CreateSubscriptionRequest, Subscription>({
    method: "post",
    path: `/subscriptions`,
    invalidate: ["/subscriptions"],
    onSuccess: () => {
      toast("Subscription created", {
        description: "The user has been subscribed to the plan.",
      });
    },
    onError: () => {
      toast("Failed to create subscription", {
        description: "Please try again later.",
      });
    },
  });
}

export function useCancelSubscription(subscriptionId: string, userId: string) {
  return useAppMutation<Record<string, never>, Subscription>({
    method: "post",
    path: `/subscriptions/${subscriptionId}/cancel`,
    invalidate: [
      // `/subscriptions/${subscriptionId}`,
      `/users/${userId}/subscription`,
    ],
  });
}
