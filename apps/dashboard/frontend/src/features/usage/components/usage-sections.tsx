import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { BillingAddressCard } from "@/features/usage/billing-address-card";
import { useBillingInvoices, useDowngradeToFree, useUpgradeToPro, useUserBillingSubscription, useUserBillingSubscriptionQuota, useUserBillingUserQuota } from "@/hooks/api/quota";
import { CreditCard, FileText, Users } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { BillingInvoicesTable } from "./billing-invoices-table";
import { CurrentPlanCard } from "./current-plan-card";
import { QuotaCard } from "./quota-card";

export function UserQuotaCard() {
  const { t } = useTranslation();
  const { data } = useUserBillingUserQuota();
  return <QuotaCard title={t("views.usage.quotas.users.title")} description={t("views.usage.quotas.users.description")} icon={Users} quota={data?.data} />;
}

export function SubscriptionQuotaCard() {
  const { t } = useTranslation();
  const { data } = useUserBillingSubscriptionQuota();
  return <QuotaCard title={t("views.usage.quotas.subscriptions.title")} description={t("views.usage.quotas.subscriptions.description")} icon={CreditCard} quota={data?.data} />;
}

export function CurrentPlanSection({ hasCheckoutQueryParams }: { hasCheckoutQueryParams: boolean }) {
  const { data: subscription } = useUserBillingSubscription();
  const upgradeMutation = useUpgradeToPro();
  const downgradeMutation = useDowngradeToFree();
  const [downgradeDialogOpen, setDowngradeDialogOpen] = useState(false);
  const [upgradeDialogOpen, setUpgradeDialogOpen] = useState(false);

  return (
    <CurrentPlanCard
      subscription={subscription?.data}
      upgradeDialogOpen={upgradeDialogOpen}
      setUpgradeDialogOpen={setUpgradeDialogOpen}
      downgradeDialogOpen={downgradeDialogOpen}
      setDowngradeDialogOpen={setDowngradeDialogOpen}
      handleUpgrade={() => {
        const currentUrl = window.location.origin + window.location.pathname;
        upgradeMutation.mutate({ cancel_url: currentUrl + "?upgrade=cancelled", redirect_url: currentUrl + "?upgrade=success" });
        setUpgradeDialogOpen(false);
      }}
      handleDowngrade={() => {
        downgradeMutation.mutate({});
        setDowngradeDialogOpen(false);
      }}
      upgradeMutation={upgradeMutation}
      downgradeMutation={downgradeMutation}
      hasCheckoutQueryParams={hasCheckoutQueryParams}
    />
  );
}

export function BillingAddressSection() {
  return <BillingAddressCard />;
}

export function InvoicesSection() {
  const { t } = useTranslation();
  const { data: invoicesData } = useBillingInvoices();
  return (
    <Card>
      <CardHeader>
        <div className="flex items-center gap-2">
          <FileText className="h-5 w-5 text-primary" />
          <CardTitle>{t("views.usage.invoices.title")}</CardTitle>
        </div>
        <CardDescription>{t("views.usage.invoices.description")}</CardDescription>
      </CardHeader>
      <CardContent>
        <BillingInvoicesTable invoices={invoicesData?.data?.invoices ?? []} />
      </CardContent>
    </Card>
  );
}
