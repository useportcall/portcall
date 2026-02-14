import { CheckoutSession } from "@/types/api";
import OrderTotal from "./order-total";
import { Separator } from "./ui/separator";
import { ChevronRight } from "lucide-react";
import { useTranslation } from "react-i18next";
import { IntervalLabel } from "./interval-label";
import { FreeTrialLabel } from "./free-trial-label";
import { UnitAmount } from "./unit-amount";
import { QuotaLabel } from "./quota-label";

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
  const { t } = useTranslation();
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
        <h2>{t("order_summary.what_you_get")}</h2>
        <div className="gap-8 flex flex-col w-full bg-white rounded-lg p-6 border border-slate-200">
          <div className="flex flex-col gap-2 w-full">
            <p className="font-medium">{t("order_summary.features")}</p>
            <ul>
              {!plan?.features?.length && (
                <li>{t("order_summary.no_features")}</li>
              )}
              {plan?.features?.map((feature) => (
                <li key={feature.id}>
                  <span className="flex gap-1 justify-start items-center">
                    <span className="w-4 h-4">
                      <ChevronRight className="w-4 h-4" />
                    </span>
                    {toTitleCase(feature.id)}
                  </span>
                </li>
              ))}
            </ul>
          </div>
          <Separator />
          <div className="flex flex-col gap-2 w-full text-sm">
            <p className="font-medium">{t("order_summary.metered_features")}</p>
            <div className="flex flex-col gap-4">
              {!session.plan?.metered_features?.length && (
                <p>{t("order_summary.no_metered")}</p>
              )}
              {session.plan?.metered_features?.map((feature) => (
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
