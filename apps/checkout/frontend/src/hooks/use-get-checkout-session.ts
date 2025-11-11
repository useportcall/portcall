"use client";

import { CheckoutSession } from "@/types/api";
import { useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import axios from "axios";

export function useGetCheckoutSession() {
  const [id, setId] = useState<string | null>(null);

  useEffect(() => {
    const queryId = new URLSearchParams(window.location.search).get("id");

    if (queryId) {
      setId(queryId);
    }
  }, []);

  return useQuery({
    queryKey: ["checkout-sessions", id],
    queryFn: async () => {
      const path = "/api/checkout-sessions/" + id;

      const { data } = await axios.get<{ data: CheckoutSession }>(path);

      return data.data;
    },
    enabled: !!id,
  });
}
