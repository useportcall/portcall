export interface QuotaStatus {
  feature_id: string;
  usage: number;
  quota: number;
  remaining: number;
  is_exceeded: boolean;
  is_unlimited: boolean;
}

export interface AllQuotasResponse {
  subscriptions?: QuotaStatus;
  users?: QuotaStatus;
}

export interface BillingPlan {
  id: string;
  name: string;
}

export interface UserBillingSubscription {
  plan_name: string;
  is_free: boolean;
  has_payment_method: boolean;
  next_reset_at?: string;
  scheduled_plan_id?: string;
  scheduled_plan?: BillingPlan;
  current_plan?: BillingPlan;
}

export interface UpgradeToProResponse {
  checkout_url?: string;
  session_id?: string;
  success: boolean;
}

export interface DowngradeToFreeResponse {
  success: boolean;
  scheduled: boolean;
  message?: string;
}

export interface DowngradeToFreeRequest {
  immediate?: boolean;
}

export interface BillingInvoice {
  id: string;
  invoice_number: string;
  currency: string;
  total: number;
  status: string;
  pdf_url?: string;
  created_at: string;
}

export interface BillingInvoicesResponse {
  invoices: BillingInvoice[];
}

export interface BillingAddress {
  id: string;
  line1: string;
  line2?: string;
  city: string;
  state?: string;
  postal_code: string;
  country: string;
}

export interface BillingAddressResponse {
  billing_address: BillingAddress | null;
}

export interface UpsertBillingAddressRequest {
  line1: string;
  line2?: string;
  city: string;
  state?: string;
  postal_code: string;
  country: string;
}
