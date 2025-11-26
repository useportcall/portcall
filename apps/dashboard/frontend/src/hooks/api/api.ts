import { useAuth } from "@/lib/keycloak/auth";
import {
  useMutation,
  useQueryClient,
  useSuspenseQuery,
} from "@tanstack/react-query";
import axios, { AxiosInstance } from "axios";
import { useMemo } from "react";
import { useApp } from "../use-app";

export function useAppQuery<T = unknown>(props: {
  path: string;
  queryKey: string[];
}) {
  const client = useAxiosClient();

  return useSuspenseQuery({
    queryKey: [props.path],
    queryFn: async () => {
      const result = await client.get<{ data: T }>(props.path);

      return result.data;
    },
  });
}

export function useAppMutation<T, V = unknown>(props: {
  method: "post" | "delete";
  path: string;
  invalidate: string | string[];
  onSuccess?: (arg0: { data: V }) => void;
  onError?: () => void;
}) {
  const client = useAxiosClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (dto: T) => {
      const result = await client![props.method]<{ data: V }>(props.path, dto);
      return result.data;
    },
    onSuccess: (arg0: { data: V }) => {
      if (Array.isArray(props.invalidate)) {
        props.invalidate.forEach((path) => {
          queryClient.invalidateQueries({ queryKey: [path] });
        });
      } else {
        queryClient.invalidateQueries({ queryKey: [props.invalidate] });
      }

      props.onSuccess?.(arg0);
    },
    onError: props.onError,
  });
}

export function useAxiosClient() {
  const { token } = useAuth();
  const app = useApp();

  const client: AxiosInstance = useMemo(() => {
    const baseURL = "/api/apps/" + app.id;

    const headers = {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    };

    return axios.create({ baseURL, headers });
  }, [app.id, token]);

  return client;
}
