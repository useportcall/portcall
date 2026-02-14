import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import { useDeleteEntitlement, useUpdateEntitlement } from "@/hooks/api/entitlements";
import { cn } from "@/lib/utils";
import { Entitlement } from "@/models/entitlement";
import { Gauge, Loader2, MoreHorizontal, RotateCcw, Trash2 } from "lucide-react";
import { useState } from "react";

export function MeteredFeatureCard({ userId, entitlement }: { userId: string; entitlement: Entitlement }) {
  const [quotaValue, setQuotaValue] = useState(entitlement.quota.toString());
  const [usageValue, setUsageValue] = useState(entitlement.usage.toString());
  const { mutate: updateEntitlement, isPending: isUpdating } = useUpdateEntitlement({ userId, featureId: entitlement.id });
  const { mutate: deleteEntitlement, isPending: isDeleting } = useDeleteEntitlement({ userId, featureId: entitlement.id });
  const usagePercentage = entitlement.quota > 0 ? Math.min(100, (entitlement.usage / entitlement.quota) * 100) : 0;
  const isOverLimit = entitlement.quota > 0 && entitlement.usage >= entitlement.quota;

  const saveValue = (raw: string, current: number, field: "quota" | "usage") => {
    const next = Number.parseInt(raw, 10);
    if (!Number.isNaN(next) && next !== current) updateEntitlement({ [field]: next });
  };

  return (
    <Card className={cn("transition-colors", isDeleting && "opacity-50")}><CardContent className="p-4">
      <div className="flex items-start justify-between mb-4"><div className="flex items-center gap-2"><Gauge className="size-4 text-muted-foreground" /><span className="font-mono text-sm font-medium">{entitlement.id}</span></div><DropdownMenu><DropdownMenuTrigger asChild><Button variant="ghost" size="icon" className="size-7"><MoreHorizontal className="size-4" /></Button></DropdownMenuTrigger><DropdownMenuContent align="end"><DropdownMenuItem onClick={() => { updateEntitlement({ usage: 0 }); setUsageValue("0"); }} disabled={isUpdating}><RotateCcw className="size-4 mr-2" />Reset Usage</DropdownMenuItem><DropdownMenuItem onClick={() => deleteEntitlement({})} disabled={isDeleting} className="text-destructive focus:text-destructive"><Trash2 className="size-4 mr-2" />Remove</DropdownMenuItem></DropdownMenuContent></DropdownMenu></div>
      <div className="mb-4"><div className="flex items-center justify-between text-xs mb-1.5"><span className="text-muted-foreground">Usage</span><span className={cn("font-medium", isOverLimit && "text-destructive")}>{entitlement.usage.toLocaleString()} / {entitlement.quota === -1 ? "∞" : entitlement.quota.toLocaleString()}</span></div><div className="h-2 bg-muted rounded-full overflow-hidden"><div className={cn("h-full transition-all rounded-full", isOverLimit ? "bg-destructive" : "bg-primary")} style={{ width: `${Math.min(100, usagePercentage)}%` }} /></div></div>
      <div className="grid grid-cols-2 gap-4"><InputField label="Quota Limit" value={quotaValue} onChange={setQuotaValue} onSave={() => saveValue(quotaValue, entitlement.quota, "quota")} placeholder="-1 for unlimited" /><InputField label="Current Usage" value={usageValue} onChange={setUsageValue} onSave={() => saveValue(usageValue, entitlement.usage, "usage")} /></div>
      {entitlement.next_reset_at && <p className="text-xs text-muted-foreground mt-3">Auto-resets: {new Date(entitlement.next_reset_at).toLocaleDateString()}</p>}
      {isUpdating && <div className="flex items-center gap-1.5 mt-3 text-xs text-muted-foreground"><Loader2 className="size-3 animate-spin" />Saving...</div>}
    </CardContent></Card>
  );
}

function InputField(props: { label: string; value: string; onChange: (value: string) => void; onSave: () => void; placeholder?: string }) {
  return (
    <div className="space-y-1.5">
      <label className="text-xs text-muted-foreground">{props.label}</label>
      <Input type="number" value={props.value} onChange={(e) => props.onChange(e.target.value)} onBlur={props.onSave} onKeyDown={(e) => e.key === "Enter" && props.onSave()} className="h-8 text-sm" placeholder={props.placeholder} />
    </div>
  );
}
