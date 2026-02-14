import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useCopyPlan, useListApps } from "@/hooks";
import { useApp } from "@/hooks/use-app";
import { Plan } from "@/models/plan";
import { Loader2 } from "lucide-react";
import { useState } from "react";

export function CopyPlanDialog({ plan, open, onOpenChange }: { plan: Plan; open: boolean; onOpenChange: (open: boolean) => void }) {
  const currentApp = useApp();
  const { data: apps } = useListApps();
  const [targetAppId, setTargetAppId] = useState("");
  const { mutateAsync, isPending } = useCopyPlan(plan.id);

  return (
    <Dialog open={open} onOpenChange={onOpenChange}><DialogContent className="flex flex-col gap-4"><DialogHeader><DialogTitle>Copy plan to another app</DialogTitle><DialogDescription>Copy "{plan.name}" and all its items and features to another app.</DialogDescription></DialogHeader><form onSubmit={async (e) => { e.preventDefault(); e.stopPropagation(); if (!targetAppId) return; await mutateAsync({ target_app_id: targetAppId }); setTargetAppId(""); onOpenChange(false); }} className="w-full flex flex-col gap-4"><div className="flex flex-col gap-2"><label htmlFor="target-app" className="text-sm font-medium">Target app</label><Select value={targetAppId} onValueChange={setTargetAppId}><SelectTrigger id="target-app" onClick={(e) => e.stopPropagation()}><SelectValue placeholder="Select an app" /></SelectTrigger><SelectContent>{apps?.data.filter((app) => app.id !== currentApp.id).map((app) => <SelectItem key={app.id} value={app.id}>{app.name} {app.is_live ? "(Live)" : "(Test)"}</SelectItem>)}</SelectContent></Select></div><div className="w-fit mt-4 flex flex-row gap-2 justify-end self-end"><Button type="button" variant="ghost" onClick={(e) => { e.stopPropagation(); onOpenChange(false); }} disabled={isPending}>Cancel</Button><Button type="submit" size="sm" disabled={isPending || !targetAppId} onClick={(e) => e.stopPropagation()}>Copy plan{isPending && <Loader2 className="size-4 animate-spin" />}</Button></div></form></DialogContent></Dialog>
  );
}
