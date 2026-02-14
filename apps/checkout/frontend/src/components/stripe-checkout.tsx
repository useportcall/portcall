import { CheckoutSessionCredentials } from "@/hooks/checkout-session-params";
import { CheckoutSession } from "@/types/api";
import { CheckoutFormSchema } from "@/types/schema";
import { Elements } from "@stripe/react-stripe-js";
import { loadStripe, StripeElementsOptions } from "@stripe/stripe-js";
import { useMemo } from "react";
import { UseFormReturn } from "react-hook-form";
import { StripeCheckoutForm } from "./stripe-checkout-form";

type Props = {
  form: UseFormReturn<CheckoutFormSchema>;
  session: CheckoutSession;
  credentials: CheckoutSessionCredentials;
};

export function StripeCheckout(props: Props) {
  const stripePromise = useMemo(() => {
    if (!props.session.external_public_key) return null;
    return loadStripe(props.session.external_public_key);
  }, [props.session.external_public_key]);

  const options: StripeElementsOptions = {
    clientSecret: props.session.external_client_secret,
    appearance: {},
  };

  return (
    <Elements stripe={stripePromise} options={options}>
      <StripeCheckoutForm
        form={props.form}
        session={props.session}
        credentials={props.credentials}
      />
    </Elements>
  );
}
