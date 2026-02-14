import { cn } from "@/lib/utils";
import { Plan } from "@/models/plan";
import { Check } from "lucide-react";

export function PlanSelectRow({ plan, isSelected, onSelect }: { plan: Plan; isSelected: boolean; onSelect: () => void }) {
  return (
    <div className={cn("flex items-center justify-between px-4 py-3 transition-colors cursor-pointer hover:bg-accent/50", isSelected && "bg-primary/5 dark:bg-primary/10")} onClick={onSelect}>
      <div className="flex items-center gap-3"><div className={cn("flex items-center justify-center size-5 rounded-full border transition-colors", isSelected ? "bg-primary border-primary text-primary-foreground" : "border-input bg-background")}>{isSelected && <Check className="size-3" />}</div><div className="flex flex-col"><span className="text-sm font-medium">{plan.name}</span><span className="text-xs text-muted-foreground">{plan.interval_count > 1 ? `${plan.interval_count} ${plan.interval}s` : plan.interval} · {plan.currency.toUpperCase()}</span></div></div>
      <span className={cn("text-xs px-2 py-0.5 rounded-full", plan.status === "published" ? "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400" : "bg-muted text-muted-foreground")}>{plan.status}</span>
    </div>
  );
}
