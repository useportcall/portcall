import { CheckoutSession, Plan } from "@/types/api";

export function formatPlanTotals(plan: Plan, locale: string) {
  const mainPlanItem = plan.items.find((item) => item.pricing_model === "fixed");
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

export function formatSessionTotals(session: CheckoutSession, locale: string) {
  if (!session.plan) return { subtotal: "$0.00", tax: "$0.00", total: "$0.00" };
  return formatPlanTotals(session.plan, locale);
}
