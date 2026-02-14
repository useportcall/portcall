import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { BillingAddress, useBillingAddress, useUpsertBillingAddress } from "@/hooks/api/quota";
import { Loader2, MapPin, Pencil } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";

type BillingAddressFormData = { line1: string; line2: string; city: string; state: string; postal_code: string; country: string };

export function BillingAddressCard() {
  const { t } = useTranslation();
  const { data: addressData } = useBillingAddress();
  const address = addressData?.data?.billing_address || null;
  return (
    <Card><CardHeader className="pb-3"><div className="flex items-center justify-between"><div className="flex items-center gap-2"><MapPin className="h-5 w-5 text-primary" /><CardTitle>{t("views.usage.address.title")}</CardTitle></div><EditBillingAddressDialog address={address} /></div><CardDescription>{t("views.usage.address.description")}</CardDescription></CardHeader><CardContent><AddressDisplay address={address} /></CardContent></Card>
  );
}

function EditBillingAddressDialog({ address }: { address: BillingAddress | null }) {
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);
  const [formData, setFormData] = useState<BillingAddressFormData>({ line1: address?.line1 || "", line2: address?.line2 || "", city: address?.city || "", state: address?.state || "", postal_code: address?.postal_code || "", country: address?.country || "" });
  const { mutateAsync, isPending } = useUpsertBillingAddress();
  const isValid = !!(formData.line1 && formData.city && formData.postal_code && formData.country);
  const setField = (key: keyof BillingAddressFormData, value: string) => setFormData((prev) => ({ ...prev, [key]: value }));

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild><Button variant="ghost" size="icon" className="h-8 w-8"><Pencil className="h-4 w-4" /></Button></DialogTrigger>
      <DialogContent className="sm:max-w-md"><DialogHeader><DialogTitle>{address ? t("views.usage.address.edit_title") : t("views.usage.address.add_title")}</DialogTitle><DialogDescription>{t("views.usage.address.dialog_description")}</DialogDescription></DialogHeader>
        <div className="space-y-4 py-4">
          <AddressField id="line1" label={t("views.usage.address.fields.line1")} value={formData.line1} onChange={(v) => setField("line1", v)} placeholder="123 Main St" />
          <AddressField id="line2" label={t("views.usage.address.fields.line2")} value={formData.line2} onChange={(v) => setField("line2", v)} placeholder={t("views.usage.address.fields.line2_placeholder")} />
          <div className="grid grid-cols-2 gap-4"><AddressField id="city" label={t("views.usage.address.fields.city")} value={formData.city} onChange={(v) => setField("city", v)} placeholder="San Francisco" /><AddressField id="state" label={t("views.usage.address.fields.state")} value={formData.state} onChange={(v) => setField("state", v)} placeholder="CA" /></div>
          <div className="grid grid-cols-2 gap-4"><AddressField id="postal_code" label={t("views.usage.address.fields.postal_code")} value={formData.postal_code} onChange={(v) => setField("postal_code", v)} placeholder="94102" /><AddressField id="country" label={t("views.usage.address.fields.country")} value={formData.country} onChange={(v) => setField("country", v)} placeholder={t("views.usage.address.fields.country_placeholder")} /></div>
        </div>
        <DialogFooter><Button variant="outline" onClick={() => setOpen(false)} disabled={isPending}>{t("common.cancel")}</Button><Button disabled={!isValid || isPending} onClick={async () => { await mutateAsync({ line1: formData.line1, line2: formData.line2 || undefined, city: formData.city, state: formData.state || undefined, postal_code: formData.postal_code, country: formData.country }); setOpen(false); }}>{isPending ? <><Loader2 className="mr-2 h-4 w-4 animate-spin" />{t("common.saving")}</> : t("views.usage.address.save_action")}</Button></DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

function AddressField({ id, label, value, onChange, placeholder }: { id: string; label: string; value: string; onChange: (value: string) => void; placeholder: string }) {
  return <div className="space-y-2"><Label htmlFor={id}>{label}</Label><Input id={id} value={value} onChange={(e) => onChange(e.target.value)} placeholder={placeholder} /></div>;
}

function AddressDisplay({ address }: { address: BillingAddress | null }) {
  const { t } = useTranslation();
  if (!address) return <p className="text-sm text-muted-foreground">{t("views.usage.address.empty")}</p>;
  return <div className="text-sm space-y-0.5"><p>{address.line1}</p>{address.line2 && <p>{address.line2}</p>}<p>{address.city}{address.state && `, ${address.state}`} {address.postal_code}</p><p>{address.country}</p></div>;
}
