export type Entitlement = {
  id: string;
  usage: number;
  quota: number;
  interval: string;
  next_reset_at: string | null;
  feature: any;
  created_at: string;
  updated_at: string;
};
