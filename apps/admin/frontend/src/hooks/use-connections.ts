import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useApiClient } from "./use-api-client";

interface MockConnectionInfo {
  app_id: number;
  app_public_id: string;
  app_name: string;
  is_live: boolean;
  mock_connection_id?: number;
  connection_name?: string;
}

interface CheckMockConnectionsResponse {
  has_mock_connections: boolean;
  live_apps_with_mock: MockConnectionInfo[];
  total_live_apps: number;
  total_with_mock: number;
}

interface ClearMockConnectionsResponse {
  cleared: number;
  errors?: string[];
  app_ids: number[];
  messages: string[];
}

interface ConnectionInfo {
  id: number;
  public_id: string;
  name: string;
  source: string;
  public_key: string;
  is_default: boolean;
}

export function useCheckMockConnections() {
  const api = useApiClient();
  return useQuery<CheckMockConnectionsResponse>({
    queryKey: ["mock-connections-check"],
    queryFn: async () => {
      const response = await api.get("/api/connections/check-mock");
      return response.data.data;
    },
    staleTime: 30000, // 30 seconds
  });
}

export function useClearMockConnections() {
  const api = useApiClient();
  const queryClient = useQueryClient();

  return useMutation<ClearMockConnectionsResponse>({
    mutationFn: async () => {
      const response = await api.post("/api/connections/clear-mock");
      return response.data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["mock-connections-check"] });
      queryClient.invalidateQueries({ queryKey: ["apps"] });
    },
  });
}

export function useAppConnections(appId: string | number) {
  const api = useApiClient();
  return useQuery<ConnectionInfo[]>({
    queryKey: ["app-connections", appId],
    queryFn: async () => {
      const response = await api.get(`/api/apps/${appId}/connections`);
      return response.data.data;
    },
    enabled: !!appId,
  });
}
