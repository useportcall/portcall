/** Returns placeholder text for provider-specific credential fields. */
export function getPublicKeyPlaceholder(provider: string): string {
  switch (provider) {
    case "stripe":
      return "pk_...";
    case "braintree":
      return "Public key from Braintree Control Panel";
    default:
      return "Add public key";
  }
}

export function getSecretKeyPlaceholder(provider: string): string {
  switch (provider) {
    case "stripe":
      return "sk_...";
    case "braintree":
      return '{"merchant_id":"...","private_key":"..."}';
    default:
      return "Add secret key";
  }
}

export function getSecretKeyLabel(provider: string): string {
  if (provider === "braintree") return "Secret (JSON or merchant_id:private_key)";
  return "Secret Key";
}

export function getWebhookSecretPlaceholder(provider: string): string {
  if (provider === "stripe") return "whsec_...";
  return "Add webhook secret";
}

/** Whether the provider needs manual webhook copy-paste. */
export function needsManualWebhook(provider: string): boolean {
  return provider === "braintree";
}
