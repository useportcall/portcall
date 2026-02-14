import { PaymentLink } from "@/models/payment-link";
import { toast } from "sonner";
import { useAppMutation } from "./api";

export interface CreatePaymentLinkRequest {
  plan_id: string;
  user_id?: string;
  user_email?: string;
  user_name?: string;
  cancel_url?: string;
  redirect_url?: string;
  expires_at?: string;
}

export function useCreatePaymentLink() {
  return useAppMutation<CreatePaymentLinkRequest, PaymentLink>({
    method: "post",
    path: `/payment-links`,
    invalidate: ["/subscriptions"],
    onSuccess: () => {
      toast("Payment link created", {
        description: "Share this link so the user can complete checkout.",
      });
    },
    onError: () => {
      toast("Failed to create payment link", {
        description: "Please try again later.",
      });
    },
  });
}
