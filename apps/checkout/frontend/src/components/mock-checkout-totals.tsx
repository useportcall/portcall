import { CheckoutSession } from "@/types/api";
import { useTranslation } from "react-i18next";
import { useSubmitState } from "./submit-state";
import { formatSessionTotals } from "@/lib/checkout-totals";
import { resolveCheckoutConsentMode } from "@/lib/payment-consent";
import { useState } from "react";
import { CheckoutRecurringConsent } from "./checkout-recurring-consent";
import { MockSubmitButton } from "./mock-submit-button";

export function MockCheckoutTotals({ session }: { session: CheckoutSession }) {
  const { i18n, t } = useTranslation();
  const consentMode = resolveCheckoutConsentMode(session);
  const [hasConsent, setHasConsent] = useState(false);
  const [showConsentError, setShowConsentError] = useState(false);
  const locale = i18n.resolvedLanguage || i18n.language || "en";
  const { state } = useSubmitState();
  const requiresConsent = consentMode === "save";

  const totals = formatSessionTotals(session, locale);

  return (
    <>
      <div className="flex flex-col justify-start items-start gap-4 mt-4">
        <span className="w-full flex text-sm justify-between">
          <p>{t("subtotal")}</p> <p>{totals.subtotal}</p>
        </span>
        <span className="w-full text-sm flex justify-between">
          <p>{t("form.taxes")}</p> <p>{totals.tax}</p>
        </span>
        <span className="w-full font-semibold flex justify-between">
          <p>{t("total_due")}</p> <p>{totals.total}</p>
        </span>
      </div>
      <CheckoutRecurringConsent
        mode={consentMode}
        checked={hasConsent}
        showError={showConsentError}
        onChange={(checked) => {
          setHasConsent(checked);
          if (checked) setShowConsentError(false);
        }}
      />
      <MockSubmitButton
        total={totals.total}
        state={state}
        t={t}
        consentMode={consentMode}
        onTrySubmit={(event) => {
          if (requiresConsent && !hasConsent) {
            event.preventDefault();
            setShowConsentError(true);
          }
        }}
      />
    </>
  );
}
