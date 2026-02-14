import { Button } from "@/components/ui/button";
import { Tier } from "@/models/plan-item";
import { MinusCircle } from "lucide-react";

export function TierInput({ value, onChange }: { value: Tier[]; onChange: (arg0: Tier[]) => void }) {
  return (
    <div><button onClick={() => onChange(normalizeTiers([...value, { unit_amount: 0, start: 0, end: 0 }]))}>Add</button><table><thead><tr><th className="text-left text-sm font-medium">Charge</th><th className="text-left text-sm font-medium">From</th><th className="text-left text-sm font-medium">To</th><th /></tr></thead><tbody>{value.map((tier, index) => <tr key={index}><td className="text-left text-sm font-medium flex"><span>$</span><input type="text" className="outline-none w-full text-sm" defaultValue={(tier.unit_amount / 100).toFixed(2)} onChange={(e) => { if (!isNaN(Number(e.target.value))) { const next = [...value]; next[index].unit_amount = Number(e.target.value) * 100; onChange(next); } }} /></td><td className="text-left text-sm font-medium"><input type="text" className="outline-none w-full text-sm" value={tier.start} onChange={(e) => { const next = [...value]; next[index].start = Number(e.target.value); onChange(next); }} onBlur={() => onChange(normalizeTiers(value))} /></td><td className="text-left text-sm font-medium"><input type="text" className="outline-none w-full text-sm" value={tier.end < 0 ? "∞" : tier.end} onChange={(e) => { const next = [...value]; next[index].end = Number(e.target.value); onChange(next); }} onBlur={() => onChange(normalizeTiers(value))} /></td><td><Button size="icon" variant="ghost" className="rounded-md hover:bg-accent size-6" disabled={value.length === 1} onClick={() => { const next = [...value]; next.splice(index, 1); onChange(normalizeTiers(next)); }}><MinusCircle className="size-3 text-red-500" /></Button></td></tr>)}</tbody></table></div>
  );
}

function normalizeTiers(tiers: Tier[]) {
  if (!tiers.length) return tiers;
  if (tiers.length === 1) { tiers[0].start = 0; tiers[0].end = -1; return tiers; }
  for (let i = 0; i < tiers.length; i++) {
    if (i === 0) { if (tiers[i].end === -1) tiers[i].end = 1000; tiers[i].start = 0; continue; }
    if (i === tiers.length - 1) { tiers[i].start = tiers[i - 1].end + 1; tiers[i].end = -1; continue; }
    if (tiers[i].end === -1) tiers[i].end = tiers[i].start + 1000;
    tiers[i].start = tiers[i - 1].end + 1;
  }
  return tiers;
}
