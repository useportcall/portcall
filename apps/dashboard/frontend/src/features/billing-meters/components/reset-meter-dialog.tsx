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
import { Loader2, RotateCcw } from "lucide-react";
import { useTranslation } from "react-i18next";

export function ResetMeterDialog({
  disabled,
  isResetting,
  onReset,
}: {
  disabled: boolean;
  isResetting: boolean;
  onReset: () => void;
}) {
  const { t } = useTranslation();
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button
          size="sm"
          variant="ghost"
          className="h-8 w-8 p-0 text-muted-foreground hover:text-destructive"
          disabled={disabled}
        >
          <RotateCcw className="h-3 w-3" />
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t("views.user.billing_meters.reset_title")}</DialogTitle>
          <DialogDescription>{t("views.user.billing_meters.reset_description")}</DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="destructive" onClick={onReset} disabled={isResetting}>
            {isResetting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            {t("views.user.billing_meters.reset_confirm")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
