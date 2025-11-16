import { Button } from "@/components/ui/button";
import { useCreatePlanItem, useListMeteredPlanItems } from "@/hooks";
import { Plus } from "lucide-react";
import { useParams } from "react-router-dom";
import { PlanItemCard } from "./metered-feature-card";

export function PlanItems() {
  const { data: items } = useListMeteredPlanItems();

  return (
    <div className="h-full space-y-2 px-4 animate-fade-in">
      <div className="flex justify-between">
        <h2 className="text-sm text-muted-foreground px-2">Metered features</h2>
        <div className="pr-1">
          <CreateMeteredFeatureButton />
        </div>
      </div>
      <div className="py-2 md:px-2 space-y-2 overflow-y-auto">
        {items.map((planItem) => (
          <PlanItemCard key={planItem.id} planItem={planItem} />
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

function CreateMeteredFeatureButton() {
  const { id } = useParams();
  const { mutateAsync } = useCreatePlanItem(id!);

  if (!id) return null;

  return (
    <Button
      size="icon"
      variant="ghost"
      className="size-[22px] rounded-md hover:bg-slate-100"
      onClick={async () => {
        await mutateAsync({
          public_title: "New Metered Feature",
          public_description: "",
          unit_amount: 0,
          pricing_model: "none",
          feature_id: null,
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
