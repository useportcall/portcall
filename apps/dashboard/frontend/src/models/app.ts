export type App = {
  id: string;
  name: string;
  public_api_key: string;
  created_at: string;
  updated_at: string;
};

export type CreateAppRequest = {
  name: string;
};
