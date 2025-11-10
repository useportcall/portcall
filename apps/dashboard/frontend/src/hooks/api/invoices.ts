import { Invoice } from "@/models/invoice";
import { useAppQuery } from "./api";
import { useMemo } from "react";

export function useListInvoices(props: { userId?: string } = {}) {
  const path = useMemo(() => {
    if (props.userId) {
      return `/invoices?user_id=${props.userId}`;
    }

    return "/invoices";
  }, [props.userId]);

  return useAppQuery<Invoice[]>({
    path,
    queryKey: [path],
  });
}

export function useListSubscriptionInvoices(subscriptionId: string) {
  const path = `/invoices?subscription_id=${subscriptionId}`;
  return useAppQuery<Invoice[]>({
    path,
    queryKey: [path],
  });
}
