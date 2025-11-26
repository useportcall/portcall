import { toast } from "sonner";
import { useAppMutation } from "./api";
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
