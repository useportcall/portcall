import { useApiClient } from "./use-api-client";
import {
  useSuspenseQuery,
  useMutation,
  useQueryClient,
  useQuery,
} from "@tanstack/react-query";

export interface TaskInfo {
  id: string;
  queue: string;
  type: string;
  payload: string;
  state: string;
  max_retry: number;
  retried: number;
  last_err: string;
  last_failed_at?: string;
  next_process_at?: string;
  completed_at?: string;
}

export interface QueueStats {
  queue: string;
  active: number;
  pending: number;
  scheduled: number;
  retry: number;
  archived: number;
  completed: number;
}

export interface ListTasksParams {
  queue: string;
  state?: string;
  taskType?: string;
  limit?: number;
}

export interface ListTasksResponse {
  tasks: TaskInfo[];
  queue: string;
  state: string;
  limit: number;
}

export function useQueueStats() {
  const client = useApiClient();

  return useQuery({
    queryKey: ["queue-stats"],
    queryFn: async () => {
      const result = await client.get<{ data: { stats: QueueStats[] } }>(
        "/queues/stats",
      );
      return result.data.data.stats;
    },
    refetchInterval: 10000, // Refresh every 10 seconds
  });
}

export function useTasks(params: ListTasksParams) {
  const client = useApiClient();

  return useSuspenseQuery({
    queryKey: [
      "queue-tasks",
      params.queue,
      params.state,
      params.taskType,
      params.limit,
    ],
    queryFn: async () => {
      const searchParams = new URLSearchParams();
      searchParams.set("queue", params.queue);
      if (params.state) searchParams.set("state", params.state);
      if (params.taskType) searchParams.set("task_type", params.taskType);
      if (params.limit) searchParams.set("limit", params.limit.toString());

      const result = await client.get<{ data: ListTasksResponse }>(
        `/queues/tasks?${searchParams.toString()}`,
      );
      return result.data.data;
    },
  });
}

export function useTask(queue: string, taskId: string) {
  const client = useApiClient();

  return useQuery({
    queryKey: ["queue-task", queue, taskId],
    queryFn: async () => {
      const result = await client.get<{ data: TaskInfo }>(
        `/queues/tasks/${taskId}?queue=${queue}`,
      );
      return result.data.data;
    },
    enabled: !!queue && !!taskId,
  });
}

export function useArchiveTask() {
  const client = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      queue,
      taskId,
    }: {
      queue: string;
      taskId: string;
    }) => {
      const result = await client.post("/queues/tasks/archive", {
        queue,
        task_id: taskId,
      });
      return result.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["queue-tasks"] });
      queryClient.invalidateQueries({ queryKey: ["queue-stats"] });
    },
  });
}

export function useDeleteTask() {
  const client = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      queue,
      taskId,
    }: {
      queue: string;
      taskId: string;
    }) => {
      const result = await client.post("/queues/tasks/delete", {
        queue,
        task_id: taskId,
      });
      return result.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["queue-tasks"] });
      queryClient.invalidateQueries({ queryKey: ["queue-stats"] });
    },
  });
}

export function useRunTask() {
  const client = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      queue,
      taskId,
    }: {
      queue: string;
      taskId: string;
    }) => {
      const result = await client.post("/queues/tasks/run", {
        queue,
        task_id: taskId,
      });
      return result.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["queue-tasks"] });
      queryClient.invalidateQueries({ queryKey: ["queue-stats"] });
    },
  });
}

export function useRetryWithPayload() {
  const client = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      queue,
      taskType,
      payload,
    }: {
      queue: string;
      taskType: string;
      payload: string;
    }) => {
      const result = await client.post("/queues/tasks/retry", {
        queue,
        task_type: taskType,
        payload: JSON.parse(payload),
      });
      return result.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["queue-tasks"] });
      queryClient.invalidateQueries({ queryKey: ["queue-stats"] });
    },
  });
}
