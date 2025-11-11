import { CheckoutSession } from "@/types/api";
import { useMemo } from "react";

export default function OrderTotal({ session }: { session: CheckoutSession }) {
  const total = useMemo(() => {
    const subtotal = calculateSubtotal(session);
    const tax = subtotal * 0;
    return subtotal + tax;
  }, [session]);

  return (
    <div className="">
      <p className="text-slate-800">Pay ${(total / 100).toFixed(2)} today</p>
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
