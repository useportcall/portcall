import { Entitlement } from "@/models/entitlement";
import { toast } from "sonner";
import { useAppMutation, useAppQuery } from "./api";

function entitlementsPath(userId: string, isMetered: boolean) {
  return `/entitlements?user_id=${userId}&is_metered=${isMetered}`;
}

function entitlementInvalidate(userId: string) {
  return [entitlementsPath(userId, true), entitlementsPath(userId, false)];
}

export function useListMeteredEntitlements(props: { userId: string }) {
  const path = entitlementsPath(props.userId, true);
  return useAppQuery<Entitlement[]>({ path, queryKey: [path] });
}

export function useListBasicEntitlements(props: { userId: string }) {
  const path = entitlementsPath(props.userId, false);
  return useAppQuery<Entitlement[]>({ path, queryKey: [path] });
}

export function useCreateEntitlement(props: { userId: string }) {
  return useAppMutation<{ feature_id: string; quota?: number; interval?: string }, Entitlement>({
    method: "post",
    path: `/entitlements/${props.userId}`,
    invalidate: entitlementInvalidate(props.userId),
    onSuccess: () => toast("Entitlement created"),
    onError: () => toast("Failed to create entitlement"),
  });
}

export function useUpdateEntitlement(props: { userId: string; featureId: string }) {
  return useAppMutation<{ quota?: number; usage?: number }, Entitlement>({
    method: "post",
    path: `/entitlements/${props.userId}/${props.featureId}`,
    invalidate: entitlementInvalidate(props.userId),
    onSuccess: () => toast("Entitlement updated"),
    onError: () => toast("Failed to update entitlement"),
  });
}

export function useDeleteEntitlement(props: { userId: string; featureId: string }) {
  return useAppMutation<Record<string, never>, { deleted: boolean }>({
    method: "delete",
    path: `/entitlements/${props.userId}/${props.featureId}`,
    invalidate: entitlementInvalidate(props.userId),
    onSuccess: () => toast("Entitlement removed"),
    onError: () => toast("Failed to remove entitlement"),
  });
}

export function useToggleUserFeature(props: { userId: string }) {
  return useAppMutation<{ feature_id: string; enabled: boolean }, Entitlement | { deleted: boolean }>({
    method: "post",
    path: `/entitlements/${props.userId}/toggle`,
    invalidate: entitlementInvalidate(props.userId),
    onSuccess: () => toast("Feature updated"),
    onError: () => toast("Failed to update feature"),
  });
}
