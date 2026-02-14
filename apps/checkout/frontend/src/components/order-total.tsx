import { CheckoutSession } from "@/types/api";
import { useMemo } from "react";
import { useTranslation } from "react-i18next";

export default function OrderTotal({ session }: { session: CheckoutSession }) {
  const { t, i18n } = useTranslation();
  const total = useMemo(() => {
    const subtotal = calculateSubtotal(session);
    const tax = subtotal * 0;
    return subtotal + tax;
  }, [session]);
  const amount = useMemo(() => formatTotal(total, session, i18n.resolvedLanguage || i18n.language || "en"), [i18n.language, i18n.resolvedLanguage, session, total]);

  return (
    <div className="">
      <p className="text-slate-800">{t("order_summary.pay_today", { amount })}</p>
    </div>
  );
}

function calculateSubtotal(session: CheckoutSession) {
  const { plan } = session;
  if (!plan) return 0;

  const subtotal = plan.items.reduce((acc, item) => {
    if (item.pricing_model !== "fixed") {
      return acc;
    }
    return acc + item.unit_amount * item.quantity;
  }, 0);
  return subtotal;
}

function formatTotal(total: number, session: CheckoutSession, locale: string): string {
  const currency = (session.plan?.currency || "usd").toUpperCase();
  const divisor = currency === "JPY" || currency === "KRW" ? 1 : 100;
  return new Intl.NumberFormat(locale, { style: "currency", currency }).format(total / divisor);
}
