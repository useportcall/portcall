import { Plan } from "@/types/api";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";

export function IntervalLabel({ plan }: { plan: Plan }) {
  const { t } = useTranslation();
  const label = useMemo(() => {
    switch (plan.interval) {
      case "day":
        return t("intervals.day");
      case "week":
        return t("intervals.week");
      case "month":
        return t("intervals.month");
      case "year":
        return t("intervals.year");
      default:
        return t("order_summary.unknown_interval");
    }
  }, [plan, t]);

  if (plan.trial_period_days) {
    return null;
  }

  return (
    <p className="text-xs text-slate-400">
      {t("order_summary.recurs", { interval: label })}
    </p>
  );
}
