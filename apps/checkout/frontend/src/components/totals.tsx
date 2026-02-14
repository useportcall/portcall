import { CheckoutSession, Plan } from "@/types/api";
import { Loader2, Lock } from "lucide-react";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";

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
  const locale = i18n.resolvedLanguage || i18n.language || "en";
  const totals = useMemo(() => {
    if (!session.plan)
      return { subtotal: "$0.00", tax: "$0.00", total: "$0.00" };
    return formatTotals(session.plan, locale);
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
        <div className="w-full p-3 bg-red-50 border border-red-200 rounded-md mt-2">
          <p className="text-sm text-red-600">{errorMessage}</p>
        </div>
      )}
      <button
        type="submit"
        disabled={isSubmitting || isDisabled}
        className="w-full bg-slate-900 p-2 flex justify-center items-center gap-2 text-white font-semibold mt-2 rounded-md hover:bg-slate-700 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
      >
        {isSubmitting ? (
          <>
            <Loader2 className="w-4 h-4 animate-spin" />
            <span>{t("submit.processing")}</span>
          </>
        ) : (
          <>
            <span>{t("submit.pay", { total: totals.total })}</span>
            <Lock className="w-4 h-4" />
          </>
        )}
      </button>
    </>
  );
}

function formatTotals(plan: Plan, locale: string) {
  const mainPlanItem = plan.items.find(
    (item) => item.pricing_model === "fixed",
  );
  const subtotal = mainPlanItem?.unit_amount || 0;
  const tax = 0;
  const total = subtotal + tax;
  const fmt = (amount: number) => {
    const cur = plan.currency.toUpperCase();
    const divisor = cur === "JPY" || cur === "KRW" ? 1 : 100;
    return new Intl.NumberFormat(locale, {
      style: "currency",
      currency: cur,
    }).format(amount / divisor);
  };
  return { subtotal: fmt(subtotal), tax: fmt(tax), total: fmt(total) };
}
