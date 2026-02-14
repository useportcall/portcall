import { LabeledInput } from "@/components/base-view";
import { Button } from "@/components/ui/button";
import {
  Dialog, DialogContent, DialogDescription,
  DialogHeader, DialogTitle, DialogTrigger,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useCreateConnection } from "@/hooks/api/connections";
import { zodResolver } from "@hookform/resolvers/zod";
import { Loader2, Plus } from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { getProviderDisplayName, getProviderIcon } from "./provider-meta";
import { ProviderSelectItems } from "./provider-select-items";
import { ConnectionCredentialFields } from "./connection-credential-fields";
import { ProviderInfoBanner } from "./provider-info-banner";

const schema = z.object({
  name: z.string(), provider: z.string().min(1),
  public_key: z.string(), secret_key: z.string(),
  webhook_secret: z.string().optional(),
});

export function CreateConnectionDialog() {
  const { mutateAsync, isPending } = useCreateConnection();
  const [open, setOpen] = useState(false);
  const form = useForm({
    resolver: zodResolver(schema),
    defaultValues: { name: "", provider: "", public_key: "", secret_key: "", webhook_secret: "" },
  });
  const selectedProvider = form.watch("provider");

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <div className="pr-1">
          <Button data-testid="add-provider-button" size="sm" variant="outline" className="gap-2">
            <Plus className="w-4 h-4" /> Add Provider
          </Button>
        </div>
      </DialogTrigger>
      <DialogContent className="flex flex-col gap-4 sm:max-w-[480px]">
        <DialogHeader>
          <DialogTitle>Add Payment Provider</DialogTitle>
          <DialogDescription>
            Connect a payment provider to start accepting payments from your customers.
          </DialogDescription>
        </DialogHeader>
        <form className="w-full flex flex-col gap-5" onSubmit={form.handleSubmit(async (f) => {
          await mutateAsync({
            name: f.name || getProviderDisplayName(f.provider),
            source: f.provider, public_key: f.public_key,
            secret_key: f.secret_key, webhook_secret: f.webhook_secret || "",
          });
          form.reset(); setOpen(false);
        })}>
          <div className="w-full flex flex-col gap-2">
            <Label htmlFor="provider">Payment Provider</Label>
            <Select value={selectedProvider} onValueChange={(v) => {
              form.setValue("provider", v);
              if (!form.getValues("name")) form.setValue("name", getProviderDisplayName(v));
            }}>
              <SelectTrigger id="provider" className="h-12">
                <SelectValue placeholder="Select a provider">
                  {selectedProvider && (
                    <div className="flex items-center gap-3">
                      {getProviderIcon(selectedProvider, "w-6 h-6")}
                      <span>{getProviderDisplayName(selectedProvider)}</span>
                    </div>
                  )}
                </SelectValue>
              </SelectTrigger>
              <SelectContent><ProviderSelectItems /></SelectContent>
            </Select>
          </div>
          <LabeledInput id="name" label="Display Name"
            placeholder={selectedProvider ? getProviderDisplayName(selectedProvider) : "e.g., Stripe Production"}
            className="outline-none w-full" helperText="Optional - defaults to provider name"
            {...form.register("name")} />
          <ConnectionCredentialFields provider={selectedProvider} form={form} />
          <ProviderInfoBanner provider={selectedProvider} />
          <div className="flex flex-row gap-2 justify-end pt-2 border-t">
            <Button type="button" variant="ghost" onClick={() => { form.reset(); setOpen(false); }}>Cancel</Button>
            <Button data-testid="submit-provider-button" type="submit" disabled={isPending || !selectedProvider} className="gap-2">
              {isPending ? (<><Loader2 className="w-4 h-4 animate-spin" /> Adding...</>) : (<><Plus className="w-4 h-4" /> Add Provider</>)}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
