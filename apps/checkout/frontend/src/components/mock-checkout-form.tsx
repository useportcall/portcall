import { CheckoutSessionCredentials } from "@/hooks/checkout-session-params";
import { CheckoutSession } from "@/types/api";
import { CheckoutFormSchema } from "@/types/schema";
import { UseFormReturn } from "react-hook-form";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { SubmitStateContext, SubmitState } from "./submit-state";
import { CountryInput } from "./country-input";
import { FormInput } from "./form-input";
import { MockCheckoutTotals } from "./mock-checkout-totals";
import { MockCheckoutFormWrapper } from "./mock-checkout-form-wrapper";
import { applyMockCheckoutAutofill } from "./mock-checkout-autofill";
import { MockCardFields } from "./mock-card-fields";
import { CheckoutComplianceNotice } from "./checkout-compliance-notice";

type Props = {
  form: UseFormReturn<CheckoutFormSchema>;
  session: CheckoutSession;
  credentials: CheckoutSessionCredentials;
};

export function MockCheckout({ form, session, credentials }: Props) {
  const { t } = useTranslation();
  const [submitState, setSubmitState] = useState<SubmitState>("idle");
  const [cardNumber, setCardNumber] = useState("");
  const [expiry, setExpiry] = useState("");
  const [cvc, setCvc] = useState("");
  const onAutofill = () => {
    const values = applyMockCheckoutAutofill(form);
    setCardNumber(values.cardNumber);
    setExpiry(values.expiry);
    setCvc(values.cvc);
  };

  return (
    <SubmitStateContext.Provider
      value={{ state: submitState, setState: setSubmitState }}
    >
      <MockCardFields
        submitState={submitState}
        cardNumber={cardNumber}
        setCardNumber={setCardNumber}
        expiry={expiry}
        setExpiry={setExpiry}
        cvc={cvc}
        setCvc={setCvc}
        onAutofill={onAutofill}
      />
      <div className="w-full flex flex-col gap-4 md:max-w-sm mx-auto">
        <p className="text-sm font-medium">{t("form.billing_details")}</p>
        <MockCheckoutFormWrapper
          form={form}
          session={session}
          credentials={credentials}
        >
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
          <MockCheckoutTotals session={session} />
          <CheckoutComplianceNotice session={session} />
        </MockCheckoutFormWrapper>
      </div>
    </SubmitStateContext.Provider>
  );
}
