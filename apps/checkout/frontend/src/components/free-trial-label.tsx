import { Plan } from "@/types/api";
import { useTranslation } from "react-i18next";
import { UnitAmount } from "./unit-amount";

function calculateStartDate(trialPeriodDays: number, locale: string): string {
  const currentDate = new Date();
  currentDate.setDate(currentDate.getDate() + trialPeriodDays);
  return currentDate.toLocaleDateString(locale);
}

export function FreeTrialLabel({ plan }: { plan: Plan }) {
  const { t, i18n } = useTranslation();
  const locale = i18n.resolvedLanguage || i18n.language || "en";
  if (!plan.trial_period_days) {
    return null;
  }

  return (
    <p className="text-xs text-slate-400">
      <span>{t("order_summary.then")} </span>
      <UnitAmount plan={plan} planItem={plan.items[0]} />
      <span> {t(`intervals.per_${plan.interval}`)} </span>
      <span>
        {t("order_summary.starting_on", {
          date: calculateStartDate(plan.trial_period_days, locale),
        })}
      </span>
    </p>
  );
}
