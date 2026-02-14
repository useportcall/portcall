import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { useCreateFeature, useListFeatures, useUpdatePlanFeature } from "@/hooks";
import { PlanFeature } from "@/models/plan-feature";
import { ChevronDown } from "lucide-react";
import React, { useState } from "react";
import { useParams } from "react-router-dom";
import { toSnakeCase } from "../utils";

export function PlanFeatureComboBox({ planFeature }: { planFeature?: PlanFeature }) {
  const [open, setOpen] = React.useState(false);
  const [value, setValue] = useState("");
  const { data: features } = useListFeatures({ isMetered: true });
  const { mutateAsync: addPlanFeature } = useCreateFeature({ isMetered: true });
  const { id } = useParams();
  if (!features?.data) return null;

  return (
    <Popover open={open} onOpenChange={setOpen}><PopoverTrigger asChild><button type="button" data-testid="metered-feature-button" aria-label="Metered feature selector" className="cursor-pointer text-sm w-20 rounded py-1 flex items-center justify-between"><span className="overflow-hidden text-ellipsis whitespace-nowrap">{planFeature?.feature.id || "Select feature"}</span><ChevronDown className="size-3 shrink-0" /></button></PopoverTrigger><PopoverContent className="w-[200px] p-0" align="start"><Command><CommandInput data-testid="metered-feature-input" value={value} onValueChange={(v) => setValue(toSnakeCase(v))} placeholder="Search features..." onKeyDown={async (e) => { if (e.key === "Enter" && !features.data.find((f) => f.id === e.currentTarget.value)) { await addPlanFeature({ plan_id: id, feature_id: toSnakeCase(value), is_metered: true, plan_feature_id: planFeature?.id }); setOpen(false); } }} /><CommandList><CommandEmpty>Press "Enter" to add</CommandEmpty><CommandGroup>{planFeature && features.data.map((f) => <PlanFeatureCommandItem key={f.id} planId={id!} planFeatureId={planFeature.id} featureId={f.id} />)}</CommandGroup></CommandList></Command></PopoverContent></Popover>
  );
}

function PlanFeatureCommandItem({ planId, planFeatureId, featureId }: { planId: string; planFeatureId: string; featureId: string }) {
  const { mutate } = useUpdatePlanFeature(planFeatureId, planId);
  return <CommandItem onSelect={() => mutate({ feature_id: featureId })}>{featureId}</CommandItem>;
}
