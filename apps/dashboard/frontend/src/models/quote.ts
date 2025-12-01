import { Plan } from "./plan";
import { User } from "./user";

export interface Quote {
  id: string;
  created_at: string; // ISO date string
  updated_at: string; // ISO date string
  expires_at: string | null; // ISO date string or null
  plan: Plan;
  user?: User;
  status: "draft" | "sent" | "voided" | "accepted" | "rejected";
  company_name: string | null;
  direct_checkout_enabled: boolean;
}
