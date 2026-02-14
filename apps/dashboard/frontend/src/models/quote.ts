import { Address } from "./address";
import { Plan } from "./plan";
import { User } from "./user";

export interface Quote {
  id: string;
  url?: string | null;
  signature_url?: string | null;
  created_at: string; // ISO date string
  updated_at: string; // ISO date string
  expires_at: string | null; // ISO date string or null
  plan: Plan;
  user?: User;
  user_id: string | null;
  status: "draft" | "sent" | "voided" | "accepted" | "rejected";
  company_name: string | null;
  direct_checkout_enabled: boolean;
  recipient_address: Address;
  recipient_email: string | null;
  recipient_name: string | null;
  recipient_title: string | null;
  toc: string;
  prepared_by_email: string | null;
}
