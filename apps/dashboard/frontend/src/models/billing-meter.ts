export type Tier = {
  start: number;
  end: number; // -1 means unlimited
  unit_amount: number; // in cents
};

export type BillingMeter = {
  id: string;
  subscription_id: string;
  feature_id: string;
  plan_item_id: string;
  user_id: string;
  usage: number;
  pricing_model: "unit" | "tiered" | "block" | "volume";
  unit_amount: number; // in cents
  free_quota: number;
  tiers?: Tier[];
  last_reset_at: string | null;
  next_reset_at: string | null;
  created_at: string;
  updated_at: string;

  // Computed/display fields
  projected_cost: number; // in cents
  feature_name?: string;
  plan_item_title?: string;
};
