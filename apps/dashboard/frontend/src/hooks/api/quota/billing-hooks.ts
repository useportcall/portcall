import { useAppMutation, useAppQuery } from "../api";
import {
  BillingAddress,
  BillingAddressResponse,
  BillingInvoicesResponse,
  UpsertBillingAddressRequest,
} from "./types";

export function useBillingInvoices() {
  return useAppQuery<BillingInvoicesResponse>({
    path: "/billing/invoices",
    queryKey: ["billing", "invoices"],
  });
}

export function useBillingAddress() {
  return useAppQuery<BillingAddressResponse>({
    path: "/billing/address",
    queryKey: ["billing", "address"],
  });
}

export function useUpsertBillingAddress() {
  return useAppMutation<UpsertBillingAddressRequest, BillingAddress>({
    method: "post",
    path: "/billing/address",
    invalidate: ["/billing/address"],
  });
}
