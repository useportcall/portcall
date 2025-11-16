"use client";

import { Form } from "@/components/ui/form";
import { useUpdateBillingAddress } from "@/hooks/use-update-billing-address";
import { CheckoutSession } from "@/types/api";
import { CheckoutFormSchema } from "@/types/schema";
import { UseFormReturn } from "react-hook-form";
import { CountryInput } from "./country-input";
import { FormInput } from "./form-input";
import { Totals } from "./totals";
import { Input } from "./ui/input";
import { useState } from "react";
import ExpiryInput from "./expiry-input";

type Props = {
  form: UseFormReturn<CheckoutFormSchema>;
  session: CheckoutSession;
};

export function MockCheckout({ form, session }: Props) {
  return (
    <>
      <div className="w-full flex flex-col gap-4 mb-2 md:max-w-sm mx-auto">
        <p className="text-sm font-medium">Card details</p>
        <div className="flex flex-col gap-2">
          <CreditCardInput />
          <span className="flex w-full justify-between gap-2">
            <ExpiryInput />
            <Input
              placeholder="CVC"
              className="bg-white text-sm"
              maxLength={4}
            />
          </span>
        </div>
      </div>
      <div className="w-full flex flex-col gap-4 md:max-w-sm mx-auto">
        <p className="text-sm font-medium">Billing details</p>
        <MockCheckoutForm form={form} session={session}>
          {/* TODO: [For after launch] add option to add business details and sales tax/VAT number */}
          <CountryInput control={form.control} name="country" />
          <FormInput
            control={form.control}
            name="line1"
            title="Address line 1"
            placeholder="1 Oxford Street"
          />
          <FormInput
            control={form.control}
            name="line2"
            title="Address line 2"
          />
          <div className="flex w-full gap-4">
            <FormInput
              control={form.control}
              name="postal_code"
              title="Postal Code"
              placeholder="WC1A 1NU"
            />
            <FormInput
              control={form.control}
              name="city"
              title="City"
              placeholder="London"
            />
          </div>
          <Totals session={session} />
        </MockCheckoutForm>
      </div>
    </>
  );
}

function CreditCardInput() {
  const [cardNumber, setCardNumber] = useState("");

  const isAmex =
    cardNumber.replace(/\s+/g, "").startsWith("34") ||
    cardNumber.replace(/\s+/g, "").startsWith("37");

  const formatCardNumber = (value: string) => {
    const digitsOnly = value.replace(/\D/g, "").slice(0, 16);

    if (isAmex) {
      const truncated = digitsOnly.slice(0, 15);
      return truncated.replace(/^(\d{4})(\d{6})(\d{5})$/, "$1 $2 $3");
    }

    const truncated = digitsOnly.slice(0, 16);
    return truncated.match(/.{1,4}/g)?.join(" ") || truncated;
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const formattedValue = formatCardNumber(e.target.value);
    setCardNumber(formattedValue);
  };

  return (
    <Input
      type="text"
      value={cardNumber}
      onChange={handleChange}
      placeholder="1234 5678 9012 3456"
      maxLength={isAmex ? 17 : 19} // 16 digits + 3 spaces
      className="bg-white text-sm"
    />
  );
}

function MockCheckoutForm({
  form,
  session,
  children,
}: {
  form: UseFormReturn<CheckoutFormSchema>;
  session: CheckoutSession;
  children: React.ReactNode;
}) {
  const { mutateAsync } = useUpdateBillingAddress();

  const onSubmit = async (data: CheckoutFormSchema) => {
    await mutateAsync(data);

    window.location.href = session.redirect_url || "/success";
  };

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="w-full md:max-w-sm mx-auto flex flex-col gap-2 mb-2"
      >
        {children}
      </form>
    </Form>
  );
}
