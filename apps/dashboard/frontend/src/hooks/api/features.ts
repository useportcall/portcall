import { Feature } from "@/models/feature";
import { useParams } from "react-router-dom";
import { useAppQuery, useAppMutation } from "./api";
import { toast } from "sonner";

export function useListFeatures(props: { isMetered: boolean }) {
  return useAppQuery<Feature[]>({
    path: "/features?is_metered=" + props.isMetered,
    queryKey: ["/features?is_metered=" + props.isMetered],
  });
}

export function useCreateFeature(props: { isMetered: boolean }) {
  const { id } = useParams();
  return useAppMutation<any, any>({
    method: "post",
    path: "/features",
    invalidate: [
      `/plan-items?plan_id=${id}`,
      `/features?is_metered=${props.isMetered}`,
    ],
    onError: () => toast("Failed to add feature"),
  });
}
