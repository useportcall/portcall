import { Badge } from "@/components/ui/badge";
import { useDeletePlanFeature } from "@/hooks";
import { PlanFeature } from "@/models/plan-feature";
import { cn } from "@/lib/utils";
import { X } from "lucide-react";

export function LinkedFeatureBadge({
  planFeature,
  planId,
}: {
  planFeature: PlanFeature;
  planId: string;
}) {
  const { mutate: removeFeature } = useDeletePlanFeature(planFeature, planId);

  return (
    <Badge variant="outline" className={cn("gap-1 cursor-default")}>
      <h3>{planFeature.feature.id}</h3>
      <button
        type="button"
        className="rounded-sm hover:text-red-800"
        aria-label={`Remove feature ${planFeature.feature.id}`}
        onClick={() => removeFeature({})}
      >
        <X className="w-3 h-3" />
      </button>
    </Badge>
  );
}
