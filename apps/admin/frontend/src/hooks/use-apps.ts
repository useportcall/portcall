import { useApiClient } from "./use-api-client";
import { useSuspenseQuery } from "@tanstack/react-query";

export interface App {
  id: number;
  public_id: string;
  name: string;
  is_live: boolean;
  created_at: string;
  user_count: number;
  subscription_count: number;
  plan_count: number;
  quote_count: number;
  invoice_count: number;
}

export interface AppDetail extends App {
  account_id: number;
  account_email: string;
  feature_count: number;
  connection_count: number;
  updated_at: string;
}

export function useApps() {
  const client = useApiClient();

  return useSuspenseQuery({
    queryKey: ["apps"],
    queryFn: async () => {
      const result = await client.get<{ data: App[] }>("/apps");
      return result.data.data;
    },
  });
}

export function useApp(appId: string) {
  const client = useApiClient();

  return useSuspenseQuery({
    queryKey: ["apps", appId],
    queryFn: async () => {
      const result = await client.get<{ data: AppDetail }>(`/apps/${appId}`);
      return result.data.data;
    },
  });
}
