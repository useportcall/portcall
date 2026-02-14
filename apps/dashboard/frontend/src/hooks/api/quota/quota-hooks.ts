import { useAppQuery } from "../api";
import { QuotaStatus } from "./types";

export function useUserBillingSubscriptionQuota() {
  return useAppQuery<QuotaStatus>({
    path: "/billing/quota/subscriptions",
    queryKey: ["billing", "quota", "subscriptions"],
  });
}

export function useUserBillingUserQuota() {
  return useAppQuery<QuotaStatus>({
    path: "/billing/quota/users",
    queryKey: ["billing", "quota", "users"],
  });
}
