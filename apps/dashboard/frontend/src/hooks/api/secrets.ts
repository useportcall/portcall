import { Secret } from "@/models/secret";
import { useAppMutation, useAppQuery } from "./api";
import { toast } from "sonner";

export function useListSecrets() {
  return useAppQuery<Secret[]>({
    path: "/secrets",
    queryKey: ["/secrets"],
  });
}

export function useCreateSecret() {
  return useAppMutation<unknown, { key: string }>({
    method: "post",
    path: "/secrets",
    invalidate: ["/secrets"],
    onSuccess: () => {
      toast("Secret created.", {
        description: "The secret is now valid for API calls.",
      });
    },
    onError: () => {
      toast("Failed to create secret.", {
        description: "Please try again later.",
      });
    },
  });
}

export function useDisableSecret(id: string) {
  return useAppMutation<unknown, Secret>({
    method: "post",
    path: "/secrets/" + id + "/disable",
    invalidate: ["/secrets"],
    onSuccess: () => {
      toast("Secret disabled.", {
        description: "The secret is no longer valid for API calls.",
      });
    },
    onError: () => {
      toast("Failed to disable secret.", {
        description: "Please try again later.",
      });
    },
  });
}
