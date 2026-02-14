import type { Address } from "./address";

export type Company = {
  name: string;
  alias: string;
  first_name: string;
  last_name: string;
  email: string;
  phone: string;
  vat_number: string;
  business_category: string;
  billing_address: Address;
  shipping_address: Address | null;
  created_at: string;
  updated_at: string;
  icon_logo_url: string | null;
};
