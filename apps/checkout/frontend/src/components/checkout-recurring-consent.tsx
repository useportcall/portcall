import { CheckoutConsentMode } from "@/lib/payment-consent";
import { useTranslation } from "react-i18next";

export function CheckoutRecurringConsent(props: {
  mode: CheckoutConsentMode;
  checked: boolean;
  showError: boolean;
  onChange: (checked: boolean) => void;
}) {
  const { t } = useTranslation();
  if (props.mode !== "save") return null;

  return (
    <div className="mt-3">
      <label className="flex items-start gap-2 text-xs text-slate-700">
        <input
          type="checkbox"
          checked={props.checked}
          onChange={(e) => props.onChange(e.target.checked)}
          aria-describedby="recurring-consent-error"
          className="mt-0.5 h-4 w-4"
        />
        <span>{t("compliance.consent_label")}</span>
      </label>
      {props.showError ? (
        <p id="recurring-consent-error" role="alert" className="mt-1 text-xs text-red-600">
          {t("validation.recurring_consent_required")}
        </p>
      ) : null}
    </div>
  );
}
