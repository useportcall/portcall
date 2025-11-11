import { Form } from "@/components/ui/form";
import { useUpdateBillingAddress } from "@/hooks/use-update-billing-address";
import { CheckoutSession } from "@/types/api";
import { CheckoutFormSchema } from "@/types/schema";
import {
  CardElement,
  Elements,
  useElements,
  useStripe,
} from "@stripe/react-stripe-js";
import { loadStripe, StripeElementsOptions } from "@stripe/stripe-js";
import { Lock } from "lucide-react";
import { useMemo } from "react";
import { UseFormReturn } from "react-hook-form";
import { CheckoutInput } from "./checkout-input";
import { CountryInput } from "./country-input";
import { FormInput } from "./form-input";
import { Totals } from "./totals";

type Props = {
  form: UseFormReturn<CheckoutFormSchema>;
  session: CheckoutSession;
};

export function StripeCheckout(props: Props) {
  const stripePromise = useMemo(() => {
    console.log(
      "Loading stripe with public id:",
      props.session.external_public_key
    );
    if (!props.session.external_public_key) return null;

    return loadStripe(props.session.external_public_key);
  }, [props.session.external_public_key]);

  const options: StripeElementsOptions = {
    clientSecret: props.session.external_client_secret,
    appearance: {},
  };

  return (
    <Elements stripe={stripePromise} options={options}>
      <StripeCheckoutForm form={props.form} session={props.session} />
    </Elements>
  );
}
function StripeCheckoutForm({ form, session }: Props) {
  const stripe = useStripe();
  const elements = useElements();
  const { mutateAsync } = useUpdateBillingAddress();

  const onSubmit = async (values: CheckoutFormSchema) => {
    if (!stripe || !elements || !session?.external_client_secret) {
      return null;
    }

    try {
      await mutateAsync({
        line1: values.line1,
        line2: values.line2,
        city: values.city,
        state: values.state,
        postal_code: values.postal_code,
        country: values.country,
      });
    } catch (error) {
      console.error("Error setting address:", error);
      return;
    }

    const { error } = await stripe.confirmCardSetup(
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
      }
    );

    if (error) {
      console.error("Error confirming card setup:", error);
    } else {
      window.location.href = session.redirect_url;
    }
  };

  return (
    <Form {...form}>
      <form
        onSubmit={(e) => {
          console.log("valid:", form.formState.errors, form.getValues());
          form.handleSubmit(onSubmit)(e);
        }}
        className="w-full max-w-sm flex flex-col gap-2"
      >
        <CheckoutInput />
        <CountryInput control={form.control} name="country" />
        <FormInput
          control={form.control}
          name="line1"
          title="Address line 1"
          placeholder="1 Oxford Street"
        />
        <FormInput control={form.control} name="line2" title="Address line 2" />
        <div className="flex w-full gap-4">
          <FormInput
            control={form.control}
            name="postal_code"
            title="Postal Code"
          />
          <FormInput control={form.control} name="city" title="City" />
        </div>
        <Totals session={session} />
        <button className="w-full bg-slate-900 p-2 flex justify-center items-center gap-2 text-white font-semibold mt-2 rounded-md hover:bg-slate-700 cursor-pointer">
          Pay $0.00
          <span className="w-4 h-4">
            <Lock className="w-4 h-4" />
          </span>
        </button>
      </form>
    </Form>
  );
}
