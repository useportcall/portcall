"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { useAuth } from "./use-auth";
import axios from "axios";

export function useUpdateSubscription(subscriptionId: string) {
  return useMutation({
    mutationFn: async ({ planId }: any) => {
      try {
        const url = `/api/subscriptions/${subscriptionId}`;
        const result = await axios.post(url, { plan_id: planId });
        return result.data;
      } catch (error) {
        toast("Error updating subscription");
      }
    },
  });
}

export function useGetFeature(featureId: string) {
  const { user } = useAuth();

  return useQuery({
    queryKey: [`/api/users/${user!.id}/entitlements/${featureId}`],
    queryFn: async () => {
      const url = `/api/users/${user!.id}/entitlements/${featureId}`;
      const result = await axios.get(url);
      return result.data;
    },
    enabled: !!user,
    refetchInterval: 5000,
  });
}

export function useCreateMeterEvent(featureId: string) {
  const { user } = useAuth();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      const result = await axios.post("/api/meter-events", {
        usage: 20,
        user_id: user!.id,
        feature_id: featureId,
      });

      return result.data;
    },
    onSuccess: (_) => {
      queryClient.invalidateQueries({
        queryKey: [`/api/users/${user!.id}/entitlements/${featureId}`],
      });
    },
    onError: (err) => {
      console.error("Feature usage increment failed:", err);
      toast("Error incrementing feature usage");
    },
  });
}
