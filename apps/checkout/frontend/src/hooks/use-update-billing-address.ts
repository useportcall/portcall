import { useMutation } from "@tanstack/react-query";
import axios from "axios";

export function useUpdateBillingAddress() {
  return useMutation({
    mutationFn: async (address: {
      line1: string;
      line2?: string;
      city: string;
      state?: string;
      postal_code: string;
      country: string;
    }) => {
      const id = new URLSearchParams(window.location.search).get("id");

      if (!id) {
        throw new Error("missing checkout session id");
      }

      const url = "/api/checkout-sessions/" + id + "/address";

      const { data } = await axios.post(url, address, {
        headers: { "Content-Type": "application/json" },
      });

      return data.data;
    },
  });
}
