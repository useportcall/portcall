import { useAuth } from "@/lib/keycloak/auth";
import { Account } from "@/models/account";
import { useSuspenseQuery } from "@tanstack/react-query";
import axios, { AxiosInstance } from "axios";
import { useMemo } from "react";

export function useGetAccount() {
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

  return useSuspenseQuery({
    queryKey: ["/account"],
    queryFn: async () => {
      const result = await client!.get<{ data: Account }>("/account");

      return result.data;
    },
  });
}
