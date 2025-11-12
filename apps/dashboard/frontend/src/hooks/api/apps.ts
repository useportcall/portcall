import { useAuth } from "@/lib/keycloak/auth";
import { App } from "@/models/app";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import axios, { AxiosInstance } from "axios";
import { useMemo } from "react";
import { useAppQuery } from "./api";

export function useListApps() {
  const client = useClient();

  return useQuery({
    queryKey: ["/apps"],
    queryFn: async () => {
      const result = await client!.get<{ data: App[] }>("/apps");

      return result.data;
    },
    enabled: !!client,
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
    if (!token) {
      return null;
    }

    const baseURL = "/api";

    const headers = {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    };

    return axios.create({ baseURL, headers });
  }, [token]);

  return client;
}
