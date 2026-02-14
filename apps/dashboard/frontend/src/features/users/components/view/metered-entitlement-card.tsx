import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { formatShortDate } from "@/features/billing-meters/utils/format";
import { useTranslation } from "react-i18next";

export function MeteredEntitlementCard({ entitlement }: { entitlement: any }) {
  const { t } = useTranslation();
  const usagePercentage = entitlement.quota > 0
    ? Math.min(100, (entitlement.usage / entitlement.quota) * 100)
    : 0;

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm flex items-center justify-between">
          <Badge variant="outline">{entitlement.id}</Badge>
          <span className="text-xs font-normal text-muted-foreground">
            {entitlement.usage} / {entitlement.quota === -1 ? "∞" : entitlement.quota}
          </span>
        </CardTitle>
      </CardHeader>
      <CardContent className="flex flex-col gap-2">
        <Tooltip>
          <TooltipTrigger asChild>
            <div className="h-2 bg-muted rounded-full overflow-hidden">
              <div className="h-full bg-primary transition-all" style={{ width: `${usagePercentage}%` }} />
            </div>
          </TooltipTrigger>
          <TooltipContent>
            <p className="bg-primary text-primary-foreground px-4 text-sm rounded-md">
              {t("views.user.metered_features.used", { count: entitlement.usage })}
            </p>
          </TooltipContent>
        </Tooltip>
        {entitlement.next_reset_at && (
          <p className="text-xs text-muted-foreground">
            {t("views.user.metered_features.resets", {
              date: formatShortDate(entitlement.next_reset_at),
            })}
          </p>
        )}
      </CardContent>
    </Card>
  );
}
