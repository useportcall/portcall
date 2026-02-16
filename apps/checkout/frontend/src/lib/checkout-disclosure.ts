import { Plan } from "@/types/api";

function addInterval(date: Date, interval: string, count: number) {
  if (interval === "day") date.setDate(date.getDate() + count);
  if (interval === "week") date.setDate(date.getDate() + 7 * count);
  if (interval === "month") date.setMonth(date.getMonth() + count);
  if (interval === "year") date.setFullYear(date.getFullYear() + count);
}

export function recurringAmountMode(plan: Plan): "fixed" | "variable" {
  if (plan.metered_features.some((f) => f.plan_item.unit_amount > 0)) {
    return "variable";
  }
  if (plan.items.some((item) => item.pricing_model !== "fixed")) {
    return "variable";
  }
  return "fixed";
}

export function cadenceLabel(plan: Plan, intervalLabel: string): string {
  if (plan.interval_count <= 1) return intervalLabel;
  return `${plan.interval_count} ${intervalLabel}`;
}

export function renewalDate(plan: Plan): Date {
  const date = new Date();
  if (plan.trial_period_days > 0) {
    date.setDate(date.getDate() + plan.trial_period_days);
    return date;
  }
  addInterval(date, plan.interval, plan.interval_count || 1);
  return date;
}
