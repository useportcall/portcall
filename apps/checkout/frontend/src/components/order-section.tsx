import { CheckoutSession, Plan, PlanItem } from "@/types/api";
import { useMemo } from "react";
import OrderTotal from "./order-total";
import { Separator } from "./ui/separator";
import { ChevronRight } from "lucide-react";

function IntervalLabel({ plan }: { plan: Plan }) {
  const label = useMemo(() => {
    switch (plan.interval) {
      case "day":
        return "Daily";
      case "week":
        return "Weekly";
      case "month":
        return "Monthly";
      case "year":
        return "Yearly";
      default:
        return "Unknown";
    }
  }, [plan]);

  if (plan.trial_period_days) {
    return null;
  }

  return <p className="text-xs text-slate-400">Recurs {label}</p>;
}

function calculateStartDate(trialPeriodDays: number): string {
  const currentDate = new Date();
  currentDate.setDate(currentDate.getDate() + trialPeriodDays);
  return currentDate.toLocaleDateString();
}

function FreeTrialLabel({ plan }: { plan: Plan }) {
  if (!plan.trial_period_days) {
    return null;
  }

  return (
    <p className="text-xs text-slate-400">
      <span>Then </span>
      <UnitAmount plan={plan} planItem={plan.items[0]} />
      <span>/</span>
      <span> {plan.interval} </span>
      starting on <span>{calculateStartDate(plan.trial_period_days)}</span>
    </p>
  );
}

function toTitleCase(str: string) {
  return str
    .split("_")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");
}

export default function OrderSummary({
  session,
}: {
  session: CheckoutSession;
}) {
  const plan = session.plan;

  if (!plan) return null;

  return (
    <div className="text-slate-800 flex flex-col gap-4 w-full md:max-w-sm mx-auto">
      <h1 className="text-2xl font-medium">{session.company?.name}</h1>
      <div className="flex flex-col gap-2 mt-2 bg-white rounded-lg p-6 border border-slate-200">
        <p className="font-semibold">{session.plan?.name}</p>
        <OrderTotal session={session} />
        <IntervalLabel plan={plan} />
        <FreeTrialLabel plan={plan} />
      </div>
      <div className="mt-4 gap-4 flex flex-col w-full text-sm">
        <h2>Here&apos;s what you get:</h2>
        <div className="gap-8 flex flex-col w-full bg-white rounded-lg p-6 border border-slate-200">
          <div className="flex flex-col gap-2 w-full">
            <p className="font-medium">Features</p>
            <ul>
              {!plan.features.length && <li>No features included</li>}
              {plan?.features.map((feature) => (
                <li key={feature.id}>
                  <span className="flex gap-1 justify-start items-center">
                    <span className="w-4 h-4">
                      <ChevronRight className="w-4 h-4" />
                    </span>
                    {/* convert snake case to first letter uppercase with spaces */}
                    {toTitleCase(feature.id)}
                  </span>
                </li>
              ))}
            </ul>
          </div>
          <Separator />
          <div className="flex flex-col gap-2 w-full text-sm">
            <p className="font-medium">Metered features</p>
            <div className="flex flex-col gap-4">
              {session.plan?.metered_features.length === 0 && (
                <p>No metered features included</p>
              )}
              {session.plan?.metered_features.map((feature) => (
                <div key={feature.id} className="flex flex-col gap-1">
                  <span className="flex flex-wrap w-full justify-between">
                    <span className="flex gap-1 justify-start items-center">
                      <span className="w-4 h-4">
                        <ChevronRight className="w-4 h-4" />
                      </span>
                      {toTitleCase(feature.feature.id)}
                    </span>
                    {feature.plan_item.unit_amount !== 0 && (
                      <span className="text-xs">
                        <UnitAmount plan={plan} planItem={feature.plan_item} />
                      </span>
                    )}
                  </span>
                  <QuotaLabel
                    quota={feature.quota}
                    interval={feature.interval}
                  />
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

// Show limit if there is limit
function QuotaLabel(props: { quota: number; interval: string }) {
  const interval = useMemo(() => {
    switch (props.interval) {
      case "day":
        return "per day";
      case "week":
        return "per week";
      case "month":
        return "per month";
      case "year":
        return "per year";
      default:
        return "";
    }
  }, [props.interval]);

  if (!props.quota) {
    return null;
  }

  if (props.quota === -1) {
    return (
      <p className="text-xs text-slate-400">
        <span>Unlimited </span>
      </p>
    );
  }

  return (
    <p className="text-xs text-slate-400">
      <span>Limit: </span>
      <span>{props.quota} </span>
      <span>{interval}</span>
    </p>
  );
}

function UnitAmount({ planItem, plan }: { planItem: PlanItem; plan: Plan }) {
  const amount = useMemo(() => {
    const decimals = getCurrencyDecimalPlaces(plan.currency);
    const divisor = getCurrencyDivisor(plan.currency);
    const unit = planItem.public_unit_label || "unit";

    const currencyOptions: Intl.NumberFormatOptions = {
      style: "currency",
      currency: plan.currency.toUpperCase(),
      minimumFractionDigits: decimals,
    };

    if (planItem.unit_amount && planItem.pricing_model === "fixed") {
      return (planItem.unit_amount / divisor).toLocaleString(
        "en-US",
        currencyOptions
      );
    }

    if (planItem.unit_amount && planItem.pricing_model === "unit") {
      return (
        (planItem.unit_amount / divisor).toLocaleString(
          "en-US",
          currencyOptions
        ) +
        " / " +
        unit
      );
    }

    if (planItem.tiers?.length) {
      const tier = planItem.tiers[0];
      const value = (tier.unit_amount / divisor).toLocaleString(
        "en-US",
        currencyOptions
      );

      return "starting at " + value + " / " + unit;
    }

    return "free";
  }, [planItem]);

  return amount;
}

function getCurrencyDivisor(currency: string): number {
  switch (currency.toLowerCase()) {
    case "jpy":
    case "krw":
      return 1;
    default:
      return 100;
  }
}

function getCurrencyDecimalPlaces(currency: string): number {
  switch (currency.toLowerCase()) {
    case "jpy":
    case "krw":
      return 0;
    default:
      return 2;
  }
}

export function toCurrencySymbol(currency: string) {
  switch (currency) {
    case "usd":
      return "$";
    case "eur":
      return "€";
    case "gbp":
      return "£";
    case "jpy":
      return "¥";
    case "aud":
      return "A$";
    case "cad":
      return "C$";
    case "cny":
      return "¥";
    case "inr":
      return "₹";
    case "rub":
      return "₽";
    case "brl":
      return "R$";
    case "mxn":
      return "MX$";
    case "sgd":
      return "S$";
    case "hkd":
      return "HK$";
    case "krw":
      return "₩";
    case "sek":
    case "dkk":
    case "nok":
      return "kr";
    default:
      return currency.toUpperCase();
  }
}
