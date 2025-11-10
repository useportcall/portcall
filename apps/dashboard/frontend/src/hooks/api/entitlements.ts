import { Entitlement } from "@/models/entitlement";
import { useAppQuery } from "./api";

export function useListMeteredEntitlements(props: { userId: string }) {
  return useAppQuery<Entitlement[]>({
    path: "/entitlements" + "?user_id=" + props.userId + "&is_metered=true",
    queryKey: [
      "/entitlements" + "?user_id=" + props.userId + "&is_metered=true",
    ],
  });
}

export function useListBasicEntitlements(props: { userId: string }) {
  return useAppQuery<Entitlement[]>({
    path: "/entitlements" + "?user_id=" + props.userId + "&is_metered=false",
    queryKey: [
      "/entitlements" + "?user_id=" + props.userId + "&is_metered=false",
    ],
  });
}
