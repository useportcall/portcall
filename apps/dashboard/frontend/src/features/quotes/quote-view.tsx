import { SaveIndicator } from "@/components/save-indicator";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { useGetQuote } from "@/hooks/api/quotes";
import { ArrowLeft } from "lucide-react";
import { useTranslation } from "react-i18next";
import { useNavigate, useParams } from "react-router-dom";
import { PlanItems } from "../plans/components/plan-items";
import { PlanNameInput } from "../plans/components/plan-name-input";
import {
  DownloadSignatureButton,
  PreviewQuoteButton,
  SendQuoteButton,
  VoidQuoteButton,
} from "./components/buttons";
import { AnimatedSuspense } from "./components/animated-suspense";
import { QuoteViewSidebar } from "./components/quote-view-sidebar";
import { isQuoteReadonly } from "./utils/quote-ui";

export function QuoteView() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { id } = useParams();
  const { data: quote } = useGetQuote(id!);
  const readonly = isQuoteReadonly(quote.data.status);

  return (
    <div className="w-full h-full flex flex-col gap-4 p-4 overflow-auto">
      <Button
        variant="ghost"
        className="mb-4 p-0 self-start text-sm flex flex-row items-center gap-2"
        onClick={() => navigate("/quotes")}
      >
        <ArrowLeft className="w-4 h-4" /> {t("views.quote.back")}
      </Button>
      <div className="flex gap-4 w-full justify-between px-4">
        <AnimatedSuspense>
          <div className={readonly ? "pointer-events-none opacity-75" : ""}>
            <PlanNameInput id={quote.data.plan.id} />
          </div>
        </AnimatedSuspense>
        <div className="flex flex-row justify-end gap-2">
          <SaveIndicator />
          <DownloadSignatureButton quote={quote.data} />
          <PreviewQuoteButton id={quote.data.id} quoteURL={quote.data.url} />
          <Separator orientation="vertical" />
          <VoidQuoteButton id={quote.data.id} />
          <SendQuoteButton id={quote.data.id} />
        </div>
      </div>
      <Separator />
      <div className="grid grid-cols-1 md:grid-cols-[20rem_auto_2fr] w-full h-full">
        <div className={readonly ? "pointer-events-none opacity-75" : ""}>
          <QuoteViewSidebar quote={quote.data} />
        </div>
        <Separator orientation="vertical" className="hidden md:block" />
        <div className={readonly ? "pointer-events-none opacity-75" : ""}>
          <PlanItems id={quote.data.plan.id} pricingDisabled />
        </div>
      </div>
    </div>
  );
}
