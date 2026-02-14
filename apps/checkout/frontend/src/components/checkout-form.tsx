import { CheckoutSession } from "@/types/api";
import { CheckoutFormSchema } from "@/types/schema";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { CheckoutSessionCredentials } from "@/hooks/checkout-session-params";
import { MockCheckout } from "./mock-checkout-form";
import { StripeCheckout } from "./stripe-checkout";
import { BraintreeCheckout } from "./braintree-checkout";
import { useTranslation } from "react-i18next";

export function CheckoutForm(props: {
  session: CheckoutSession;
  credentials: CheckoutSessionCredentials;
}) {
  const { t } = useTranslation();
  const form = useForm<CheckoutFormSchema>({
    resolver: zodResolver(CheckoutFormSchema),
  });

  switch (props.session.external_provider) {
    case "stripe":
      return (
        <StripeCheckout
          form={form}
          session={props.session}
          credentials={props.credentials}
        />
      );
    case "braintree":
      return (
        <BraintreeCheckout
          form={form}
          session={props.session}
          credentials={props.credentials}
        />
      );
    case "local":
      return (
        <MockCheckout
          form={form}
          session={props.session}
          credentials={props.credentials}
        />
      );
    default:
      return <div>{t("unsupported_payment_provider")}</div>;
  }
}
