import {
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { useTranslation } from "react-i18next";

export function PaymentLinkDialogHeader() {
  const { t } = useTranslation();
  return (
    <DialogHeader>
      <DialogTitle>{t("views.user.payment_link.title")}</DialogTitle>
      <DialogDescription>{t("views.user.payment_link.description")}</DialogDescription>
    </DialogHeader>
  );
}
