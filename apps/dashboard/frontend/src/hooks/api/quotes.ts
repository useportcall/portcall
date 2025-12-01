import { toast } from "sonner";
import { useAppMutation, useAppQuery } from "./api";
import { useNavigate } from "react-router-dom";
import { Quote } from "@/models/quote";

const PATH = "/quotes";

export function useCreateQuote() {
  const navigate = useNavigate();

  return useAppMutation<unknown, Quote>({
    method: "post",
    path: PATH,
    invalidate: [PATH],
    onSuccess: (result) => {
      toast("Quote created", {
        description: "The quote has been successfully created.",
      });

      navigate(`${PATH}/${result.data.id}`);
    },
    onError: () => {
      toast("Failed to create quote", {
        description: "Please try again later.",
      });
    },
  });
}

export function useListQuotes() {
  return useAppQuery<Quote[]>({ path: PATH, queryKey: [PATH] });
}

export function useGetQuote(id: string) {
  return useAppQuery<Quote>({
    path: `${PATH}/${id}`,
    queryKey: [`${PATH}/${id}`],
  });
}

export function useUpdateQuote(id: string) {
  return useAppMutation<
    {
      user_id?: string;
      expires_at?: Date | null;
      company_name?: string | null;
      direct_checkout_enabled?: boolean;
    },
    Quote
  >({
    method: "post",
    path: `${PATH}/${id}`,
    invalidate: `${PATH}/${id}`,
    onError: () => toast("Failed to update quote"),
    onSuccess: () => {
      window.dispatchEvent(new Event("saved"));
    },
  });
}

export function useSendQuote(id: string) {
  return useAppMutation<unknown, Quote>({
    method: "post",
    path: `${PATH}/${id}/send`,
    invalidate: `${PATH}/${id}`,
    onError: () => toast("Failed to send quote"),
    onSuccess: () => {
      window.dispatchEvent(new Event("sent"));
    },
  });
}

export function useVoidQuote(id: string) {
  return useAppMutation<unknown, Quote>({
    method: "post",
    path: `${PATH}/${id}/void`,
    invalidate: `${PATH}/${id}`,
    onError: () => toast("Failed to void quote"),
    onSuccess: () => {
      window.dispatchEvent(new Event("voided"));
    },
  });
}
