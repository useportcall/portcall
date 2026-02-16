import { useApiClient } from "./use-api-client";
import { useSuspenseQuery } from "@tanstack/react-query";

export interface Subscription {
  id: string;
  user_id: string;
  user_name: string;
  user_email: string;
  plan_id: string;
  plan_name: string;
  status: string;
  currency: string;
  billing_interval: string;
  last_reset_at: string;
  next_reset_at: string;
  created_at: string;
  item_count: number;
}

export interface SubscriptionDetail {
  id: string;
  user_id: string;
  user_name: string;
  user_email: string;
  plan_id: string;
  plan_name: string;
  status: string;
  currency: string;
  billing_interval: string;
  billing_interval_count: number;
  invoice_due_by_days: number;
  last_reset_at: string;
  next_reset_at: string;
  final_reset_at: string | null;
  created_at: string;
  updated_at: string;
  discount_pct: number;
  discount_qty: number;
  items: {
    id: string;
    title: string;
    description: string;
    pricing_model: string;
    unit_amount: number;
    quantity: number;
    usage: number;
  }[];
  invoices: {
    id: string;
    invoice_number: string;
    status: string;
    total: number;
    currency: string;
    due_by: string;
  }[];
}

export function useSubscriptions(appId: string) {
  const client = useApiClient();

  return useSuspenseQuery({
    queryKey: ["apps", appId, "subscriptions"],
    queryFn: async () => {
      const result = await client.get<{ data: Subscription[] }>(
        `/apps/${appId}/subscriptions`,
      );
      return result.data.data;
    },
  });
}

export function useSubscription(appId: string, subscriptionId: string) {
  const client = useApiClient();

  return useSuspenseQuery({
    queryKey: ["apps", appId, "subscriptions", subscriptionId],
    queryFn: async () => {
      const result = await client.get<{ data: SubscriptionDetail }>(
        `/apps/${appId}/subscriptions/${subscriptionId}`,
      );
      return result.data.data;
    },
  });
}
