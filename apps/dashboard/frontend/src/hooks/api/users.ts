import { User } from "@/models/user";
import { keepPreviousData, useQuery } from "@tanstack/react-query";
import { useMemo } from "react";
import { toast } from "sonner";
import { useAppMutation, useAppQuery, useAxiosClient } from "./api";

export type PaginatedUsers = {
  users: User[];
  page: number;
  limit: number;
  total: number;
  total_pages: number;
};

export function useListUsers(props: { email: string }) {
  const query = useMemo(() => {
    let q = "?";
    if (props.email) {
      q += `email=${encodeURIComponent(props.email)}`;
    }
    return q;
  }, [props.email]);

  const path = `/users${query}`;

  return useAppQuery<User[]>({
    path,
    queryKey: [path],
  });
}

export function useListUsersPaginated(props: {
  search: string;
  subscribed: "all" | "true" | "false";
  paymentMethodAdded: "all" | "true" | "false";
  page: number;
  limit: number;
}) {
  const query = useMemo(() => {
    const p = new URLSearchParams();
    if (props.search) p.set("q", props.search);
    if (props.subscribed !== "all") p.set("subscribed", props.subscribed);
    if (props.paymentMethodAdded !== "all") {
      p.set("payment_method_added", props.paymentMethodAdded);
    }
    p.set("page", String(props.page));
    p.set("limit", String(props.limit));
    return `?${p.toString()}`;
  }, [props.search, props.subscribed, props.paymentMethodAdded, props.page, props.limit]);

  const client = useAxiosClient();
  const path = `/users${query}`;
  return useQuery({
    queryKey: [path],
    queryFn: async () => (await client.get<{ data: PaginatedUsers }>(path)).data,
    placeholderData: keepPreviousData,
  });
}

export function useCreateUser() {
  return useAppMutation<{ email: string; name?: string }, User>({
    method: "post",
    path: `/users`,
    invalidate: ["/users"],
    onError: () => {
      toast("Failed to create user", {
        description: "Please try again later.", // TODO: update error message
      });
    },
  });
}

export function useGetUser(userId: string) {
  return useAppQuery<User>({
    path: `/users/${userId}`,
    queryKey: [`/users/${userId}`],
  });
}

export function useUpdateUser(userId: string) {
  return useAppMutation<Partial<User>, any>({
    method: "post",
    path: `/users/${userId}`,
    invalidate: [`/users/${userId}`],
  });
}

export function useDeleteUser(userId: string) {
  return useAppMutation<
    Record<string, never>,
    { deleted: boolean; id: string }
  >({
    method: "delete",
    path: `/users/${userId}`,
    invalidate: ["/users", "/users?", `/users/${userId}`],
  });
}
