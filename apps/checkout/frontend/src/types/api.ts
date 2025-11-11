export interface PlanItem {
  id: string;
  quantity: number;
  pricing_model: string;
  unit_amount: number;
  tiers: Tier[] | null;
  minimum: number | null;
  maximum: number | null;
  public_title: string;
  public_description: string;
  public_unit_label: string;
  created_at: string;
  updated_at: string;
}

type Tier = {
  start: number;
  end: number;
  unit_amount: number;
};

export interface Plan {
  id: string;
  name: string;
  currency: string;
  status: string;
  trial_period_days: number;
  created_at: string;
  updated_at: string;
  interval: string;
  interval_count: number;
  plan_group: { id: string } | null;
  items: PlanItem[];
  features: { id: string; is_metered: boolean }[];
  metered_features: MeteredFeature[];
}

type MeteredFeature = {
  id: string;
  quota: number;
  interval: string;
  feature: { id: string };
  plan_item: PlanItem;
};

export interface CheckoutSession {
  id: string;
  url: string;
  created_at: string;
  expires_at: string;
  external_client_secret: string;
  external_public_key: string;
  external_provider: string; // e.g., stripe, local
  external_session_id: string;
  redirect_url: string; // URL to redirect after checkout
  cancel_url: string | null; // URL to redirect if the user cancels the checkout
  plan: Plan | null;
  billing_address: Address | null;
  company_address: Address | null;
  company?: {
    id: string;
    name: string;
  };
}

export interface Address {
  line1: string;
  line2?: string;
  city: string;
  postal_code: string;
  country_code: string;
  state?: string;
}
