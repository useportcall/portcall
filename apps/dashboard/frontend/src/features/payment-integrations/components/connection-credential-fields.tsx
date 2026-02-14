import { LabeledInput } from "@/components/base-view";
import { UseFormReturn } from "react-hook-form";
import {
  getPublicKeyPlaceholder,
  getSecretKeyLabel,
  getSecretKeyPlaceholder,
  getWebhookSecretPlaceholder,
  needsManualWebhook,
} from "./provider-field-config";

type Props = {
  provider: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  form: UseFormReturn<any>;
};

/** Credential fields shown for non-local providers. */
export function ConnectionCredentialFields({ provider, form }: Props) {
  if (!provider || provider === "local") return null;
  return (
    <>
      <LabeledInput
        type="text"
        id="public_key"
        label="Publishable Key"
        placeholder={getPublicKeyPlaceholder(provider)}
        className="outline-none w-full text-sm font-mono"
        helperText=""
        {...form.register("public_key")}
      />
      <LabeledInput
        type="password"
        id="secret_key"
        label={getSecretKeyLabel(provider)}
        placeholder={getSecretKeyPlaceholder(provider)}
        className="outline-none w-full text-sm font-mono"
        helperText="Your secret key is encrypted and stored securely"
        {...form.register("secret_key")}
      />
      {!needsManualWebhook(provider) && (
        <LabeledInput
          type="password"
          id="webhook_secret"
          label="Webhook Signing Secret"
          placeholder={getWebhookSecretPlaceholder(provider)}
          className="outline-none w-full text-sm font-mono"
          helperText="Optional - Required to receive payment events."
          {...form.register("webhook_secret")}
        />
      )}
    </>
  );
}
