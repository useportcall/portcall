import { useMutation } from "@tanstack/react-query";
import axios from "axios";
import { CheckoutSessionCredentials } from "./checkout-session-params";

export function useUpdateBillingAddress(
  credentials: CheckoutSessionCredentials | null,
) {
  return useMutation({
    mutationFn: async (address: {
      line1: string;
      line2?: string;
      city: string;
      state?: string;
      postal_code: string;
      country: string;
    }) => {
      if (!credentials) {
        throw new Error("missing checkout session credentials");
      }

      const url = "/api/checkout-sessions/" + credentials.id + "/address";

      const { data } = await axios.post(url, address, {
        headers: {
          "Content-Type": "application/json",
          "X-Checkout-Session-Token": credentials.token,
          "Cache-Control": "no-store",
        },
      });

      return data.data;
    },
  });
}
