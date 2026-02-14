import { Button } from "@/components/ui/button";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Select, SelectContent, SelectGroup, SelectItem } from "@/components/ui/select";
import { SelectTrigger } from "@radix-ui/react-select";
import { ChevronDown, Coins } from "lucide-react";
import { TierInput } from "./tier-input";
import { currencyCodetoSymbol } from "../utils";

const DEFAULT_TIER_1 = { unit_amount: 200, start: 0, end: 1000 };
const DEFAULT_TIER_2 = { unit_amount: 100, start: 1001, end: -1 };

export function PriceCreate({ form, currency }: { form: any; currency: string }) {
  return <Popover><PopoverTrigger asChild><Button variant="outline" className="px-2 py-2 h-6 text-xs font-semibold"><Coins className="size-3" /><CurrentAmount currency={currency} pricingModel={form.watch("pricing_model")} amount={form.watch("unit_amount")} /></Button></PopoverTrigger><PopoverContent side="left" align="start" className="px-2 py-4 space-y-4"><div className="flex justify-between px-2"><PricingModelSelect value={form.watch("pricing_model")} onChange={(value) => { if (value === "unit") { form.setValue("tiers", []); form.setValue("unit_amount", 0); } else { form.setValue("tiers", [DEFAULT_TIER_1, DEFAULT_TIER_2]); form.setValue("unit_amount", DEFAULT_TIER_1.unit_amount); } form.setValue("pricing_model", value); }} /></div><PriceInput form={form} /></PopoverContent></Popover>;
}

function CurrentAmount({ pricingModel, amount, currency }: { pricingModel: string; amount: any; currency: string }) {
  const symbol = currencyCodetoSymbol(currency); return <span>{symbol}{(amount / 100).toFixed(2)}{pricingModel === "unit" ? "" : "*"}</span>;
}
function PriceInput({ form }: { form: any }) { return form.watch("pricing_model") === "unit" ? <UnitAmountInput value={(form.watch("unit_amount") / 100).toFixed(2)} onChange={(value) => form.setValue("unit_amount", value * 100)} /> : <TierInput value={form.watch("tiers")} onChange={(value) => form.setValue("tiers", value)} />; }
function UnitAmountInput({ value, onChange }: { value: string; onChange: (value: number) => void }) { return <input type="text" className="outline-none w-full text-sm" placeholder="0.00" value={value} onChange={(e) => !isNaN(Number(e.target.value)) && onChange(Number(e.target.value))} />; }

export function PricingModelSelect({ value, onChange, disabled }: { value: string; onChange: (value: string) => void; disabled?: boolean }) {
  return <Select value={value} onValueChange={(v) => onChange(v)}><SelectTrigger data-testid="metered-pricing-model-button" aria-label="Metered pricing model" disabled={disabled} className="cursor-pointer text-sm w-16 overflow-hidden text-ellipsis whitespace-nowrap rounded py-1 flex items-center justify-between disabled:text-muted-foreground disabled:cursor-not-allowed">{value}<ChevronDown className="size-3" /></SelectTrigger><SelectContent><SelectGroup><SelectItem value="none">None</SelectItem><SelectItem value="unit">Unit</SelectItem><SelectItem value="tier">Tier</SelectItem><SelectItem value="block">Block</SelectItem></SelectGroup></SelectContent></Select>;
}
