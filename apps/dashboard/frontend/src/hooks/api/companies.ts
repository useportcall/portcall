import { Company } from "@/models/company";
import { useAppMutation, useAppQuery } from "./api";
import { toast } from "sonner";

export function useGetCompany() {
  return useAppQuery<Company>({
    path: "/company",
    queryKey: ["/company"],
  });
}

export function useUpdateCompany() {
  return useAppMutation<Partial<FormData>, Company>({
    method: "post",
    contentType: "multipart/form-data",
    path: "/company",
    invalidate: ["/company"],
    onSuccess: () => {
      window.dispatchEvent(new Event("saved"));
    },
    onError: () => {
      toast("Failed to update company details");
    },
  });
}
