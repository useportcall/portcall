import { useTranslation } from "react-i18next";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { CircleFlag } from "react-circle-flags";

export function LanguageSwitcher() {
    const { i18n, t } = useTranslation();
    const currentLanguage = (i18n.resolvedLanguage || i18n.language || "en").split("-")[0];

    return (
        <Select
            value={currentLanguage}
            onValueChange={(value) => i18n.changeLanguage(value)}
        >
            <SelectTrigger
                className="relative h-9 w-[128px] pl-8 text-sm font-medium"
                data-testid="dashboard-language-switcher"
            >
                <span className="pointer-events-none absolute left-2.5 inline-flex h-4 w-4 items-center justify-center overflow-hidden rounded-full">
                    <CircleFlag
                        countryCode={currentLanguage === "ja" ? "jp" : "us"}
                        height={16}
                        width={16}
                    />
                </span>
                <SelectValue placeholder={t("language")}>
                    {currentLanguage === "ja" ? t("japanese") : t("english")}
                </SelectValue>
            </SelectTrigger>
            <SelectContent>
                <SelectItem value="en" className="text-sm">
                    <div className="flex items-center gap-2">
                        <span className="inline-flex h-4 w-4 items-center justify-center overflow-hidden rounded-full">
                            <CircleFlag countryCode="us" height={16} width={16} />
                        </span>
                        <span>{t("english")}</span>
                    </div>
                </SelectItem>
                <SelectItem value="ja" className="text-sm">
                    <div className="flex items-center gap-2">
                        <span className="inline-flex h-4 w-4 items-center justify-center overflow-hidden rounded-full">
                            <CircleFlag countryCode="jp" height={16} width={16} />
                        </span>
                        <span>{t("japanese")}</span>
                    </div>
                </SelectItem>
            </SelectContent>
        </Select>
    );
}
