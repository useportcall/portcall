"use client";

import Link from "next/link";
import { CheckCircle2 } from "lucide-react";
import { useTranslation } from "react-i18next";

export default function SuccessPage() {
  const { t } = useTranslation();

  return (
    <div className="flex items-center justify-center min-h-screen bg-slate-50 px-4">
      <div className="flex flex-col items-center justify-center max-w-md w-full bg-white rounded-xl shadow-lg p-10 md:p-12 text-center space-y-6">
        <div className="flex items-center justify-center w-16 h-16 bg-emerald-100 rounded-full">
          <CheckCircle2 className="w-10 h-10 text-emerald-600" />
        </div>

        <div className="space-y-2">
          <h1
            data-testid="checkout-success-heading"
            className="text-2xl md:text-3xl font-semibold text-slate-900"
          >
            {t("success.title")}
          </h1>
          <p className="text-slate-600 text-sm md:text-base">
            {t("success.message")}
          </p>
        </div>

        <div className="w-full pt-4">
          <div className="bg-slate-50 rounded-lg p-4 space-y-2">
            <p className="text-xs text-slate-500 font-medium uppercase tracking-wide">
              {t("success.whats_next")}
            </p>
            <p className="text-sm text-slate-700">
              {t("success.confirmation_email")}
            </p>
          </div>
        </div>

        <Link
          href="https://useportcall.com"
          className="inline-flex items-center justify-center w-full bg-slate-900 hover:bg-slate-800 text-white font-medium py-3 px-6 rounded-lg transition-colors duration-200"
        >
          {t("success.return_dashboard")}
        </Link>
        <Link
          href="https://useportcall.com"
          className="text-xs text-slate-400 hover:text-slate-600 transition-colors"
        >
          {t("powered_by")}
        </Link>
      </div>
    </div>
  );
}
