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
import CopyButton from "@/features/developer/components/copy-button";
import { useGetApp } from "@/hooks";
import {
  useCreateSecret,
  useDisableSecret,
  useListSecrets,
} from "@/hooks/api/secrets";
import { Secret } from "@/models/secret";
import { useState } from "react";
import { useTranslation } from "react-i18next";

export default function DeveloperView() {
  const { t } = useTranslation();
  const { data: app } = useGetApp();
  return (
    <div className="w-full p-4 lg:p-10 flex flex-col gap-6 space-mono-regular">
      <div className="flex flex-row justify-between">
        <h1 className="text-xl font-semibold">{t("views.developer.title")}</h1>
      </div>
      <div className="grid w-full max-w-sm items-center gap-1.5">
        <Label htmlFor="public_key">{t("views.developer.public_key")}</Label>
        <div className="relative">
          <Input
            type="text"
            id="public_key"
            value={app.data.public_api_key}
            readOnly
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
  const { t } = useTranslation();
  const { data: secrets } = useListSecrets();
  return (
    <div className="grid w-full max-w-sm items-center gap-1.5">
      <Label htmlFor="secret_keys">{t("views.developer.secret_keys")}</Label>
      {secrets?.data?.map((secret, index) => (
        <SecretRow key={secret.id} secret={secret} index={index} />
      ))}
      <CreateApiSecretButton />
    </div>
  );
}

function SecretRow({ secret, index }: { secret: Secret; index: number }) {
  const { t } = useTranslation();
  return (
    <div className="grid w-full max-w-sm items-center gap-1.5">
      <div className="flex w-full max-w-sm items-center justify-between space-x-2">
        <Label>{t("views.developer.secret_key_number", { index: index + 1 })}</Label>
        <span className="text-xs text-muted-foreground">
          {new Date(secret.created_at).toLocaleDateString()}
        </span>
      </div>
      <div className="flex w-full max-w-sm items-center space-x-2">
        <Input
          type="text"
          id={`secret_${secret.id}`}
          value="••••••••••••••••••••••••••••••••"
          readOnly
          disabled={!!secret.disabled_at}
        />
        <DisableApiSecretButton secret={secret} />
      </div>
    </div>
  );
}

function DisableApiSecretButton({ secret }: { secret: Secret }) {
  const { t } = useTranslation();
  const { mutate } = useDisableSecret(secret.id);
  return (
    <Button
      className="min-w-20"
      disabled={!!secret.disabled_at}
      onClick={() => mutate({})}
    >
      {secret.disabled_at ? t("views.developer.disabled") : t("views.developer.disable")}
    </Button>
  );
}

function CreateApiSecretButton() {
  const { t } = useTranslation();
  const [key, setKey] = useState<string>();
  const [open, setOpen] = useState(false);
  const { mutateAsync } = useCreateSecret();

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button
          data-testid="create-secret-button"
          variant="outline"
          onClick={async () => {
            const result = await mutateAsync({});
            setKey(result.data.key);
            setOpen(true);
          }}
        >
          {t("views.developer.create_secret")}
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-md font-mono">
        <DialogHeader>
          <DialogTitle>{t("views.developer.secret_created_title")}</DialogTitle>
          <DialogDescription>
            {t("views.developer.secret_created_description")}
          </DialogDescription>
        </DialogHeader>
        <div className="flex items-center space-x-2">
          <div className="grid flex-1 gap-2">
            <Label htmlFor="key" className="sr-only">
              Key
            </Label>
            {key && (
              <div className="relative">
                <Textarea
                  data-testid="secret-key-value"
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
            <Button
              data-testid="close-secret-dialog"
              type="button"
              variant="secondary"
            >
              {t("common.close")}
            </Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
