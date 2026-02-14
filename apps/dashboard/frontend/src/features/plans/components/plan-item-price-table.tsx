import { Table, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { useUpdatePlanItem } from "@/hooks";
import { PlanItem, Tier } from "@/models/plan-item";
import { Plus, Trash2 } from "lucide-react";
import { useEffect, useState } from "react";
import { formatPriceInput, parsePriceToCents } from "./plan-item-price-utils";

export function MutablePlanItemPrice({ planItem, planId }: { planItem: PlanItem; planId: string }) {
  const { mutateAsync } = useUpdatePlanItem(planItem.id, planId);
  return <Table className="text-sm"><TableHeader className="border-b"><TableHead className="text-sm text-foreground font-medium border-r">From</TableHead><TableHead className="text-sm text-foreground border-r">To</TableHead><TableHead className="text-sm text-foreground">Price</TableHead><TableHead className="text-sm text-foreground flex justify-end"><button type="button" aria-label="Add pricing tier" disabled={planItem.pricing_model === "unit"} className="p-0 disabled:text-accent-foreground/50" onClick={async () => { const tiers = [...(planItem.tiers || [])]; if (tiers.length > 0 && tiers[tiers.length - 1].end === -1) tiers[tiers.length - 1].end = tiers[tiers.length - 1].start + 1000; tiers.push({ start: tiers.length > 0 ? tiers[tiers.length - 1].end + 1 : 0, end: -1, unit_amount: 0 }); await mutateAsync({ tiers }); }}><Plus className="size-4" /></button></TableHead></TableHeader><MutableAmount planItem={planItem} planId={planId} /></Table>;
}

function MutableAmount({ planItem, planId }: { planItem: PlanItem; planId: string }) {
  if (["unit", "fixed"].includes(planItem.pricing_model)) return <TierRow index={0} tier={{ start: 0, end: -1, unit_amount: planItem.unit_amount }} planItem={planItem} planId={planId} />;
  return planItem.tiers?.map((tier, index) => <TierRow tier={tier} key={index} planItem={planItem} index={index} planId={planId} />) || null;
}

function TierRow({ planItem, tier, index, planId }: { planItem: PlanItem; tier: Tier; index: number; planId: string }) {
  const { mutateAsync } = useUpdatePlanItem(planItem.id, planId);
  const [endValue, setEndValue] = useState(tier.end); useEffect(() => setEndValue(tier.end), [tier.end]);
  const [value, setValue] = useState(formatPriceInput((tier.unit_amount / 100).toString()));
  const isLast = index === (planItem.tiers?.length || 1) - 1;

  return <TableRow className="border-b-0"><TableCell className="text-sm border-r"><input type="text" value={tier.start} readOnly className="w-12 outline-none" onBlur={(e) => { let start = Number(e.target.value); if (index === 0) start = 0; const tiers = planItem.tiers?.map((t, i) => i === index ? { ...t, start } : t); mutateAsync({ tiers }); }} /></TableCell><TableCell className="text-sm border-r"><input type="text" value={tier.end === -1 ? "∞" : endValue} readOnly={tier.end === -1} className="w-12 outline-none" onChange={(e) => { const next = parseInt(e.target.value); if (!isNaN(next)) setEndValue(next); }} onBlur={(e) => { const end = Number(e.target.value); const tiers = planItem.tiers?.map((t, i) => i === index ? (isLast ? { ...t, end: -1 } : { ...t, end }) : i === index + 1 ? { ...t, start: end + 1 } : t); mutateAsync({ tiers }); }} /></TableCell><TableCell className="text-sm">$<input data-testid="metered-price-input" aria-label="Metered price input" className="w-16 outline-none" type="text" value={value} onChange={(e) => !isNaN(Number(e.target.value)) && setValue(e.target.value)} onBlur={async (e) => { const unit_amount = parsePriceToCents(e.target.value); setValue(formatPriceInput(e.target.value)); const tiers = planItem.tiers?.map((t, i) => i === index ? { ...t, unit_amount } : t) ?? []; await mutateAsync({ unit_amount, tiers }); }} /></TableCell><TableCell className="flex justify-end"><button type="button" aria-label="Remove pricing tier" className="p-0 disabled:text-accent-foreground/50 text-foreground hover:text-red-400" disabled={index === 0} onClick={async () => { const tiers = planItem.tiers?.filter((_, i) => i !== index) ?? []; if (tiers.length > 0 && index === tiers.length) tiers[tiers.length - 1].end = -1; await mutateAsync({ tiers }); }}><Trash2 className="w-4 h-4 cursor-pointer" /></button></TableCell></TableRow>;
}
