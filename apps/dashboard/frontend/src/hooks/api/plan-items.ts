import { MeteredPlanItem, PlanItem } from "@/models/plan-item";
import { useMemo } from "react";
import { useParams } from "react-router-dom";
import { useAppQuery, useAppMutation } from "./api";
import { toast } from "sonner";

export const PLAN_ITEMS_PATH = "/plan-items";

export function useListMeteredPlanItems() {
  const { id } = useParams();
  const path = `${PLAN_ITEMS_PATH}?plan_id=${id}&is_metered=true`;
  const query = useAppQuery<PlanItem[]>({ path, queryKey: [path] });

  const data: MeteredPlanItem[] = useMemo(() => {
    if (!query.data || !Array.isArray(query.data.data)) return [];

    return query.data.data.map((item) => {
      return { ...item, feature: item.features[0] };
    });
  }, [query.data]);

  return { ...query, data };
}

export function useListPlanItems(
  planId: string,
  options?: {
    limit?: number;
    pricing_model?: string;
    isMetered?: boolean;
  }
) {
  const query = useMemo(() => {
    let q = `?plan_id=${planId}`;
    if (options?.limit) q += `&limit=${options.limit}`;
    if (options?.pricing_model) q += `&pricing_model=${options.pricing_model}`;
    if (options?.isMetered !== undefined)
      q += `&is_metered=${options.isMetered}`;
    return q;
  }, [options, planId]);
  const path = `${PLAN_ITEMS_PATH}${query}`;
  return useAppQuery<{ data: PlanItem[] }>({ path, queryKey: [path] });
}

export function useFirstFixedPlanItem(planId: string) {
  const { data } = useListPlanItems(planId, {
    pricing_model: "fixed",
    limit: 1,
  });
  const planItem = useMemo(() => {
    if (!data || !Array.isArray(data.data) || data.data.length === 0)
      return null;
    return data.data[0];
  }, [data]);
  return { planItem };
}

export function useCreatePlanItem(planId: string) {
  return useAppMutation<any, any>({
    method: "post",
    path: PLAN_ITEMS_PATH,
    invalidate: [`${PLAN_ITEMS_PATH}?plan_id=${planId}&is_metered=true`],
  });
}

export function useUpdatePlanItem(planItemId: string) {
  const { id } = useParams();
  return useAppMutation<any, any>({
    method: "post",
    path: `${PLAN_ITEMS_PATH}/${planItemId}`,
    invalidate: `${PLAN_ITEMS_PATH}?plan_id=${id}&is_metered=true`,
    onError: () => toast("Failed to update metered feature"),
    onSuccess: () => {
      window.dispatchEvent(new Event("saved"));
    },
  });
}

export function useDeletePlanItem(planItemId: string) {
  const { id } = useParams();
  return useAppMutation<any, any>({
    method: "delete",
    path: `${PLAN_ITEMS_PATH}/${planItemId}`,
    invalidate: [`${PLAN_ITEMS_PATH}?plan_id=${id}&is_metered=true`],
    onError: () => toast("Failed to delete metered feature"),
    onSuccess: () => toast("Metered feature deleted"),
  });
}
