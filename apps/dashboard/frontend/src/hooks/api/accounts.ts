import { Account } from "@/models/account";
import { useAppQuery } from "./api";

export function useGetAccount() {
  return useAppQuery<Account>({
    path: "/account",
    queryKey: ["/account"],
  });
}
