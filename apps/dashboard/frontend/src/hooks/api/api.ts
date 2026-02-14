import { useAuth } from "@/lib/keycloak/auth";
import { getKeycloak } from "@/lib/keycloak/keycloak";
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
  contentType?: "multipart/form-data";
  path: string;
  invalidate: string | string[];
  onSuccess?: (arg0: { data: V }) => void;
  onError?: () => void;
}) {
  const client = useAxiosClient({ contentType: props.contentType });
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

export function useAxiosClient(options?: {
  contentType?: "multipart/form-data";
}) {
  const { token, login } = useAuth();
  const app = useApp();

  const client: AxiosInstance = useMemo(() => {
    const env = (window as any).__ENV__ || {};
    const isE2E = !!env.E2E_MODE;
    const baseURL = "/api/apps/" + app.id;

    const headers: Record<string, string> = {
      "Content-Type": options?.contentType || "application/json",
    };
    if (isE2E) {
      headers["X-Admin-API-Key"] = env.E2E_ADMIN_KEY || "";
      headers["X-Target-App-ID"] = String(env.E2E_APP_ID || app.id);
    } else {
      headers["Authorization"] = `Bearer ${token}`;
    }

    const instance = axios.create({ baseURL, headers });
    if (!isE2E) {
      instance.interceptors.request.use(
        async (config) => {
          try {
            const kc = await getKeycloak();
            await kc.updateToken(5);
            config.headers.Authorization = `Bearer ${kc.token}`;
          } catch (error) {
            login();
            throw error;
          }
          return config;
        },
        (error) => Promise.reject(error),
      );
    }

    return instance;
  }, [app.id, token, options?.contentType, login]);

  return client;
}
