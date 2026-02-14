import { Quote } from "@/models/quote";
import { Scale } from "lucide-react";
import { DEFAULT_TERMS } from "../constants";
import { summarizeTerms } from "../utils/quote-ui";

export function TermsSummary({ quote }: { quote: Quote }) {
  return (
    <div className="px-2 flex flex-col gap-2">
      <div className="flex w-full gap-1.5">
        <Scale className="w-4 h-4" />
        <p className="text-sm">{summarizeTerms(quote.toc || DEFAULT_TERMS)}</p>
      </div>
    </div>
  );
}
