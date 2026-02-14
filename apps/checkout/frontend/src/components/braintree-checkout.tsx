import { CheckoutSessionCredentials } from "@/hooks/checkout-session-params";
import { completeCheckoutSession } from "@/hooks/complete-checkout-session";
import { useUpdateBillingAddress } from "@/hooks/use-update-billing-address";
import { safeRedirect } from "@/lib/redirect";
import { CheckoutSession } from "@/types/api";
import { CheckoutFormSchema } from "@/types/schema";
import { useState, useEffect, useRef, useCallback } from "react";
import { UseFormReturn } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { Form } from "@/components/ui/form";
import { CountryInput } from "./country-input";
import { FormInput } from "./form-input";
import { Totals } from "./totals";

type Props = {
  form: UseFormReturn<CheckoutFormSchema>;
  session: CheckoutSession;
  credentials: CheckoutSessionCredentials;
};

export function BraintreeCheckout({ form, session, credentials }: Props) {
  const { t } = useTranslation();
  const { mutateAsync } = useUpdateBillingAddress(credentials);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const dropinRef = useRef<HTMLDivElement>(null);
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const instanceRef = useRef<any>(null);

  useEffect(() => {
    if (!session.external_client_secret || !dropinRef.current) return;
    let cancelled = false;
    import("braintree-web-drop-in").then((mod) => {
      if (cancelled) return;
      mod
        .create({
          authorization: session.external_client_secret,
          container: dropinRef.current!,
          card: { vault: { vaultCard: true } },
          // eslint-disable-next-line @typescript-eslint/no-explicit-any
        } as any)
        .then((inst) => {
          if (!cancelled) instanceRef.current = inst;
        });
    });
    return () => {
      cancelled = true;
      instanceRef.current?.teardown();
    };
  }, [session.external_client_secret]);

  const onSubmit = useCallback(
    async (values: CheckoutFormSchema) => {
      if (!instanceRef.current) return;
      setIsSubmitting(true);
      setErrorMessage(null);
      try {
        await mutateAsync(values);
      } catch {
        setErrorMessage(t("error.save_billing_address"));
        setIsSubmitting(false);
        return;
      }
      try {
        const { nonce } = await instanceRef.current.requestPaymentMethod();
        await completeCheckoutSession(credentials, nonce);
        safeRedirect(session.redirect_url);
      } catch (err: unknown) {
        setErrorMessage(
          err instanceof Error ? err.message : t("error.payment_failed_retry"),
        );
        setIsSubmitting(false);
      }
    },
    [credentials, mutateAsync, session.redirect_url, t],
  );

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="w-full max-w-sm flex flex-col gap-2"
      >
        <div ref={dropinRef} data-testid="braintree-dropin" />
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
          isDisabled={false}
          errorMessage={errorMessage}
        />
      </form>
    </Form>
  );
}
