import { useMemo } from "react";
import { useTranslation } from "react-i18next";

export function QuotaLabel(props: { quota: number; interval: string }) {
  const { t } = useTranslation();
  const interval = useMemo(() => {
    switch (props.interval) {
      case "day":
        return t("intervals.per_day");
      case "week":
        return t("intervals.per_week");
      case "month":
        return t("intervals.per_month");
      case "year":
        return t("intervals.per_year");
      default:
        return "";
    }
  }, [props.interval, t]);

  if (!props.quota) {
    return null;
  }

  if (props.quota === -1) {
    return (
      <p className="text-xs text-slate-400">
        <span>{t("order_summary.unlimited")} </span>
      </p>
    );
  }

  return (
    <p className="text-xs text-slate-400">
      {t("order_summary.limit", { quota: props.quota, interval: interval })}
    </p>
  );
}
