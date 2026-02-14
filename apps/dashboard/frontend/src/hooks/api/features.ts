import { CreateFeatureRequest, Feature } from "@/models/feature";
import { useParams } from "react-router-dom";
import { useAppQuery, useAppMutation } from "./api";
import { toast } from "sonner";
import { PLAN_ITEMS_PATH } from "./plan-items";

export function useListFeatures(props: { isMetered: boolean }) {
  return useAppQuery<Feature[]>({
    path: "/features?is_metered=" + props.isMetered,
    queryKey: ["/features?is_metered=" + props.isMetered],
  });
}

export function useCreateFeature(props: { isMetered: boolean }) {
  const { id } = useParams();
  return useAppMutation<CreateFeatureRequest, Feature>({
    method: "post",
    path: "/features",
    invalidate: [
      `${PLAN_ITEMS_PATH}?plan_id=${id}&is_metered=true`,
      `/features?is_metered=${props.isMetered}`,
    ],
    onError: () => toast("Failed to add feature"),
  });
}

export function useCreateBasicFeature() {
  const { id } = useParams();
  return useAppMutation<CreateFeatureRequest, Feature>({
    method: "post",
    path: "/features",
    invalidate: [
      `${PLAN_ITEMS_PATH}?plan_id=${id}&is_metered=true`,
      `/plan-features?plan_id=${id}`,
      `/features?is_metered=false`,
    ],
    onError: () => toast("Failed to add feature"),
  });
}
