import axios from "axios";
import { CheckoutSessionCredentials } from "./checkout-session-params";

/**
 * Calls the checkout complete endpoint to trigger subscription
 * creation via the billing queue. Called after card setup (Stripe)
 * or form submission (local/mock).
 *
 * For Stripe, completion is webhook-driven from setup_intent.succeeded
 * to ensure trusted SCA confirmation before subscription activation.
 * For local/mock providers, the endpoint enqueues resolve_checkout_session.
 */
export async function completeCheckoutSession(
  credentials: CheckoutSessionCredentials,
  paymentMethodId?: string,
): Promise<void> {
  const url = "/api/checkout-sessions/" + credentials.id + "/complete";
  await axios.post(
    url,
    { payment_method_id: paymentMethodId ?? "" },
    {
      headers: {
        "X-Checkout-Session-Token": credentials.token,
        "Cache-Control": "no-store",
      },
    },
  );
}
