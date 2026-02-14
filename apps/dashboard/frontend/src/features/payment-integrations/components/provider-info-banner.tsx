/** Info banners shown for specific providers in the connection dialog. */
export function ProviderInfoBanner({ provider }: { provider: string }) {
  if (provider === "local") {
    return (
      <div className="rounded-lg border border-amber-200 bg-amber-50 p-3">
        <p className="text-sm text-amber-800">
          <strong>Test Mode:</strong> The mock provider simulates payment flows
          without processing real transactions.
        </p>
      </div>
    );
  }
  if (provider === "braintree") {
    return (
      <div className="rounded-lg border border-blue-200 bg-blue-50 p-3">
        <p className="text-sm text-blue-800">
          <strong>Manual webhook setup:</strong> After adding, copy the webhook
          URL from the connection card and paste it into the Braintree Control
          Panel → Settings → Webhooks.
        </p>
      </div>
    );
  }
  return null;
}
