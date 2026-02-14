"use client";

import { useTranslation } from "react-i18next";
import { CircleFlag } from "react-circle-flags";
import { ChevronDown } from "lucide-react";

export function LanguageToggle() {
    const { i18n, t } = useTranslation();
    const currentLanguage = (i18n.resolvedLanguage || i18n.language || "en").split("-")[0];

    return (
        <label className="relative flex h-9 min-w-[128px] items-center rounded-md border border-slate-200 bg-white/95 text-slate-700">
            <span className="pointer-events-none absolute left-2.5 inline-flex h-4 w-4 items-center justify-center overflow-hidden rounded-full">
                <CircleFlag
                    countryCode={currentLanguage === "ja" ? "jp" : "us"}
                    height={16}
                    width={16}
                />
            </span>
            <select
                aria-label={t("language")}
                value={currentLanguage}
                onChange={(e) => i18n.changeLanguage(e.target.value)}
                className="h-full w-full appearance-none bg-transparent py-1 pl-8 pr-8 text-sm font-medium outline-none"
            >
                <option value="en">{t("language_option_english")}</option>
                <option value="ja">{t("language_option_japanese")}</option>
            </select>
            <ChevronDown className="pointer-events-none absolute right-2 h-4 w-4 text-slate-400" />
        </label>
    );
}
