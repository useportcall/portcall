import { cn } from "@/lib/utils";
import { Feature } from "@/models/feature";
import { Check, Loader2 } from "lucide-react";

export function FeatureToggleRow({
  feature,
  isEnabled,
  isPending,
  onToggle,
}: {
  feature: Feature;
  isEnabled: boolean;
  isPending: boolean;
  onToggle: () => void;
}) {
  return (
    <div className={cn("flex items-center justify-between px-4 py-3 transition-colors cursor-pointer hover:bg-accent/50", isEnabled && "bg-primary/5 dark:bg-primary/10")} onClick={onToggle}>
      <div className="flex items-center gap-3">
        <div className={cn("flex items-center justify-center size-5 rounded border transition-colors", isEnabled ? "bg-primary border-primary text-primary-foreground" : "border-input bg-background")}>
          {isPending ? <Loader2 className="size-3 animate-spin" /> : isEnabled ? <Check className="size-3" /> : null}
        </div>
        <span className="font-mono text-sm">{feature.id}</span>
      </div>
      <span className={cn("text-xs px-2 py-0.5 rounded-full", isEnabled ? "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400" : "bg-muted text-muted-foreground")}>
        {isEnabled ? "Enabled" : "Disabled"}
      </span>
    </div>
  );
}
