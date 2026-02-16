import { useApiClient } from "./use-api-client";
import { useSuspenseQuery } from "@tanstack/react-query";

export interface Stats {
  total_apps: number;
  total_users: number;
  total_subscriptions: number;
  total_plans: number;
  total_quotes: number;
  total_invoices: number;
  total_accounts: number;
  active_subscriptions: number;
  draft_quotes: number;
  paid_invoices: number;
  pending_invoices: number;
}

export function useStats() {
  const client = useApiClient();

  return useSuspenseQuery({
    queryKey: ["stats"],
    queryFn: async () => {
      const result = await client.get<{ data: Stats }>("/stats");
      return result.data.data;
    },
  });
}
