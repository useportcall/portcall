import { Suspense } from "react";
import { useTranslation } from "react-i18next";
import { StatusMessages } from "./status-messages";
import { BillingAddressSection, CurrentPlanSection, InvoicesSection, SubscriptionQuotaCard, UserQuotaCard } from "./usage-sections";
import { InvoicesCardSkeleton, PlanCardSkeleton, QuotaCardSkeleton } from "./usage-skeletons";

export function UsageContent({ isLiveApp }: { isLiveApp: boolean }) {
  const { t } = useTranslation();
  const urlParams = new URLSearchParams(window.location.search);
  const upgradeStatus = urlParams.get("upgrade");
  const downgradeStatus = urlParams.get("downgrade");
  const downgradeScheduled = urlParams.get("scheduled") === "true";
  const hasCheckoutQueryParams = upgradeStatus !== null || downgradeStatus !== null;

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">{t("views.usage.title")}</h1>
        <p className="text-muted-foreground">{isLiveApp ? t("views.usage.description_live") : t("views.usage.description_test")}</p>
      </div>
      <StatusMessages upgradeStatus={upgradeStatus} downgradeStatus={downgradeStatus} downgradeScheduled={downgradeScheduled} />
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <Suspense fallback={<QuotaCardSkeleton />}><UserQuotaCard /></Suspense>
        <Suspense fallback={<QuotaCardSkeleton />}><SubscriptionQuotaCard /></Suspense>
      </div>
      {isLiveApp && <div className="grid gap-4 lg:grid-cols-[2fr_1fr]"><Suspense fallback={<PlanCardSkeleton />}><CurrentPlanSection hasCheckoutQueryParams={hasCheckoutQueryParams} /></Suspense><Suspense fallback={<PlanCardSkeleton />}><BillingAddressSection /></Suspense></div>}
      {isLiveApp && <Suspense fallback={<InvoicesCardSkeleton />}><InvoicesSection /></Suspense>}
    </div>
  );
}
