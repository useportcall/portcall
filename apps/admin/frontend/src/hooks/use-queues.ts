import { useApiClient } from "./use-api-client";
import {
  useSuspenseQuery,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";

export interface ParamInfo {
  name: string;
  type: string;
  required: boolean;
}

export interface TaskInfo {
  name: string;
  description: string;
  params: ParamInfo[];
}

export interface QueueInfo {
  name: string;
  tasks: TaskInfo[];
}

export interface TaskPayload {
  queue_name: string;
  task_type: string;
  payload: Record<string, unknown>;
}

export function useQueues() {
  const client = useApiClient();

  return useSuspenseQuery({
    queryKey: ["queues"],
    queryFn: async () => {
      const result = await client.get<{ data: QueueInfo[] }>("/queues");
      return result.data.data;
    },
  });
}

export function useEnqueueTask() {
  const client = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (payload: TaskPayload) => {
      const result = await client.post("/queues/enqueue", payload);
      return result.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["queues"] });
    },
  });
}
