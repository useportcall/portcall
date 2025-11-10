export type Feature = {
  id: string;
  is_metered: boolean;
  created_at: string;
  updated_at: string;
};

export type CreateFeatureRequest = {
  feature_id: string;
  is_metered: boolean;
  plan_id?: string;
  plan_feature_id?: string;
};
