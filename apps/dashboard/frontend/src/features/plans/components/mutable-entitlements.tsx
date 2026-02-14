import { useListBasicPlanFeatures } from "@/hooks";
import { AddFeatureButton } from "./mutable-entitlements-add-feature";
import { LinkedFeatureBadge } from "./mutable-entitlements-linked-feature";

export default function MutableEntitlements({ id }: { id: string }) {
  const { data: planFeatures } = useListBasicPlanFeatures(id);
  if (!planFeatures?.data) return null;

  return (
    <div className="space-y-2 px-2" data-testid="included-features">
      <div className="flex items-center justify-between">
        <h2 className="text-sm text-muted-foreground">Included features</h2>
        <AddFeatureButton planFeatures={planFeatures.data} planId={id} />
      </div>
      <div className="flex flex-wrap gap-2" data-testid="feature-badges">
        {planFeatures.data.map((pf) => (
          <LinkedFeatureBadge key={pf.id} planFeature={pf} planId={id} />
        ))}
      </div>
    </div>
  );
}
