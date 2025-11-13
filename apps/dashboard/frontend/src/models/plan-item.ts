import { PlanFeature } from "./plan-feature";

export interface PlanItem {
  id: string;
  quantity: number;
  pricing_model: string;
  unit_amount: number;
  tiers?: Tier[];
  minimum?: number;
  maximum?: number;
  public_title: string;
  public_description: string;
  public_unit_label: string;
  created_at: string; // ISO date string
  updated_at: string; // ISO date string
  features: PlanFeature[];
}

export interface MeteredPlanItem extends PlanItem {
  feature: PlanFeature;
}

export interface Tier {
  [key: string]: any;
}

export interface CreatePlanItemRequest {
  plan_id: string;
  pricing_model: string;
  unit_amount: number;
  public_title: string;
  public_description: string;
  interval: string;
  quota: number;
  rollover: number;
}

export interface UpdatePlanItemRequest {
  pricing_model: string;
  quantity: number;
  unit_amount: number;
  tiers: Tier[];
  minimum?: number;
  maximum?: number;
  public_title: string;
  public_description: string;
}
