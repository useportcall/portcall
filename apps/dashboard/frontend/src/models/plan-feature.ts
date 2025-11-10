export interface PlanFeature {
  id: string;
  interval: string;
  quota: number;
  rollover: number;
  feature: any;
  plan_item: any;
  created_at: string; // ISO date string
  updated_at: string; // ISO date string
}

export interface UpdatePlanFeatureRequest {
  plan_item_id: string;
  feature_id?: string;
  interval: string;
  quota: number;
  rollover?: number;
}

export interface CreatePlanFeatureRequest {
  feature_id: string;
  plan_id: string;
}
