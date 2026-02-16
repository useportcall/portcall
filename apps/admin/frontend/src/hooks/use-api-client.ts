import { useAuth } from "@/lib/keycloak/auth";
import axios, { AxiosInstance } from "axios";
import { useMemo } from "react";

export function useApiClient() {
  const { token } = useAuth();

  const client: AxiosInstance = useMemo(() => {
    const baseURL = "/api";
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };

    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }

    return axios.create({ baseURL, headers });
  }, [token]);

  return client;
}
