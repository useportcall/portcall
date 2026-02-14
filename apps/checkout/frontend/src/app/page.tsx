"use client";

import { CheckoutForm } from "@/components/checkout-form";
import OrderSummary from "@/components/order-section";
import Link from "next/link";
import { useGetCheckoutSession } from "../hooks/use-get-checkout-session";
import { useTranslation } from "react-i18next";
import { LanguageToggle } from "@/components/language-toggle";

export default function Page() {
  const { t } = useTranslation();
  const {
    data: session,
    isLoading,
    isInvalidLink,
    isSessionError,
    credentials,
    linkError,
    sessionError,
  } = useGetCheckoutSession();

  if (isInvalidLink) {
    return (
      <div className="flex items-center justify-center h-screen">
        <p className="text-lg">{linkError || t("error.invalid_link")}</p>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <p className="text-lg">{t("error.loading")}</p>
      </div>
    );
  }

  if (!credentials) {
    return (
      <div className="flex items-center justify-center h-screen">
        <p className="text-lg">{t("error.invalid_link")}</p>
      </div>
    );
  }

  if (isSessionError || !session) {
    return (
      <div className="flex items-center justify-center h-screen">
        <p className="text-lg">{sessionError || t("error.session_error")}</p>
      </div>
    );
  }

  return (
    <>
      <div className="flex flex-col justify-center items-center px-0 md:px-4 lg:px-0 w-full h-full md:py-20 bg-slate-50 overflow-auto">
        <div className="flex mx-auto px-auto flex-col md:flex-row w-full md:max-w-[1080px] h-full md:h-fit rounded-xl shadow-xs bg-white">
          <div className="bg-slate-50 md:bg-slate-100/80 flex flex-col w-full items-center p-10 lg:p-20 rounded-tl-xl rounded-bl-xl relative">
            <div className="absolute top-4 left-4">
              <LanguageToggle />
            </div>
            <OrderSummary session={session} />
          </div>
          <div className="h-full w-full bg-slate-50 md:bg-white flex flex-col justify-between items-center p-10 lg:p-20 rounded-tr-xl rounded-br-xl animate-fade-in-100">
            <div className="flex flex-col w-full gap-4 md:max-w-sm mx-auto">
              <div className="flex flex-col w-full justify-start items-start gap-2 md:max-w-sm mx-auto">
                <h2 className="text-lg font-medium">{t("payment_details")}</h2>
                <p className="text-xs text-slate-400">
                  {t("payment_description")}
                </p>
              </div>
              <CheckoutForm session={session} credentials={credentials} />
            </div>
            <Link
              href="https://useportcall.com"
              className="mt-4 text-xs text-slate-400 hover:text-slate-600 transition-colors"
            >
              {t("powered_by")}
            </Link>
          </div>
        </div>
      </div>
    </>
  );
}
