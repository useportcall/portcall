import { CheckoutSession } from "@/types/api";
import { CheckoutFormSchema } from "@/types/schema";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { MockCheckout } from "./mock-checkout-form";
import { StripeCheckout } from "./stripe-checkout";

export function CheckoutForm(props: { session: CheckoutSession }) {
  const form = useForm<CheckoutFormSchema>({
    resolver: zodResolver(CheckoutFormSchema),
  });

  switch (props.session.external_provider) {
    case "stripe":
      return <StripeCheckout form={form} session={props.session} />;
    case "local":
      return <MockCheckout form={form} session={props.session} />;
    default:
      return <div>Unsupported payment provider</div>;
  }
}
