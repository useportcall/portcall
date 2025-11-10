import { PlanItem } from "./plan-item";

export interface Plan {
  id: string;
  name: string;
  currency: string;
  status: string;
  trial_period_days: number;
  created_at: string; // ISO date string
  updated_at: string; // ISO date string
  items: PlanItem[];
  interval: string;
  interval_count: number;
  plan_group: any;
  features: any[];
  metered_features: any[];
}

export interface UpdatePlanRequest {
  name: string;
  currency: string;
  trial_period_days: number;
  interval: string;
  interval_count: number;
  plan_group_id: string;
}
