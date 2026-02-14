import { UseFormReturn } from "react-hook-form";
import { CheckoutFormSchema } from "@/types/schema";

const defaults = {
  cardNumber: "4242 4242 4242 4242",
  expiry: "12/30",
  cvc: "123",
  country: "US",
  line1: "1 Test Blvd",
  line2: "",
  city: "San Francisco",
  postalCode: "94105",
};

export function applyMockCheckoutAutofill(form: UseFormReturn<CheckoutFormSchema>) {
  form.setValue("country", defaults.country, { shouldValidate: true });
  form.setValue("line1", defaults.line1, { shouldValidate: true });
  form.setValue("line2", defaults.line2, { shouldValidate: true });
  form.setValue("city", defaults.city, { shouldValidate: true });
  form.setValue("postal_code", defaults.postalCode, { shouldValidate: true });
  return {
    cardNumber: defaults.cardNumber,
    expiry: defaults.expiry,
    cvc: defaults.cvc,
  };
}
