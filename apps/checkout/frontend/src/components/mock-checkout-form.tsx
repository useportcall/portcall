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
          <Input
            placeholder="1234 1234 1234 1234"
            className="bg-white text-sm"
          />
          <span className="flex w-full justify-between gap-2">
            <Input placeholder="MM / YY" className="bg-white text-sm" />
            <Input placeholder="CVC" className="bg-white text-sm" />
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
            />
            <FormInput control={form.control} name="city" title="City" />
          </div>
          <Totals session={session} />
        </MockCheckoutForm>
      </div>
    </>
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
