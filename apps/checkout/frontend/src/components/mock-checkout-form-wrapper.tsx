import { completeCheckoutSession } from "@/hooks/complete-checkout-session";
import { Form } from "@/components/ui/form";
import { CheckoutSessionCredentials } from "@/hooks/checkout-session-params";
import { useUpdateBillingAddress } from "@/hooks/use-update-billing-address";
import { safeRedirect } from "@/lib/redirect";
import { CheckoutSession } from "@/types/api";
import { CheckoutFormSchema } from "@/types/schema";
import { UseFormReturn } from "react-hook-form";
import { useSubmitState } from "./submit-state";

type Props = {
  form: UseFormReturn<CheckoutFormSchema>;
  session: CheckoutSession;
  credentials: CheckoutSessionCredentials;
  children: React.ReactNode;
};

export function MockCheckoutFormWrapper({
  form,
  session,
  credentials,
  children,
}: Props) {
  const { mutateAsync } = useUpdateBillingAddress(credentials);
  const { setState } = useSubmitState();

  const onSubmit = async (data: CheckoutFormSchema) => {
    setState("loading");
    try {
      await mutateAsync(data);
      await completeCheckoutSession(credentials);
      setState("success");
      await new Promise((resolve) => setTimeout(resolve, 1500));
      safeRedirect(session.redirect_url);
    } catch {
      setState("idle");
    }
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
