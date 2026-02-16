import { useApiClient } from "./use-api-client";
import {
  useSuspenseQuery,
  useQuery,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";

export interface DogfoodAccount {
  id: number;
  email: string;
}

export interface DogfoodApp {
  id: number;
  public_id: string;
  name: string;
  is_live: boolean;
}

export interface DogfoodPlan {
  id: number;
  public_id: string;
  name: string;
}

export interface DogfoodPlanDetail {
  id: number;
  public_id: string;
  name: string;
  description?: string;
  is_free: boolean;
  interval: string;
  interval_count: number;
  subscriber_count: number;
  features: DogfoodPlanFeature[];
}

export interface DogfoodPlanFeature {
  feature_id: string;
  quota: number;
  interval: string;
}

export interface DogfoodFeature {
  id: number;
  public_id: string;
  is_metered: boolean;
}

export interface DogfoodSecret {
  id: number;
  public_id: string;
  key_type: string;
  created_at: string;
  disabled_at: string | null;
}

export interface DogfoodStatus {
  configured: boolean;
  account?: DogfoodAccount;
  live_app?: DogfoodApp;
  test_app?: DogfoodApp;
  has_secrets: boolean;
  plan?: DogfoodPlan;
  features?: DogfoodFeature[];
  secrets?: DogfoodSecret[];
  user_count: number;
}

export interface DogfoodUser {
  id: number;
  public_id: string;
  name: string;
  email: string;
  created_at: string;
  subscription?: DogfoodUserSubscription;
}

export interface DogfoodUserSubscription {
  id: number;
  public_id: string;
  status: string;
  plan_name?: string;
  next_reset_at: string;
}

export interface DogfoodUserDetail {
  id: number;
  public_id: string;
  name: string;
  email: string;
  created_at: string;
  subscription?: DogfoodUserSubscriptionDetail;
  entitlements: DogfoodEntitlement[];
  invoices: DogfoodInvoice[];
}

export interface DogfoodUserSubscriptionDetail {
  id: number;
  public_id: string;
  status: string;
  plan_id?: number;
  plan_name?: string;
  currency: string;
  billing_interval: string;
  billing_interval_count: number;
  last_reset_at: string;
  next_reset_at: string;
  created_at: string;
}

export interface DogfoodEntitlement {
  id: number;
  feature_public_id: string;
  usage: number;
  quota: number;
  interval: string;
  is_metered: boolean;
  last_reset_at: string;
  next_reset_at: string;
}

export interface DogfoodInvoice {
  id: number;
  public_id: string;
  status: string;
  total: number;
  currency: string;
  issued_at: string;
  due_at: string;
  paid_at: string | null;
}

export interface SetupDogfoodResponse {
  account: DogfoodAccount;
  live_app: DogfoodApp;
  test_app: DogfoodApp;
  live_secret: string;
  test_secret: string;
  plan: DogfoodPlan;
  feature: DogfoodFeature;
  k8s_updated?: boolean;
  k8s_message?: string;
}

export interface CreateFeatureRequest {
  public_id: string;
  is_metered: boolean;
}

// Hook to get dogfood status
export function useDogfoodStatus() {
  const client = useApiClient();

  return useSuspenseQuery({
    queryKey: ["dogfood", "status"],
    queryFn: async () => {
      const result = await client.get<{ data: DogfoodStatus }>(
        "/dogfood/status",
      );
      return result.data.data;
    },
  });
}

// Hook to setup dogfood account
export function useSetupDogfood() {
  const client = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (updateK8s: boolean = true) => {
      const result = await client.post<{ data: SetupDogfoodResponse }>(
        "/dogfood/setup",
        { update_k8s: updateK8s },
      );
      return result.data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["dogfood"] });
    },
  });
}

// Hook to get dogfood users
export function useDogfoodUsers(appId?: number) {
  const client = useApiClient();

  return useQuery({
    queryKey: ["dogfood", "users", appId],
    queryFn: async () => {
      if (!appId) return [];
      const result = await client.get<{ data: DogfoodUser[] }>(
        `/dogfood/apps/${appId}/users`,
      );
      return result.data.data;
    },
    enabled: !!appId,
  });
}

// Hook to get a specific dogfood user
export function useDogfoodUser(appId?: number, userId?: string) {
  const client = useApiClient();

  return useQuery({
    queryKey: ["dogfood", "users", appId, userId],
    queryFn: async () => {
      if (!appId || !userId) return null;
      const result = await client.get<{ data: DogfoodUserDetail }>(
        `/dogfood/apps/${appId}/users/${userId}`,
      );
      return result.data.data;
    },
    enabled: !!appId && !!userId,
  });
}

