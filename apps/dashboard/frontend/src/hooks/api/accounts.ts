import { useAuth } from "@/lib/keycloak/auth";
import { getKeycloak } from "@/lib/keycloak/keycloak";
import { Account } from "@/models/account";
import { useSuspenseQuery } from "@tanstack/react-query";
import axios, { AxiosInstance } from "axios";
import { useMemo } from "react";

export function useGetAccount() {
  const { token, login } = useAuth();
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
  }, [token, login]);

  return useSuspenseQuery({
    queryKey: ["/account"],
    queryFn: async () => {
      const env = (window as any).__ENV__ || {};
      if (env.E2E_MODE) {
        return {
          data: {
            email: "e2e@portcall.internal",
            first_name: "E2E",
            last_name: "Runner",
            created_at: new Date(0).toISOString(),
            updated_at: new Date(0).toISOString(),
          },
        };
      }
      const result = await client!.get<{ data: Account }>("/account");

      return result.data;
    },
    // Retry on 401 errors to give token refresh time to work
    // This prevents the error boundary from triggering on stale tokens
    retry: (failureCount, error) => {
      if (axios.isAxiosError(error) && error.response?.status === 401) {
        return failureCount < 2;
      }
      return false;
    },
    retryDelay: 1000,
  });
}
