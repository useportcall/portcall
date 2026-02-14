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
import { ArrowDownCircle } from "lucide-react";
import { useTranslation } from "react-i18next";

interface DowngradeDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    onConfirm: () => void;
    isPending: boolean;
}

export function DowngradeDialog({
    open,
    onOpenChange,
    onConfirm,
    isPending,
}: DowngradeDialogProps) {
    const { t } = useTranslation();
    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogTrigger asChild>
                <Button variant="outline" size="sm">
                    <ArrowDownCircle className="mr-2 h-4 w-4" />
                    {t("views.usage.current_plan.downgrade_action")}
                </Button>
            </DialogTrigger>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>{t("views.usage.downgrade_dialog.title")}</DialogTitle>
                    <DialogDescription>
                        {t("views.usage.downgrade_dialog.description")}
                    </DialogDescription>
                </DialogHeader>
                <div className="py-4 space-y-2 text-sm">
                    <p className="text-muted-foreground">{t("views.usage.downgrade_dialog.after_title")}</p>
                    <ul className="list-disc list-inside text-muted-foreground space-y-1">
                        <li>{t("views.usage.downgrade_dialog.after_users")}</li>
                        <li>{t("views.usage.downgrade_dialog.after_subscriptions")}</li>
                        <li>{t("views.usage.downgrade_dialog.after_data")}</li>
                        <li>{t("views.usage.downgrade_dialog.after_upgrade_again")}</li>
                    </ul>
                </div>
                <DialogFooter>
                    <Button variant="outline" onClick={() => onOpenChange(false)}>
                        {t("common.cancel")}
                    </Button>
                    <Button
                        variant="destructive"
                        onClick={onConfirm}
                        disabled={isPending}
                    >
                        {isPending ? t("common.processing") : t("views.usage.downgrade_dialog.confirm")}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}