// Hook to get dogfood features
export function useDogfoodFeatures(appId?: number) {
  const client = useApiClient();

  return useQuery({
    queryKey: ["dogfood", "features", appId],
    queryFn: async () => {
      if (!appId) return [];
      const result = await client.get<{ data: DogfoodFeature[] }>(
        `/dogfood/apps/${appId}/features`,
      );
      return result.data.data;
    },
    enabled: !!appId,
  });
}

// Hook to create a dogfood feature
export function useCreateDogfoodFeature(appId?: number) {
  const client = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (feature: CreateFeatureRequest) => {
      if (!appId) throw new Error("App ID is required");
      const result = await client.post<{ data: DogfoodFeature }>(
        `/dogfood/apps/${appId}/features`,
        feature,
      );
      return result.data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["dogfood", "features", appId],
      });
      queryClient.invalidateQueries({ queryKey: ["dogfood", "status"] });
    },
  });
}

// Hook to refresh k8s secrets
export function useRefreshK8sSecrets() {
  const client = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      const result = await client.post<{
        data: { success: boolean; message: string };
      }>("/dogfood/refresh-k8s");
      return result.data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["dogfood"] });
    },
  });
}

// Hook to get dogfood plans
export function useDogfoodPlans(appId?: number) {
  const client = useApiClient();

  return useQuery({
    queryKey: ["dogfood", "plans", appId],
    queryFn: async () => {
      if (!appId) return [];
      const result = await client.get<{ data: DogfoodPlanDetail[] }>(
        `/dogfood/apps/${appId}/plans`,
      );
      return result.data.data;
    },
    enabled: !!appId,
  });
}

export interface FixUsersResult {
  user_id: string;
  email: string;
  name: string;
  action: string;
  old_email?: string;
  new_email?: string;
  subscribed?: boolean;
  plan_name?: string;
  error?: string;
}

export interface FixUsersResponse {
  app_id: number;
  app_name: string;
  total_users: number;
  emails_fixed: number;
  subscribed: number;
  already_ok: number;
  failed: number;
  dry_run: boolean;
  results: FixUsersResult[];
}

export interface ResetPasswordResponse {
  success: boolean;
  username: string;
  email: string;
  password: string;
  realm: string;
  message: string;
}

// Hook to fix users (emails and subscriptions)
export function useFixDogfoodUsers(appId?: number) {
  const client = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (dryRun: boolean = true) => {
      if (!appId) throw new Error("App ID required");
      const result = await client.post<{ data: FixUsersResponse }>(
        `/dogfood/apps/${appId}/fix-users`,
        { dry_run: dryRun },
      );
      return result.data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["dogfood", "users", appId] });
      queryClient.invalidateQueries({ queryKey: ["dogfood", "plans", appId] });
    },
  });
}

// Hook to reset the dogfood user password in Keycloak
export function useResetDogfoodPassword() {
  const client = useApiClient();

  return useMutation({
    mutationFn: async () => {
      const result = await client.post<{ data: ResetPasswordResponse }>(
        "/dogfood/reset-password",
      );
      return result.data.data;
    },
  });
}

export interface UserValidationResult {
  user_public_id: string;
  user_email: string;
  user_name: string;
  app_public_id?: string;
  app_name?: string;
  app_is_live?: boolean;
  account_email?: string;
  status: string;
  action?: string;
  issue?: string;
  fixed?: boolean;
  error?: string;
}

export interface ValidateUsersResponse {
  live_app_id: number;
  live_app_name: string;
  test_app_id: number;
  test_app_name: string;
  total_users: number;
  valid: number;
  missing_app: number;
  wrong_environment: number;
  email_mismatch: number;
  fixed: number;
  failed: number;
  dry_run: boolean;
  results: UserValidationResult[];
}

// Hook to validate dogfood users against their corresponding apps
export function useValidateDogfoodUsers() {
  const client = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (dryRun: boolean = true) => {
      const result = await client.post<{ data: ValidateUsersResponse }>(
        "/dogfood/validate-users",
        { dry_run: dryRun },
      );
      return result.data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["dogfood", "users"] });
      queryClient.invalidateQueries({ queryKey: ["dogfood", "status"] });
    },
  });
}
