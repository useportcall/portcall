import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { ArrowRight, CalendarClock } from "lucide-react";
import { useTranslation } from "react-i18next";

interface SubscriptionData { next_reset_at?: string; scheduled_plan?: { name: string }; current_plan?: { name: string }; plan_name?: string; }

export function PendingChangesDialog({ open, onOpenChange, subscription, hasPendingChange, formatNextResetDate }: { open: boolean; onOpenChange: (open: boolean) => void; subscription?: SubscriptionData; hasPendingChange: boolean; formatNextResetDate: (dateStr?: string) => string | null; }) {
  const { t } = useTranslation();
  if (!subscription?.next_reset_at) return null;

  return (
    <div className="border-t pt-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2"><CalendarClock className="h-4 w-4 text-muted-foreground" /><span className="text-sm text-muted-foreground">{t("views.usage.pending.next_cycle")}</span><span className="text-sm font-medium">{formatNextResetDate(subscription.next_reset_at)}</span></div>
        {hasPendingChange && <Dialog open={open} onOpenChange={onOpenChange}><DialogTrigger asChild><Badge variant="secondary" className="cursor-pointer bg-amber-100 text-amber-800 hover:bg-amber-200"><CalendarClock className="h-3 w-3 mr-1" />{t("views.usage.pending.badge")}</Badge></DialogTrigger><DialogContent className="sm:max-w-lg"><DialogHeader><DialogTitle>{t("views.usage.pending.dialog_title")}</DialogTitle><DialogDescription>{t("views.usage.pending.dialog_description")}</DialogDescription></DialogHeader><div className="py-4 space-y-4"><PlanChangePreview subscription={subscription} formatNextResetDate={formatNextResetDate} />{subscription?.scheduled_plan?.name === "Free" && <DowngradeImpactList />}</div><DialogFooter><Button variant="outline" onClick={() => onOpenChange(false)}>{t("common.close")}</Button></DialogFooter></DialogContent></Dialog>}
      </div>
    </div>
  );
}

function PlanChangePreview({ subscription, formatNextResetDate }: { subscription: SubscriptionData; formatNextResetDate: (dateStr?: string) => string | null; }) {
  const { t } = useTranslation();
  return <><div className="flex items-center justify-center gap-4"><div className="text-center p-4 border rounded-lg min-w-[120px]"><p className="text-xs text-muted-foreground mb-1">{t("views.usage.pending.current_plan")}</p><p className="font-semibold">{subscription.current_plan?.name ?? subscription.plan_name}</p></div><ArrowRight className="h-5 w-5 text-muted-foreground" /><div className="text-center p-4 border rounded-lg min-w-[120px] bg-amber-50 border-amber-200"><p className="text-xs text-muted-foreground mb-1">{t("views.usage.pending.scheduled_plan")}</p><p className="font-semibold text-amber-800">{subscription.scheduled_plan?.name}</p></div></div><div className="text-center text-sm text-muted-foreground"><CalendarClock className="h-4 w-4 inline mr-1" />{t("views.usage.pending.effective")}: <span className="font-medium">{formatNextResetDate(subscription.next_reset_at)}</span></div></>;
}

function DowngradeImpactList() {
  const { t } = useTranslation();
  return <div className="p-3 bg-muted rounded-lg text-sm"><p className="font-medium mb-1">{t("views.usage.pending.impact_title")}</p><ul className="list-disc list-inside text-muted-foreground space-y-1"><li>{t("views.usage.downgrade_dialog.after_users")}</li><li>{t("views.usage.downgrade_dialog.after_subscriptions")}</li><li>{t("views.usage.downgrade_dialog.after_data")}</li></ul></div>;
}
