import CopyButton from "@/features/developer/components/copy-button";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useGetApp } from "@/hooks";
import {
  useCreateSecret,
  useDisableSecret,
  useListSecrets,
} from "@/hooks/api/secrets";
import { Secret } from "@/models/secret";
import { useState } from "react";

export default function DeveloperView() {
  const { data: app } = useGetApp();

  return (
    <div className="w-full p-4 lg:p-10 flex flex-col gap-6 space-mono-regular">
      <div className="flex flex-row justify-between">
        <h1 className="text-xl font-semibold">Developer</h1>
      </div>
      <div className="grid w-full max-w-sm items-center gap-1.5">
        <Label htmlFor="email">Public key</Label>
        <div className="relative">
          <Input
            type="text"
            id="public_key"
            placeholder="Loading..."
            value={app.data.public_api_key}
          />
          <div className="absolute top-1/2 -translate-y-1/2 right-3">
            <CopyButton text={app.data.public_api_key} />
          </div>
        </div>
      </div>
      <Secrets />
    </div>
  );
}

function Secrets() {
  const { data: secrets } = useListSecrets();

  return (
    <div className="grid w-full max-w-sm items-center gap-1.5">
      <Label htmlFor="email">Secret keys</Label>
      {secrets?.data &&
        secrets.data.map((secret, index) => {
          return (
            <div
              key={secret.id}
              className="grid w-full max-w-sm items-center gap-1.5"
            >
              <div className="flex w-full max-w-sm items-center justify-between space-x-2">
                <Label htmlFor="email">{`Secret key #${index + 1}`}</Label>
                <span className="text-xs text-slate-400">
                  {"" + new Date(secret.created_at).toLocaleDateString()}
                </span>
              </div>
              <div className="flex w-full max-w-sm items-center space-x-2">
                <Input
                  type="text"
                  id={"secret_" + secret.id}
                  placeholder="Loading..."
                  className="pr-2"
                  value="••••••••••••••••••••••••••••••••"
                  readOnly={!secret.disabled_at}
                  disabled={!!secret.disabled_at}
                />
                <DisableApiSecretButton secret={secret} />
              </div>
            </div>
          );
        })}
      <CreateApiSecretButton />
    </div>
  );
}

function DisableApiSecretButton(props: { secret: Secret }) {
  const { mutate } = useDisableSecret(props.secret.id);

  return (
    <Button
      className="min-w-20"
      disabled={!!props.secret.disabled_at}
      onClick={() => {
        mutate({});
      }}
    >
      {!!props.secret.disabled_at ? "Disabled" : "Disable"}
    </Button>
  );
}

function CreateApiSecretButton() {
  const [key, setKey] = useState<string>();
  const [open, setOpen] = useState(false);
  const { mutateAsync } = useCreateSecret();

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button
          variant="outline"
          onClick={async () => {
            const result = await mutateAsync({});
            setKey(result.data.key);
            setOpen(true);
          }}
        >
          Create new secret key
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-md font-mono">
        <DialogHeader>
          <DialogTitle>API secret created</DialogTitle>
          <DialogDescription>
            {
              "Copy and store this key somewhere safe. Once this message is closed you won't be able to see this key again."
            }
          </DialogDescription>
        </DialogHeader>
        <div className="flex items-center space-x-2">
          <div className="grid flex-1 gap-2">
            <Label htmlFor="link" className="sr-only">
              Key
            </Label>
            {!!key && (
              <div className="relative">
                <Textarea
                  id="key"
                  value={key}
                  className="text-xs pr-8"
                  readOnly
                />
                <div className="absolute top-1/2 -translate-y-1/2 right-3">
                  <CopyButton text={key} />
                </div>
              </div>
            )}
          </div>
        </div>
        <DialogFooter className="sm:justify-start">
          <DialogClose asChild>
            <Button type="button" variant="secondary">
              Close
            </Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
