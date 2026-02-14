import axios from "axios";
import { PaymentLinkCredentials } from "./checkout-session-params";

export async function redeemPaymentLink(
  credentials: PaymentLinkCredentials,
  isCancelled: () => boolean,
  setIsRedeemingLink: (value: boolean) => void,
  setIsInvalidLink: (value: boolean) => void,
  setLinkError: (value: string | null) => void,
) {
  setIsRedeemingLink(true);
  try {
    const path = "/api/payment-links/" + credentials.id + "/redeem";
    const headers: Record<string, string> = { "Cache-Control": "no-store" };
    if (credentials.token.length > 0) {
      headers["X-Payment-Link-Token"] = credentials.token;
    }
    const { data } = await axios.post<{ data: { checkout_url: string } }>(
      path,
      {},
      {
        headers,
      },
    );
    if (!isCancelled() && data?.data?.checkout_url) {
      window.location.replace(data.data.checkout_url);
      return;
    }
    if (!isCancelled()) {
      setLinkError(null);
      setIsInvalidLink(true);
    }
  } catch (error: unknown) {
    if (!isCancelled()) {
      const apiError = axios.isAxiosError<{ error?: string }>(error)
        ? error.response?.data?.error
        : null;
      setLinkError(typeof apiError === "string" ? apiError : null);
      setIsInvalidLink(true);
    }
  } finally {
    if (!isCancelled()) {
      setIsRedeemingLink(false);
    }
  }
}
