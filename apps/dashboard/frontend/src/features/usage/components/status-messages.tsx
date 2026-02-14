import { AlertTriangle, ArrowDownCircle, Zap } from "lucide-react";
import { useTranslation } from "react-i18next";

interface StatusMessagesProps {
    upgradeStatus: string | null;
    downgradeStatus: string | null;
    downgradeScheduled: boolean;
}

/**
 * Displays success/cancelled messages for upgrade/downgrade operations.
 */
export function StatusMessages({
    upgradeStatus,
    downgradeStatus,
    downgradeScheduled,
}: StatusMessagesProps) {
    const { t } = useTranslation();
    if (!upgradeStatus && !downgradeStatus) return null;

    return (
        <>
            {upgradeStatus === "success" && (
                <div className="p-4 bg-green-50 border border-green-200 rounded-lg flex items-center justify-between">
                    <div className="flex items-center gap-2 text-green-800">
                        <Zap className="h-5 w-5" />
                        <span>
                            <strong>{t("views.usage.messages.upgrade_success_title")}</strong> {t("views.usage.messages.upgrade_success_body")}
                        </span>
                    </div>
                    <a
                        href={window.location.pathname}
                        className="text-sm text-green-700 underline hover:text-green-900"
                    >
                        {t("views.usage.messages.refresh")}
                    </a>
                </div>
            )}
            {upgradeStatus === "cancelled" && (
                <div className="p-4 bg-amber-50 border border-amber-200 rounded-lg flex items-center gap-2 text-amber-800">
                    <AlertTriangle className="h-5 w-5" />
                    <span>{t("views.usage.messages.upgrade_cancelled")}</span>
                </div>
            )}
            {downgradeStatus === "success" && (
                <div className="p-4 bg-blue-50 border border-blue-200 rounded-lg flex items-center justify-between">
                    <div className="flex items-center gap-2 text-blue-800">
                        <ArrowDownCircle className="h-5 w-5" />
                        <span>
                            {downgradeScheduled ? (
                                <>
                                    <strong>{t("views.usage.messages.downgrade_scheduled_title")}</strong> {t("views.usage.messages.downgrade_scheduled_body")}
                                </>
                            ) : (
                                <>
                                    <strong>{t("views.usage.messages.downgrade_complete_title")}</strong> {t("views.usage.messages.downgrade_complete_body")}
                                </>
                            )}
                        </span>
                    </div>
                    <a
                        href={window.location.pathname}
                        className="text-sm text-blue-700 underline hover:text-blue-900"
                    >
                        {t("views.usage.messages.refresh")}
                    </a>
                </div>
            )}
        </>
    );
}
