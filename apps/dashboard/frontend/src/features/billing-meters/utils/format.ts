import { BillingMeter } from "@/models/billing-meter";

const PRICING_LABELS: Record<BillingMeter["pricing_model"], string> = {
  unit: "per unit",
  tiered: "tiered pricing",
  block: "block pricing",
  volume: "volume pricing",
};

export function formatUsd(cents: number) {
  return (cents / 100).toLocaleString("en-US", {
    style: "currency",
    currency: "USD",
  });
}

export function formatShortDate(value: string) {
  return new Date(value).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

export function getPricingLabel(model: BillingMeter["pricing_model"]) {
  return PRICING_LABELS[model];
}
