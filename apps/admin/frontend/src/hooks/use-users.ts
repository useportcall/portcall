import { useApiClient } from "./use-api-client";
import { useSuspenseQuery } from "@tanstack/react-query";

export interface User {
  id: string;
  name: string;
  email: string;
  company_title: string;
  payment_customer_id: string;
  created_at: string;
  has_subscription: boolean;
  has_payment_method: boolean;
  entitlement_count: number;
}

export interface UserDetail {
  id: string;
  name: string;
  email: string;
  company_title: string;
  payment_customer_id: string;
  created_at: string;
  updated_at: string;
  billing_address: {
    id: string;
    line1: string;
    line2: string;
    city: string;
    state: string;
    postal_code: string;
    country: string;
  } | null;
  subscriptions: {
    id: string;
    status: string;
    plan_name: string;
    created_at: string;
  }[];
  entitlements: {
    feature_id: string;
    usage: number;
    quota: number;
    interval: string;
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

export function useUsers(appId: string) {
  const client = useApiClient();

  return useSuspenseQuery({
    queryKey: ["apps", appId, "users"],
    queryFn: async () => {
      const result = await client.get<{ data: User[] }>(`/apps/${appId}/users`);
      return result.data.data;
    },
  });
}

export function useUser(appId: string, userId: string) {
  const client = useApiClient();

  return useSuspenseQuery({
    queryKey: ["apps", appId, "users", userId],
    queryFn: async () => {
      const result = await client.get<{ data: UserDetail }>(
        `/apps/${appId}/users/${userId}`,
      );
      return result.data.data;
    },
  });
}
