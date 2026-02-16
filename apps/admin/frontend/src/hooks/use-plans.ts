import { useApiClient } from "./use-api-client";
import { useSuspenseQuery } from "@tanstack/react-query";

export interface Plan {
  id: string;
  name: string;
  status: string;
  currency: string;
  interval: string;
  interval_count: number;
  is_free: boolean;
  trial_period_days: number;
  discount_pct: number;
  created_at: string;
  subscription_count: number;
  item_count: number;
}

export interface PlanDetail {
  id: string;
  name: string;
  status: string;
  currency: string;
  interval: string;
  interval_count: number;
  is_free: boolean;
  trial_period_days: number;
  invoice_due_by_days: number;
  discount_pct: number;
  discount_qty: number;
  created_at: string;
  updated_at: string;
  group_name: string;
  items: {
    id: string;
    pricing_model: string;
    quantity: number;
    unit_amount: number;
    title: string;
    description: string;
    unit_label: string;
  }[];
  features: {
    id: string;
    feature_id: string;
    interval: string;
    quota: number;
    rollover: number;
  }[];
  subscription_count: number;
}

export function usePlans(appId: string) {
  const client = useApiClient();

  return useSuspenseQuery({
    queryKey: ["apps", appId, "plans"],
    queryFn: async () => {
      const result = await client.get<{ data: Plan[] }>(`/apps/${appId}/plans`);
      return result.data.data;
    },
  });
}

export function usePlan(appId: string, planId: string) {
  const client = useApiClient();

  return useSuspenseQuery({
    queryKey: ["apps", appId, "plans", planId],
    queryFn: async () => {
      const result = await client.get<{ data: PlanDetail }>(
        `/apps/${appId}/plans/${planId}`,
      );
      return result.data.data;
    },
  });
}
