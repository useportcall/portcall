import EmptyTable from "@/components/empty-table";
import { BaseSection, BaseView, LabeledInput } from "@/components/base-view";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardHeader } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useAppConfig, useGetAppConfig } from "@/hooks/api/app-configs";
import {
  useCreateConnection,
  useDeleteConnection,
  useListConnections,
} from "@/hooks/api/connections";
import { zodResolver } from "@hookform/resolvers/zod";

import { Loader2, MoreVertical, Plus } from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import z from "zod";

export default function PaymentIntegrationsView() {
  return (
    <BaseView
      title="Payment Integrations"
      description="Manage your payment integrations here."
    >
      <BaseSection title="Intergrations" actions={<CreateConnection />}>
        <Connections />
      </BaseSection>
    </BaseView>
  );
}

export const schema = z.object({
  name: z.string(),
  provider: z.string().min(1, "Provider is required"),
  public_key: z.string(),
  secret_key: z.string(),
});

function CreateConnection() {
  const { mutateAsync, isPending } = useCreateConnection();

  const [showForm, setShowForm] = useState(false);
  const form = useForm({
    resolver: zodResolver(schema),
    defaultValues: {
      name: "New Provider",
      provider: "local",
      public_key: "",
      secret_key: "",
    },
  });

  return (
    <Dialog open={showForm} onOpenChange={setShowForm}>
      <DialogTrigger asChild>
        <div className="pr-1">
          <Button
            size="icon"
            variant="ghost"
            className="size-[22px] rounded-md hover:bg-slate-100"
            onClick={() => setShowForm(true)}
          >
            <Plus className="" />
          </Button>
        </div>
      </DialogTrigger>
      <DialogContent className="flex flex-col gap-4">
        <DialogHeader>
          <DialogTitle>Add a payment integration</DialogTitle>
          <DialogDescription>
            Connect your payment provider to start accepting payments
            seamelessly.
          </DialogDescription>
        </DialogHeader>
        <form
          onSubmit={form.handleSubmit(async (fields) => {
            await mutateAsync({
              name: fields.name,
              source: fields.provider,
              public_key: fields.public_key,
              secret_key: fields.secret_key,
              webhook_secret: "",
            });
            form.reset();
            setShowForm(false);
          })}
          className="w-full flex flex-col gap-4"
        >
          <BaseSection title="">
            <LabeledInput
              autoFocus
              id="name"
              label="Connection Name"
              placeholder="Add name"
              className="font-semibold outline-none w-full"
              helperText=""
              {...form.register("name")}
            />
            <div className="w-full flex flex-col gap-2 justify-start items-start">
              <Label htmlFor="provider">Payment provider</Label>
              <p className="text-xs text-slate-400">
                {`Choose a payment provider from the list below. If you're not
                ready to go live yet, feel free to choose the mock provider
                option.`}
              </p>
              <Select
                value={form.watch("provider")}
                onValueChange={(value) => form.setValue("provider", value)}
              >
                <SelectTrigger id="provider">
                  <SelectValue placeholder="Select a provider" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="local">Mock</SelectItem>
                  <SelectItem value="stripe">Stripe</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <LabeledInput
              type="text"
              id="public_key"
              label="Public Key"
              placeholder="Add public key"
              className="outline-none w-full text-sm"
              helperText=""
              {...form.register("public_key")}
            />
            <LabeledInput
              type="password"
              id="secret_key"
              label="Secret Key"
              placeholder="Add secret key"
              className="outline-none w-full text-sm"
              helperText=""
              {...form.register("secret_key")}
            />
          </BaseSection>
          <div className="w-fit mt-4 flex flex-row gap-2 justify-end self-end">
            <Button
              type="button"
              variant="ghost"
              onClick={() => setShowForm(false)}
            >
              Cancel
            </Button>
            <Button size="sm" disabled={isPending}>
              Add integration
              {isPending && <Loader2 className="size-4 animate-spin" />}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}

function Connections() {
  const { data: config } = useGetAppConfig();
  const { data: connections } = useListConnections();

  if (connections && !connections.data.length) {
    return (
      <EmptyTable message="No payment integrations created yet." button={""} />
    );
  }

  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-4 max-w-fit">
      {connections.data.map((connection) => {
        return (
          <Card key={connection.id}>
            <CardHeader className="flex flex-col justify-start items-start gap-2 relative">
              <span className="flex flex-row justify-between">
                <span className="flex flex-col gap-1">
                  <div className="w-10 h-10 rounded-full border-2 border-slate-200">
                    {connection.source}
                  </div>
                  <p className="outline-black">{connection.name}</p>;
                </span>
                <div className="absolute top-4 right-3">
                  <ActionsMenu connection={connection} />
                </div>
              </span>
              <IsDefault
                checked={connection.id === config?.data?.default_connection}
              />
            </CardHeader>
          </Card>
        );
      })}
    </div>
  );
}

function IsDefault({ checked }: { checked: boolean }) {
  if (!checked) return <></>;

  return (
    <Badge variant={"outline"} className="w-fit">
      ✅ Default
    </Badge>
  );
}

function ActionsMenu({ connection }: { connection: any }) {
  const { data: config } = useGetAppConfig();

  if (!config) return <></>;

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button size="icon" variant="ghost" className="size-8 p-0">
          <MoreVertical className="w-5 h-5" />
        </Button>
      </DropdownMenuTrigger>

      <DropdownMenuContent align="start">
        <MakeDefault
          id={connection.id}
          checked={connection.id === config?.data?.default_connection}
        />
        <DeleteConnection id={connection.id} />
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

function DeleteConnection({ id }: { id: string }) {
  const { mutate, isPending } = useDeleteConnection(id);

  return (
    <DropdownMenuItem disabled={isPending} onClick={() => mutate({})}>
      {isPending ? (
        <span className="flex items-center gap-2">
          <Loader2 className="animate-spin w-4 h-4" /> Deleting...
        </span>
      ) : (
        "Delete"
      )}
    </DropdownMenuItem>
  );
}

function MakeDefault({ id, checked }: { id: string; checked: boolean }) {
  const { mutateAsync: updateAppConfig } = useAppConfig();

  return (
    <DropdownMenuItem
      disabled={checked}
      onClick={() => updateAppConfig({ connection_id: id })}
    >
      Make default
    </DropdownMenuItem>
  );
}
