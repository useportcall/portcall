export type App = {
  id: string;
  name: string;
  public_api_key: string;
  is_live: boolean;
  billing_exempt?: boolean;
  created_at: string;
  updated_at: string;
};

export type CreateAppRequest = {
  name: string;
};
