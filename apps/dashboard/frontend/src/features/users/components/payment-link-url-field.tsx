import { Button } from "@/components/ui/button";
import { Copy } from "lucide-react";
import { useTranslation } from "react-i18next";

export function PaymentLinkURLField({ linkURL }: { linkURL: string }) {
  const { t } = useTranslation();
  return (
    <div className="space-y-2">
      <p className="text-sm text-muted-foreground">{t("views.user.payment_link.share_url")}</p>
      <div className="flex gap-2">
        <input
          data-testid="payment-link-url-input"
          readOnly
          value={linkURL}
          className="w-full h-10 px-3 rounded-lg border border-input bg-muted text-xs"
        />
        <Button
          variant="outline"
          size="icon"
          onClick={() => navigator.clipboard.writeText(linkURL)}
        >
          <Copy className="size-4" />
        </Button>
      </div>
    </div>
  );
}
