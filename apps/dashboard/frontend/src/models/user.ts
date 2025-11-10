import type { Address } from "./address";

export type User = {
  id: string;
  name: string;
  email: string | null;
  created_at: string;
  updated_at: string;
  subscribed: boolean;
  payment_method_added: boolean;
  billing_address: Address | null;
};

export type UpdateUserRequest = {
  name: string;
};

export type CreateUserRequest = {
  email: string;
  name: string;
};
