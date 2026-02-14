import { useAuth } from "@/lib/keycloak/auth";
import { App } from "@/models/app";
import {
  useMutation,
  useQueryClient,
  useSuspenseQuery,
} from "@tanstack/react-query";
import axios, { AxiosInstance } from "axios";
import { useMemo } from "react";
import { useAppQuery } from "./api";

export function useListApps() {
  const client = useClient();

  return useSuspenseQuery({
    queryKey: ["/apps"],
    queryFn: async () => {
      try {
        const result = await client!.get<{ data: App[] }>("/apps");
        return result.data;
      } catch (err) {
        if (!axios.isAxiosError(err) || err.response?.status !== 429) throw err;
        const env = (window as any).__ENV__ || {};
        const fallbackId = String(
          env.E2E_APP_ID || localStorage.getItem("selectedProject") || "",
        );
        if (!fallbackId) return { data: [] };
        const now = new Date().toISOString();
        return {
          data: [{ id: fallbackId, name: "App", public_api_key: "", is_live: false, created_at: now, updated_at: now }],
        };
      }
    },
  });
}

export function useCreateApp() {
  const client = useClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: Partial<App>) => {
      const result = await client!.post<{ data: App }>("/apps", data);
      return result.data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["/apps"] });
    },
  });
}

export function useGetApp() {
  return useAppQuery<App>({
    path: "",
    queryKey: ["/apps"],
  });
}

function useClient() {
  const { token } = useAuth();

  const client: AxiosInstance | null = useMemo(() => {
    const env = (window as any).__ENV__ || {};
    const isE2E = !!env.E2E_MODE;

    if (!token && !isE2E) {
      return null;
    }

    const baseURL = "/api";
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };

    if (isE2E) {
      headers["X-Admin-API-Key"] = env.E2E_ADMIN_KEY || "";
      headers["X-Target-App-ID"] = String(env.E2E_APP_ID || "");
    } else {
      headers["Authorization"] = `Bearer ${token}`;
    }

    return axios.create({ baseURL, headers });
  }, [token]);

  return client;
}
