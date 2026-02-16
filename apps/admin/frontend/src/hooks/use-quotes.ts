import { useApiClient } from "./use-api-client";
import { useSuspenseQuery } from "@tanstack/react-query";

export interface Quote {
  id: string;
  public_title: string;
  public_name: string;
  company_name: string;
  status: string;
  days_valid: number;
  user_id: string;
  user_name: string;
  user_email: string;
  plan_name: string;
  recipient_email: string;
  issued_at: string | null;
  expires_at: string | null;
  accepted_at: string | null;
  created_at: string;
}

export interface QuoteDetail {
  id: string;
  public_title: string;
  public_name: string;
  company_name: string;
  status: string;
  days_valid: number;
  user_id: string;
  user_name: string;
  user_email: string;
  plan_id: string;
  plan_name: string;
  recipient_email: string;
  direct_checkout: boolean;
  url: string | null;
  toc: string;
  issued_at: string | null;
  expires_at: string | null;
  accepted_at: string | null;
  voided_at: string | null;
  created_at: string;
  updated_at: string;
}

export function useQuotes(appId: string) {
  const client = useApiClient();

  return useSuspenseQuery({
    queryKey: ["apps", appId, "quotes"],
    queryFn: async () => {
      const result = await client.get<{ data: Quote[] }>(
        `/apps/${appId}/quotes`,
      );
      return result.data.data;
    },
  });
}

export function useQuote(appId: string, quoteId: string) {
  const client = useApiClient();

  return useSuspenseQuery({
    queryKey: ["apps", appId, "quotes", quoteId],
    queryFn: async () => {
      const result = await client.get<{ data: QuoteDetail }>(
        `/apps/${appId}/quotes/${quoteId}`,
      );
      return result.data.data;
    },
  });
}
