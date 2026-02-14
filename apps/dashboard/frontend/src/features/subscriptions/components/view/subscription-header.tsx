import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router-dom";

export function SubscriptionHeader({ email }: { email: string }) {
  const navigate = useNavigate();
  const { t } = useTranslation();

  return (
    <div className="flex flex-row justify-between">
      <div className="flex flex-col space-y-8">
        <Button
          variant="ghost"
          className="mb-4 p-0 self-start text-sm flex flex-row justify-start items-center gap-2"
          onClick={() => navigate(-1)}
        >
          <ArrowLeft className="w-4 h-4" /> {t("views.subscription.back")}
        </Button>
        <div className="flex flex-col justify-start space-y-2">
          <p className="text-lg md:text-xl font-semibold space-mono-bold">{t("views.subscription.heading", { email })}</p>
          <p className="text-sm text-muted-foreground">{t("views.subscription.description")}</p>
        </div>
      </div>
    </div>
  );
}
