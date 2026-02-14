import { Button } from "@/components/ui/button";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog";
import { Zap } from "lucide-react";
import { useTranslation } from "react-i18next";

interface UpgradeDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    onConfirm: () => void;
    isPending: boolean;
    isDisabled: boolean;
}

export function UpgradeDialog({
    open,
    onOpenChange,
    onConfirm,
    isPending,
    isDisabled,
}: UpgradeDialogProps) {
    const { t } = useTranslation();
    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogTrigger asChild>
                <Button disabled={isDisabled}>
                    {isPending ? (
                        <span className="animate-pulse">{t("common.processing")}</span>
                    ) : (
                        <>
                            <Zap className="mr-2 h-4 w-4" />
                            {t("views.usage.current_plan.upgrade_action")}
                        </>
                    )}
                </Button>
            </DialogTrigger>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>{t("views.usage.upgrade_dialog.title")}</DialogTitle>
                    <DialogDescription>
                        {t("views.usage.upgrade_dialog.description")}
                    </DialogDescription>
                </DialogHeader>
                <div className="py-4 space-y-2 text-sm">
                    <p className="text-muted-foreground">{t("views.usage.upgrade_dialog.after_title")}</p>
                    <ul className="list-disc list-inside text-muted-foreground space-y-1">
                        <li>{t("views.usage.upgrade_dialog.after_users")}</li>
                        <li>{t("views.usage.upgrade_dialog.after_subscriptions")}</li>
                        <li>{t("views.usage.upgrade_dialog.after_support")}</li>
                        <li>{t("views.usage.upgrade_dialog.after_price")}</li>
                    </ul>
                </div>
                <DialogFooter>
                    <Button variant="outline" onClick={() => onOpenChange(false)}>
                        {t("common.cancel")}
                    </Button>
                    <Button onClick={onConfirm} disabled={isPending}>
                        {isPending ? t("common.processing") : t("views.usage.upgrade_dialog.confirm")}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}
