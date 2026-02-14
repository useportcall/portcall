import { toast } from "sonner";
import { useAppMutation } from "./api";
import { Address } from "@/models/address";

export function useUpdateAddress(id: string) {
  return useAppMutation<Partial<Address>, Address>({
    method: "post",
    path: "/addresses/" + id,
    invalidate: ["/company"],
    onSuccess: () => {
      window.dispatchEvent(new Event("saved"));
    },
    onError: () => {
      toast("Failed to update address");
    },
  });
}
