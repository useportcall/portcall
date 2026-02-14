import { Separator } from "@/components/ui/separator";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { useDeletePlanItem, useUpdatePlanItem } from "@/hooks";
import { MeteredPlanItem } from "@/models/plan-item";
import { MoreVertical } from "lucide-react";
import { Suspense, useState } from "react";
import MutableMeteredLimit from "./mutable-metered-limit";
import MutableMeteredRollover from "./mutable-metered-rollover";
import { PricingModelSelect } from "./mutable-price";
import { MutableMeteredFeatureInterval } from "./mutable-reset";
import { PlanFeatureComboBox } from "./plan-feature-combo-box";
import { MutablePlanItemPrice } from "./plan-item-price-table";

export function PlanItemCard({ planItem, planId, pricingDisabled = false }: { planItem: MeteredPlanItem; planId: string; pricingDisabled?: boolean }) {
  const updatePlanItem = useUpdatePlanItem(planItem.id, planId);
  const deletePlanItem = useDeletePlanItem(planItem.id, planId);
  if (!planItem?.features?.length) return null;

  return (
    <Suspense><Card data-testid="metered-feature-card"><CardHeader className="flex flex-col items-start gap-0"><div className="flex w-full justify-between items-start"><div className="self-start flex-col justify-evenly gap-4 max-w-[200px]"><input data-testid="metered-title-input" aria-label="Metered feature title" type="text" className="font-semibold w-full outline-none" placeholder="API Calls" defaultValue={planItem.public_title} onBlur={(e) => e.target.value !== planItem.public_title && updatePlanItem.mutateAsync({ public_title: e.target.value })} /><MutableMultiline className="text-xs text-muted-foreground w-full resize-none outline-none" placeholder="No description" value={planItem.public_description} saveFn={(value) => updatePlanItem.mutateAsync({ public_description: value })} /></div><DropdownMenu><DropdownMenuTrigger asChild><button type="button" aria-label="Metered feature actions"><MoreVertical className="w-6 h-6 p-1 hover:bg-accent rounded-sm" /></button></DropdownMenuTrigger><DropdownMenuContent><DropdownMenuItem onClick={() => deletePlanItem.mutateAsync({})}>Delete</DropdownMenuItem></DropdownMenuContent></DropdownMenu></div><Separator /><div className="flex flex-wrap lg:flex-nowrap gap-6 mt-4"><Field title="Resets"><MutableMeteredFeatureInterval planFeature={planItem.feature} planId={planId} /></Field><Field title="Metered feature"><PlanFeatureComboBox planFeature={planItem.feature} /></Field><Field title="Pricing"><PricingModelSelect value={planItem.pricing_model} onChange={(pricing_model) => updatePlanItem.mutateAsync({ pricing_model })} disabled={pricingDisabled} /></Field><Field title="Rollover"><MutableMeteredRollover planFeature={planItem.feature} planId={planId} /></Field><Field title="Limit"><MutableMeteredLimit planFeature={planItem.feature} planId={planId} /></Field></div></CardHeader>{planItem.pricing_model !== "none" && <CardContent><MutablePlanItemPrice planItem={planItem} planId={planId} /></CardContent>}</Card></Suspense>
  );
}

function Field({ title, children }: { title: string; children: React.ReactNode }) { return <div className="flex flex-col gap-1 w-fit justify-start"><p className="text-xs text-start text-muted-foreground">{title}</p>{children}</div>; }
function MutableMultiline({ value, placeholder, saveFn, className }: { value: string; placeholder?: string; saveFn: (value: string) => Promise<any>; className?: string; }) { const [state, setState] = useState(value); return <textarea aria-label="Metered feature description" className={className} placeholder={placeholder} value={state} onChange={(e) => setState(e.target.value)} onBlur={(e) => e.target.value !== value && saveFn(e.target.value)} />; }
