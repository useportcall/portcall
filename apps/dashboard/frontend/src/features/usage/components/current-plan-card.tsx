import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { UseMutationResult } from "@tanstack/react-query";
import { Sparkles } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { DowngradeDialog } from "./downgrade-dialog";
import { PendingChangesDialog } from "./pending-changes-dialog";
import { UpgradeDialog } from "./upgrade-dialog";

interface SubscriptionData { plan_name?: string; is_free?: boolean; next_reset_at?: string; scheduled_plan?: { name: string }; current_plan?: { name: string }; }

export function CurrentPlanCard({ subscription, upgradeDialogOpen, setUpgradeDialogOpen, downgradeDialogOpen, setDowngradeDialogOpen, handleUpgrade, handleDowngrade, upgradeMutation, downgradeMutation, hasCheckoutQueryParams }: { subscription?: SubscriptionData; upgradeDialogOpen: boolean; setUpgradeDialogOpen: (open: boolean) => void; downgradeDialogOpen: boolean; setDowngradeDialogOpen: (open: boolean) => void; handleUpgrade: () => void; handleDowngrade: () => void; upgradeMutation: UseMutationResult<any, any, any, any>; downgradeMutation: UseMutationResult<any, any, any, any>; hasCheckoutQueryParams: boolean; }) {
  const { t, i18n } = useTranslation();
  const [pendingChangesDialogOpen, setPendingChangesDialogOpen] = useState(false);
  const hasPendingChange = !!subscription?.scheduled_plan;
  const nextReset = formatNextResetDate(subscription?.next_reset_at, i18n.language);

  return (
    <Card><CardHeader><div className="flex items-center gap-2"><Sparkles className="h-5 w-5 text-primary" /><CardTitle>{t("views.usage.current_plan.title")}</CardTitle></div></CardHeader><CardContent className="space-y-4"><div className="flex items-center justify-between"><div><p className="text-lg font-semibold">{subscription?.plan_name}</p><p className="text-sm text-muted-foreground">{!subscription?.is_free ? t("views.usage.current_plan.pro_description") : t("views.usage.current_plan.free_description")}</p></div><Badge variant={!subscription?.is_free ? "default" : "outline"} className={!subscription?.is_free ? "bg-linear-to-r from-amber-500 to-orange-500" : "text-primary border-primary"}>{!subscription?.is_free ? t("views.usage.current_plan.pro_badge") : t("views.usage.current_plan.active_badge")}</Badge></div>
      {subscription?.next_reset_at && <PendingChangesDialog open={pendingChangesDialogOpen} onOpenChange={setPendingChangesDialogOpen} subscription={subscription} hasPendingChange={hasPendingChange} formatNextResetDate={() => nextReset} />}
      {subscription?.is_free ? <PlanAction title={t("views.usage.current_plan.upgrade_title")} description={t("views.usage.current_plan.upgrade_description")}><UpgradeDialog open={upgradeDialogOpen} onOpenChange={setUpgradeDialogOpen} onConfirm={handleUpgrade} isPending={upgradeMutation.isPending} isDisabled={upgradeMutation.isPending || hasCheckoutQueryParams} />{hasCheckoutQueryParams && <p className="text-sm text-muted-foreground mt-2">{t("views.usage.current_plan.refresh_upgrade")}</p>}{upgradeMutation.isError && <p className="text-sm text-red-600 mt-2">{t("views.usage.current_plan.upgrade_failed")}</p>}</PlanAction> : <PlanAction title={t("views.usage.current_plan.downgrade_title")} description={t("views.usage.current_plan.downgrade_description")}><DowngradeDialog open={downgradeDialogOpen} onOpenChange={setDowngradeDialogOpen} onConfirm={handleDowngrade} isPending={downgradeMutation.isPending} />{downgradeMutation.isError && <p className="text-sm text-red-600 mt-2">{t("views.usage.current_plan.downgrade_failed")}</p>}</PlanAction>}
    </CardContent></Card>
  );
}

function PlanAction({ title, description, children }: { title: string; description: string; children: React.ReactNode }) {
  return <div className="border-t pt-4"><div className="flex items-center justify-between"><div><p className="text-sm font-medium">{title}</p><p className="text-sm text-muted-foreground">{description}</p></div>{children}</div></div>;
}

function formatNextResetDate(dateStr?: string, locale = "en-US") {
  if (!dateStr) return null;
  return new Date(dateStr).toLocaleDateString(locale, { month: "short", day: "numeric", year: "numeric" });
}
