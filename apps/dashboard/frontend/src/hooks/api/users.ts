import { User } from "@/models/user";
import { useAppMutation, useAppQuery } from "./api";
import { useMemo } from "react";

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

export function useCreateUser() {
  return useAppMutation<{ email: string; name?: string }, User>({
    method: "post",
    path: `/users`,
    invalidate: ["/users"],
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
