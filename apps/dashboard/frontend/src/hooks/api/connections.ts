import { toast } from "sonner";
import { useAppMutation, useAppQuery } from "./api";

type Connection = {
  id: string;
  name: string;
  source: "local" | "stripe";
  created_at: string;
};

export function useListConnections() {
  return useAppQuery<Connection[]>({
    path: "/connections",
    queryKey: ["/connections"],
  });
}

type CreateConnectionRequest = {
  name: string;
  source: string;
  public_key: string;
  secret_key: string;
  webhook_secret?: string;
};

export function useCreateConnection() {
  return useAppMutation<CreateConnectionRequest, Connection>({
    method: "post",
    path: "/connections",
    invalidate: ["/connections"],
    onSuccess: () => {
      toast("Connection created.", {
        description: "The connection is now valid for API calls.",
      });
    },
    onError: () => {
      toast("Failed to create connection.", {
        description: "Please try again later.",
      });
    },
  });
}

export function useDeleteConnection(id: string) {
  return useAppMutation<unknown, Connection>({
    method: "delete",
    path: "/connections/" + id,
    invalidate: ["/connections"],
    onSuccess: () => {
      toast("Connection deleted.", {
        description: "The connection has been removed.",
      });
    },
    onError: () => {
      toast("Failed to remove connection.", {
        description: "Please try again later.",
      });
    },
  });
}
