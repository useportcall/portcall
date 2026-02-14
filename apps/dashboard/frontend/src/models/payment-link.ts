export interface PaymentLink {
  id: string;
  created_at: string;
  expires_at: string;
  status: string;
  plan_id: string;
  user_id: string;
  redirect_url?: string | null;
  cancel_url?: string | null;
  require_billing_address: boolean;
  url: string;
}
