import { formatSessionTotals } from "@/lib/checkout-totals";
import {
  cadenceLabel,
  recurringAmountMode,
  renewalDate,
} from "@/lib/checkout-disclosure";
import { resolveCheckoutConsentMode } from "@/lib/payment-consent";
import { CheckoutSession } from "@/types/api";
import { useTranslation } from "react-i18next";

export function CheckoutComplianceNotice({
  session,
}: {
  session: CheckoutSession;
}) {
  const { i18n, t } = useTranslation();
  const locale = i18n.resolvedLanguage || i18n.language || "en";
  const mode = resolveCheckoutConsentMode(session);
  const merchant = session.company?.name || t("compliance.merchant_fallback");
  const plan = session.plan;
  const totals = formatSessionTotals(session, locale);
  const interval = t("intervals." + (plan?.interval || "month"));
  const amountMode = plan ? recurringAmountMode(plan) : "fixed";
  const cadence = plan ? cadenceLabel(plan, interval) : interval;
  const nextRenewal = plan
    ? renewalDate(plan).toLocaleDateString(locale, { dateStyle: "long" })
    : "-";
  const termsURL = process.env.NEXT_PUBLIC_CHECKOUT_TERMS_URL || session.cancel_url;
  const privacyURL = process.env.NEXT_PUBLIC_CHECKOUT_PRIVACY_URL || "https://useportcall.com/privacy";
  const cancelURL = process.env.NEXT_PUBLIC_CHECKOUT_CANCELLATION_URL || session.cancel_url || termsURL;

  return (
    <div className="mt-4 rounded-md border border-slate-200 bg-slate-50 p-3 text-xs text-slate-600">
      <p className="font-medium text-slate-700">{t("compliance.title")}</p>
      <p className="mt-1">
        {mode === "save"
          ? t("compliance.save_payment_method", { merchant })
          : t("compliance.one_time_charge", { merchant })}
      </p>
      <p className="mt-1">
        {t("compliance.recurring_amount", {
          amount: amountMode === "variable" ? t("compliance.variable_amount") : totals.total,
        })}
      </p>
      <p className="mt-1">{t("compliance.billing_cadence", { cadence })}</p>
      <p className="mt-1">{t("compliance.renewal_date", { date: nextRenewal })}</p>
      {plan?.trial_period_days ? (
        <p className="mt-1">{t("compliance.trial_details", { days: plan.trial_period_days })}</p>
      ) : null}
      <p className="mt-1">{t("compliance.tax_fee_disclosure", { tax: totals.tax })}</p>
      {mode === "save" ? (
        <p className="mt-1">{t("compliance.off_session_notice")}</p>
      ) : null}
      <p className="mt-1">{t("compliance.strong_authentication")}</p>
      <p className="mt-1">{t("compliance.terms_acknowledgement")}</p>
      <div className="mt-2 flex flex-wrap gap-3 text-xs">
        {termsURL ? <a className="underline" href={termsURL} target="_blank" rel="noreferrer">{t("compliance.links.terms")}</a> : null}
        {privacyURL ? <a className="underline" href={privacyURL} target="_blank" rel="noreferrer">{t("compliance.links.privacy")}</a> : null}
        {cancelURL ? <a className="underline" href={cancelURL} target="_blank" rel="noreferrer">{t("compliance.links.cancellation")}</a> : null}
      </div>
    </div>
  );
}
