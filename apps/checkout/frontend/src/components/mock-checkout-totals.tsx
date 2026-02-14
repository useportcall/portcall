import { CheckoutSession } from "@/types/api";
import { Check, Loader2, Lock } from "lucide-react";
import { cn } from "@/lib/utils";
import { useTranslation } from "react-i18next";
import { useSubmitState } from "./submit-state";

export function MockCheckoutTotals({ session }: { session: CheckoutSession }) {
  const { i18n, t } = useTranslation();
  const locale = i18n.resolvedLanguage || i18n.language || "en";
  const { state } = useSubmitState();

  const totals = calcTotals(session, locale);

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
      <MockSubmitButton total={totals.total} state={state} t={t} />
    </>
  );
}

function MockSubmitButton({
  total,
  state,
  t,
}: {
  total: string;
  state: "idle" | "loading" | "success";
  t: (key: string, options?: Record<string, unknown>) => string;
}) {
  return (
    <button
      data-testid="checkout-submit-button"
      type="submit"
      disabled={state !== "idle"}
      className={cn(
        "w-full p-3 flex justify-center items-center gap-2 text-white font-semibold mt-2 rounded-md transition-all duration-300",
        state === "idle" && "bg-slate-900 hover:bg-slate-700 cursor-pointer",
        state === "loading" && "bg-slate-700 cursor-wait",
        state === "success" && "bg-emerald-500 cursor-default",
      )}
    >
      {state === "idle" && (
        <>
          <span>{t("submit.pay", { total })}</span>
          <Lock className="w-4 h-4" />
        </>
      )}
      {state === "loading" && (
        <>
          <Loader2 className="w-5 h-5 animate-spin" />
          <span>{t("submit.processing")}</span>
        </>
      )}
      {state === "success" && (
        <>
          <Check className="w-5 h-5" />
          <span>{t("submit.payment_successful")}</span>
        </>
      )}
    </button>
  );
}

function calcTotals(session: CheckoutSession, locale: string) {
  if (!session.plan) return { subtotal: "$0.00", tax: "$0.00", total: "$0.00" };

  const mainPlanItem = session.plan.items.find(
    (item) => item.pricing_model === "fixed",
  );
  const subtotal = mainPlanItem?.unit_amount || 0;
  const tax = 0;
  const total = subtotal + tax;

  const formatAmount = (amount: number, currency: string) => {
    const cur = currency.toUpperCase();
    const divisor = cur === "JPY" || cur === "KRW" ? 1 : 100;
    return new Intl.NumberFormat(locale, {
      style: "currency",
      currency: cur,
    }).format(amount / divisor);
  };

  return {
    subtotal: formatAmount(subtotal, session.plan.currency),
    tax: formatAmount(tax, session.plan.currency),
    total: formatAmount(total, session.plan.currency),
  };
}
