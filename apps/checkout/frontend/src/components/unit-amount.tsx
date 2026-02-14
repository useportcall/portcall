import { Plan, PlanItem } from "@/types/api";
import { getCurrencyDecimalPlaces, getCurrencyDivisor } from "@/lib/currency";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";

export function UnitAmount({
  planItem,
  plan,
}: {
  planItem: PlanItem;
  plan: Plan;
}) {
  const { t, i18n } = useTranslation();
  const locale = i18n.resolvedLanguage || i18n.language || "en";

  const amount = useMemo(() => {
    const decimals = getCurrencyDecimalPlaces(plan.currency);
    const divisor = getCurrencyDivisor(plan.currency);
    const unit = planItem.public_unit_label || t("order_summary.unit");

    const currencyOptions: Intl.NumberFormatOptions = {
      style: "currency",
      currency: plan.currency.toUpperCase(),
      minimumFractionDigits: decimals,
    };

    if (planItem.unit_amount && planItem.pricing_model === "fixed") {
      return (planItem.unit_amount / divisor).toLocaleString(
        locale,
        currencyOptions,
      );
    }

    if (planItem.unit_amount && planItem.pricing_model === "unit") {
      return (
        (planItem.unit_amount / divisor).toLocaleString(
          locale,
          currencyOptions,
        ) +
        " / " +
        unit
      );
    }

    if (planItem.tiers?.length) {
      const tier = planItem.tiers[0];
      const value = (tier.unit_amount / divisor).toLocaleString(
        locale,
        currencyOptions,
      );
      return t("order_summary.starting_at", { value, unit });
    }

    return t("order_summary.free");
  }, [planItem, plan.currency, locale, t]);

  return amount;
}
