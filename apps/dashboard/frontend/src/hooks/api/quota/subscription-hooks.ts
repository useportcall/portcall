import { useAppMutation, useAppQuery } from "../api";
import { redirectToCheckout, redirectWithQuery } from "./helpers";
import {
  DowngradeToFreeRequest,
  DowngradeToFreeResponse,
  UpgradeToProResponse,
  UserBillingSubscription,
} from "./types";

export function useUserBillingSubscription() {
  return useAppQuery<UserBillingSubscription>({
    path: "/billing/subscription",
    queryKey: ["billing", "subscription"],
  });
}

export function useUpgradeToPro() {
  return useAppMutation<{ cancel_url: string; redirect_url: string }, UpgradeToProResponse>({
    method: "post",
    path: "/billing/upgrade-to-pro",
    invalidate: ["billing", "subscription"],
    onSuccess: (response) => {
      if (response.data?.success) {
        redirectWithQuery("upgrade=success");
        return;
      }
      if (response.data?.checkout_url) {
        redirectToCheckout(response.data.checkout_url);
      }
    },
  });
}

export function useDowngradeToFree() {
  return useAppMutation<DowngradeToFreeRequest, DowngradeToFreeResponse>({
    method: "post",
    path: "/billing/downgrade-to-free",
    invalidate: ["/billing/subscription", "/billing/quota/subscriptions", "/billing/quota/users"],
    onSuccess: (response) => redirectWithQuery(`downgrade=success${response.data?.scheduled ? "&scheduled=true" : ""}`),
  });
}
