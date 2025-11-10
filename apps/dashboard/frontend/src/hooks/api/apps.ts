import { App } from "@/models/app";
import { useAppMutation, useAppQuery } from "./api";
import { useQuery } from "@tanstack/react-query";
import axios, { AxiosInstance } from "axios";
import { useMemo } from "react";
import { useAuth } from "@/lib/keycloak/auth";

export function useListApps() {
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
  return useAppMutation<{ name: string }, App>({
    method: "post",
    path: "/apps",
    invalidate: ["/apps"],
  });
}

export function useGetApp() {
  return useAppQuery<App>({
    path: "",
    queryKey: ["/apps"],
  });
}
