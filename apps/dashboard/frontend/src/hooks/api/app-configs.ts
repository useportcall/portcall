import { toast } from "sonner";
import { useAppMutation, useAppQuery } from "./api";

export type AppConfig = {
  default_connection: string;
};

export function useGetAppConfig() {
  return useAppQuery<AppConfig>({
    path: "/config",
    queryKey: ["/config"],
  });
}

export function useAppConfig() {
  return useAppMutation<unknown, AppConfig>({
    method: "post",
    path: "/config",
    invalidate: ["/config"],
    onError: () => {
      toast("Failed to update config.", {
        description: "Please try again later.",
      });
    },
  });
}
