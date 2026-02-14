import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { ArrowLeft } from "lucide-react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router-dom";
import { UserInvoices } from "./components/user-invoices";
import { UserBillingMetersSection } from "./components/view/user-billing-meters-section";
import { UserEmailAndNameSection } from "./components/view/user-email-and-name-section";
import { UserEntitlementsSection } from "./components/view/user-entitlements-section";
import { UserMeteredFeaturesSection } from "./components/view/user-metered-features-section";
import { UserSubscriptionSection } from "./components/view/user-subscription-section";

export default function UserView() {
  const navigate = useNavigate();
  const { t } = useTranslation();

  return (
    <div className="w-full h-full p-4 lg:p-10 flex flex-col gap-6">
      <div className="flex flex-row justify-between">
        <div className="flex flex-col space-y-8">
          <Button
            variant="ghost"
            className="mb-4 p-0 self-start text-sm flex flex-row justify-start items-center gap-2"
            onClick={() => navigate("/users")}
          >
            <ArrowLeft className="w-4 h-4" /> {t("views.user.back")}
          </Button>
        </div>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-[20rem_auto_2fr] w-full h-full">
        <div className="h-full flex flex-col gap-6 px-2">
          <UserEmailAndNameSection />
          <Separator />
          <UserSubscriptionSection />
        </div>
        <Separator orientation="vertical" className="hidden md:block" />
        <div className="h-full mt-5 px-10 space-y-10">
          <div className="space-y-10">
            <UserEntitlementsSection />
            <UserMeteredFeaturesSection />
          </div>
          <Separator />
          <UserBillingMetersSection />
          <Separator />
          <UserInvoices />
        </div>
      </div>
    </div>
  );
}
