import { Form } from "@/components/ui/form";
import { CheckoutSessionCredentials } from "@/hooks/checkout-session-params";
import { completeCheckoutSession } from "@/hooks/complete-checkout-session";
import { useUpdateBillingAddress } from "@/hooks/use-update-billing-address";
import { safeRedirect } from "@/lib/redirect";
import { CheckoutSession } from "@/types/api";
import { CheckoutFormSchema } from "@/types/schema";
import { CardElement, useElements, useStripe } from "@stripe/react-stripe-js";
import { useState } from "react";
import { UseFormReturn } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { CheckoutInput } from "./checkout-input";
import { CountryInput } from "./country-input";
import { FormInput } from "./form-input";
import { Totals } from "./totals";

type Props = {
  form: UseFormReturn<CheckoutFormSchema>;
  session: CheckoutSession;
  credentials: CheckoutSessionCredentials;
};

export function StripeCheckoutForm({ form, session, credentials }: Props) {
  const { t } = useTranslation();
  const stripe = useStripe();
  const elements = useElements();
  const { mutateAsync } = useUpdateBillingAddress(credentials);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const onSubmit = async (values: CheckoutFormSchema) => {
    if (!stripe || !elements || !session?.external_client_secret) return null;
    setIsSubmitting(true);
    setErrorMessage(null);

    try {
      await mutateAsync(values);
    } catch {
      setErrorMessage(t("error.save_billing_address"));
      setIsSubmitting(false);
      return;
    }

    const { error, setupIntent } = await stripe.confirmCardSetup(
      session.external_client_secret!,
      {
        payment_method: {
          billing_details: {
            address: {
              country: values.country,
              line1: values.line1,
              line2: values.line2,
              city: values.city,
              state: values.state,
              postal_code: values.postal_code,
            },
          },
          card: elements.getElement(CardElement)!,
        },
      },
    );

    if (error) {
      setErrorMessage(error.message || t("error.payment_failed_retry"));
      setIsSubmitting(false);
    } else {
      const pmId =
        typeof setupIntent?.payment_method === "string"
          ? setupIntent.payment_method
          : setupIntent?.payment_method?.id;
      try {
        await completeCheckoutSession(credentials, pmId);
      } catch {
        // Best-effort: webhook will also trigger the flow.
      }
      safeRedirect(session.redirect_url);
    }
  };

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="w-full max-w-sm flex flex-col gap-2"
      >
        <CheckoutInput />
        <CountryInput control={form.control} name="country" />
        <FormInput
          control={form.control}
          name="line1"
          title={t("form.address_line_1")}
          placeholder={t("form.address_line_1_placeholder")}
        />
        <FormInput
          control={form.control}
          name="line2"
          title={t("form.address_line_2")}
        />
        <div className="flex w-full gap-4">
          <FormInput
            control={form.control}
            name="postal_code"
            title={t("form.postal_code")}
            placeholder={t("form.postal_code_placeholder")}
          />
          <FormInput
            control={form.control}
            name="city"
            title={t("form.city")}
            placeholder={t("form.city_placeholder")}
          />
        </div>
        <Totals
          session={session}
          isSubmitting={isSubmitting}
          isDisabled={!stripe}
          errorMessage={errorMessage}
        />
      </form>
    </Form>
  );
}
