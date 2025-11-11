import { CheckoutSession, Plan } from "@/types/api";
import { Lock } from "lucide-react";
import { useMemo } from "react";

export function Totals({ session }: { session: CheckoutSession }) {
  const totals = useMemo(() => {
    if (!session.plan)
      return { subtotal: "$0.00", tax: "$0.00", total: "$0.00" };
    return calculateFormattedTotals(session.plan);
  }, [session.plan]);

  return (
    <>
      <div className="flex flex-col justify-start items-start gap-4 mt-4">
        <span className="w-full flex text-sm justify-between">
          <p>Subtotal</p> <p>{totals.subtotal}</p>
        </span>
        <span className="w-full text-sm flex justify-between">
          <p>Taxes</p> <p>{totals.tax}</p>
        </span>
        <span className="w-full font-semibold flex justify-between">
          <p>Total</p> <p>{totals.total}</p>
        </span>
      </div>
      <button className="w-full bg-slate-900 p-2 flex justify-center items-center gap-2 text-white font-semibold mt-2 rounded-md hover:bg-slate-700 cursor-pointer">
        Pay {totals.total}
        <span className="w-4 h-4">
          <Lock className="w-4 h-4" />
        </span>
      </button>
    </>
  );
}

function calculateFormattedTotals(plan: Plan) {
  const mainPlanItem = plan.items.find(
    (item) => item.pricing_model === "fixed"
  );

  const subtotal = mainPlanItem?.unit_amount || 0;
  const tax = 0; // TODO: Calculate tax based on billing address
  const total = subtotal + tax;

  const formatAmount = (amount: number, currency: string) => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: currency,
    }).format(amount / 100);
  };

  return {
    subtotal: formatAmount(subtotal, plan.currency),
    tax: formatAmount(tax, plan.currency),
    total: formatAmount(total, plan.currency),
  };
}
