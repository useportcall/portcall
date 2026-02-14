import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { useListPlans } from "@/hooks";
import { useCreateSubscription } from "@/hooks/api/subscriptions";
import { useQueryClient } from "@tanstack/react-query";
import { Loader2, Search, UserPlus } from "lucide-react";
import { useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { PlanSelectRow } from "./subscribe-user/plan-select-row";

export function SubscribeUserDialog({ userId }: { userId: string }) {
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);
  const [searchValue, setSearchValue] = useState("");
  const [selectedPlanId, setSelectedPlanId] = useState<string | null>(null);
  const queryClient = useQueryClient();
  const { data: plans, isLoading: plansLoading } = useListPlans("*");
  const { mutateAsync: createSubscription, isPending } = useCreateSubscription();

  const filteredPlans = useMemo(
    () =>
      (plans?.data || [])
        .filter((plan) => plan.status === "published")
        .filter((plan) => !searchValue || plan.name.toLowerCase().includes(searchValue.toLowerCase())),
    [plans?.data, searchValue],
  );

  return (
    <Dialog open={open} onOpenChange={(isOpen) => { setOpen(isOpen); if (!isOpen) { setSelectedPlanId(null); setSearchValue(""); } }}>
      <DialogTrigger asChild><Button variant="outline" className="gap-2"><UserPlus className="size-4" />{t("views.user.subscription.subscribe_action")}</Button></DialogTrigger>
      <DialogContent className="max-w-lg max-h-[85vh] overflow-hidden flex flex-col">
        <DialogHeader><DialogTitle>{t("views.user.subscription.subscribe_title")}</DialogTitle><DialogDescription>{t("views.user.subscription.subscribe_description")}</DialogDescription></DialogHeader>
        <div className="flex-1 flex flex-col overflow-hidden"><div className="relative mt-2 mb-3"><Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" /><input type="text" value={searchValue} onChange={(e) => setSearchValue(e.target.value)} placeholder={t("views.user.subscription.search_plans")} className="w-full h-10 pl-10 pr-4 rounded-lg border border-input bg-transparent text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2" /></div>
          <div className="flex-1 overflow-auto rounded-lg border bg-card max-h-[300px]">
            {plansLoading ? <div className="flex items-center justify-center h-32"><Loader2 className="size-6 animate-spin text-muted-foreground" /></div> : filteredPlans.length === 0 ? <div className="flex flex-col items-center justify-center h-32 text-muted-foreground"><p className="text-sm">{t("views.user.subscription.no_plans_found")}</p><p className="text-xs mt-1">{t("views.user.subscription.no_plans_hint")}</p></div> : <div className="divide-y divide-border">{filteredPlans.map((plan) => <PlanSelectRow key={plan.id} plan={plan} isSelected={selectedPlanId === plan.id} onSelect={() => setSelectedPlanId(plan.id)} />)}</div>}
          </div>
        </div>
        <DialogFooter className="mt-4"><Button variant="outline" onClick={() => setOpen(false)}>{t("common.cancel")}</Button><Button disabled={!selectedPlanId || isPending} onClick={async () => { if (!selectedPlanId) return; await createSubscription({ user_id: userId, plan_id: selectedPlanId }); queryClient.invalidateQueries({ queryKey: [`/users/${userId}/subscription`] }); setOpen(false); setSelectedPlanId(null); setSearchValue(""); }}>{isPending && <Loader2 className="size-4 animate-spin mr-2" />}{t("views.user.subscription.subscribe_confirm")}</Button></DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
