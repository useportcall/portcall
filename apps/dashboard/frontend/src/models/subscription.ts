export type SubscriptionItem = {
  subscription_id: number;
  price_id: number;
  description: string | null;
  unit_amount: number;
  quantity: number;
  total: number;
  type: string;
  created_at: string;
  name: string;
  price: any;
};

export type Subscription = {
  id: string;
  app_id: number;
  user_id: string;
  plan_id?: string;
  status: string;
  last_reset_at?: string;
  next_reset_at: string;
  stripe_payment_method_id: string | null;
  trial_duration_days?: number;
  invoice_count?: number;
  auto_collection?: boolean;
  currency?: string;
  billing_interval?: string;
  billing_interval_count?: number;
  items: SubscriptionItem[];
  created_at: string;
  updated_at?: string;
  user: any;
  plan: any;
};
