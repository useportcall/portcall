"use client";

import { CheckoutForm } from "@/components/checkout-form";
import OrderSummary from "@/components/order-section";
import Link from "next/link";
import { useGetCheckoutSession } from "../hooks/use-get-checkout-session";

export default function Page() {
  const { data: session } = useGetCheckoutSession();

  if (!session) {
    return (
      <div className="flex items-center justify-center h-screen">
        <p className="text-lg">Loading...</p>
      </div>
    );
  }

  return (
    <>
      <div className="flex flex-col justify-center items-center px-0 md:px-4 lg:px-0 w-full h-full md:py-20 bg-slate-50 overflow-auto">
        <div className="flex mx-auto px-auto flex-col md:flex-row w-full md:max-w-[1080px] h-full md:h-fit rounded-xl shadow-xs bg-white">
          <div className="bg-slate-50 md:bg-slate-100/80 flex flex-col w-full items-center p-10 lg:p-20 rounded-tl-xl rounded-bl-xl">
            <OrderSummary session={session} />
          </div>
          <div className="h-full w-full bg-slate-50 md:bg-white flex flex-col justify-between items-center p-10 lg:p-20 rounded-tr-xl rounded-br-xl animate-fade-in-100">
            <div className="flex flex-col w-full gap-4 md:max-w-sm mx-auto">
              <div className="flex flex-col w-full justify-start items-start gap-2 md:max-w-sm mx-auto">
                <h2 className="text-lg font-medium">Payment details</h2>
                <p className="text-xs text-slate-400">
                  Complete your subscription by providing your payment details
                </p>
              </div>
              <CheckoutForm session={session} />
            </div>
            <Link
              href={"https://useportcall.com/ref?" + `${session.id}`}
              className="bg-slate-100 mt-4 rounded-full text-slate-800 px-4 py-2 shadow-xs"
            >
              Powered by ⚓️ Portcall
            </Link>
          </div>
        </div>
      </div>
    </>
  );
}
