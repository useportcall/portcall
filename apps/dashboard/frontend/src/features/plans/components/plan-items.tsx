import { Button } from "@/components/ui/button";
import { useCreateFeature, useCreatePlanItem, useListMeteredPlanItems } from "@/hooks";
import { Plus } from "lucide-react";
import { PlanItemCard } from "./metered-feature-card";

export function PlanItems(props: { id: string; pricingDisabled?: boolean }) {
  const { data: items } = useListMeteredPlanItems({ id: props.id });

  return (
    <div
      className="h-full space-y-2 px-4 animate-fade-in"
      data-testid="metered-features"
    >
      <div className="flex justify-between">
        <h2 className="text-sm text-muted-foreground px-2">Metered features</h2>
        <div className="pr-1">
          <CreateMeteredFeatureButton id={props.id} />
        </div>
      </div>
      <div className="py-2 md:px-2 space-y-2 overflow-y-auto">
        {items.map((planItem) => (
          <PlanItemCard
            key={planItem.id}
            planItem={planItem}
            planId={props.id}
            pricingDisabled={props.pricingDisabled}
          />
        ))}
        {items.length === 0 && (
          <span className="flex justify-center w-full text-sm text-muted-foreground">
            No metered features added yet.
          </span>
        )}
      </div>
    </div>
  );
}

function CreateMeteredFeatureButton({ id }: { id: string }) {
  const { mutateAsync: createFeature } = useCreateFeature({ isMetered: true });
  const { mutateAsync: createPlanItem } = useCreatePlanItem(id);

  return (
    <Button
      type="button"
      data-testid="add-metered-feature-button"
      aria-label="Add metered feature"
      size="icon"
      variant="ghost"
      className="size-[22px] rounded-md hover:bg-accent"
      onClick={async () => {
        const featureId = `metered_feature_${Date.now()}`;
        await createFeature({ feature_id: featureId, is_metered: true });
        await createPlanItem({
          public_title: "New Metered Feature",
          public_description: "",
          unit_amount: 0,
          pricing_model: "none",
          feature_id: featureId,
          quota: 0,
          tiers: [],
          interval: "month",
          plan_id: id,
        });
      }}
    >
      <Plus className="" />
    </Button>
  );
}
