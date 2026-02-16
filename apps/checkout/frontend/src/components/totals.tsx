import { formatPlanTotals } from "@/lib/checkout-totals";
import { resolveCheckoutConsentMode } from "@/lib/payment-consent";
import { CheckoutSession } from "@/types/api";
import { Loader2, Lock } from "lucide-react";
import { useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { CheckoutComplianceNotice } from "./checkout-compliance-notice";
import { CheckoutRecurringConsent } from "./checkout-recurring-consent";

type TotalsProps = {
  session: CheckoutSession;
  isSubmitting?: boolean;
  isDisabled?: boolean;
  errorMessage?: string | null;
};

export function Totals({
  session,
  isSubmitting = false,
  isDisabled = false,
  errorMessage,
}: TotalsProps) {
  const { i18n, t } = useTranslation();
  const consentMode = resolveCheckoutConsentMode(session);
  const [hasConsent, setHasConsent] = useState(false);
  const [showConsentError, setShowConsentError] = useState(false);
  const locale = i18n.resolvedLanguage || i18n.language || "en";
  const requiresConsent = consentMode === "save";
  const totals = useMemo(() => {
    if (!session.plan)
      return { subtotal: "$0.00", tax: "$0.00", total: "$0.00" };
    return formatPlanTotals(session.plan, locale);
  }, [session.plan, locale]);

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
      {errorMessage && (
        <div className="w-full p-3 bg-red-50 border border-red-200 rounded-md mt-2" role="alert">
          <p className="text-sm text-red-600">{errorMessage}</p>
        </div>
      )}
      <CheckoutRecurringConsent
        mode={consentMode}
        checked={hasConsent}
        showError={showConsentError}
        onChange={(checked) => {
          setHasConsent(checked);
          if (checked) setShowConsentError(false);
        }}
      />
      <button
        type="submit"
        disabled={isSubmitting || isDisabled}
        onClick={(event) => {
          if (requiresConsent && !hasConsent) {
            event.preventDefault();
            setShowConsentError(true);
          }
        }}
        className="w-full bg-slate-900 p-2 flex justify-center items-center gap-2 text-white font-semibold mt-2 rounded-md hover:bg-slate-700 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
      >
        {isSubmitting ? (
          <>
            <Loader2 className="w-4 h-4 animate-spin" />
            <span>{t("submit.processing")}</span>
          </>
        ) : (
          <>
            <span>
              {consentMode === "save"
                ? t("submit.authorize_and_continue")
                : t("submit.pay", { total: totals.total })}
            </span>
            <Lock className="w-4 h-4" />
          </>
        )}
      </button>
      <CheckoutComplianceNotice session={session} />
    </>
  );
}
