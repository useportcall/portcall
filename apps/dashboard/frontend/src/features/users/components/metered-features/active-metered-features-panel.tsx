import { Loader2, Activity } from "lucide-react";
import { Entitlement } from "@/models/entitlement";
import { MeteredFeatureCard } from "./metered-feature-card";

export function ActiveMeteredFeaturesPanel({ isLoading, entitlements, userId }: { isLoading: boolean; entitlements: Entitlement[]; userId: string }) {
  return (
    <div className="flex-1 flex flex-col gap-3 min-w-0">
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-medium text-foreground">Active Entitlements</h3>
        <span className="text-xs text-muted-foreground">{entitlements.length} configured</span>
      </div>
      <div className="flex-1 overflow-auto">
        {isLoading ? (
          <div className="flex items-center justify-center h-32"><Loader2 className="size-6 animate-spin text-muted-foreground" /></div>
        ) : entitlements.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-32 text-muted-foreground border-2 border-dashed rounded-lg"><Activity className="size-8 mb-2 opacity-50" /><p className="text-sm">No metered features assigned</p><p className="text-xs mt-1">Select a feature from the left to add it</p></div>
        ) : (
          <div className="grid gap-3">{entitlements.map((entitlement) => <MeteredFeatureCard key={entitlement.id} userId={userId} entitlement={entitlement} />)}</div>
        )}
      </div>
    </div>
  );
}
