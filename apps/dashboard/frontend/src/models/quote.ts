import { Plan } from "./plan";

export interface Quote {
  id: string;
  created_at: string; // ISO date string
  updated_at: string; // ISO date string
  plan: Plan;
}
