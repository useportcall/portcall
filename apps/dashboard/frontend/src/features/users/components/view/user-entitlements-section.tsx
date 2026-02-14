import { Badge } from "@/components/ui/badge";
import { useListBasicEntitlements } from "@/hooks/api/entitlements";
import { useTranslation } from "react-i18next";
import { useParams } from "react-router-dom";
import { EditBasicEntitlementsDialog } from "../edit-basic-entitlements-dialog";

export function UserEntitlementsSection() {
  const { t } = useTranslation();
  const { id } = useParams();
  const { data: entitlements } = useListBasicEntitlements({ userId: id! });

  if (!entitlements) return null;

  return (
    <div className="space-y-2">
      <div className="flex items-center justify-between">
        <h2 className="text-sm text-muted-foreground">{t("views.user.entitlements.title")}</h2>
        <EditBasicEntitlementsDialog userId={id!} />
      </div>
      <p className="text-xs text-muted-foreground">
        {t("views.user.entitlements.description")}
      </p>
      <div className="flex flex-wrap gap-2">
        {!entitlements.data.length && <p className="text-sm text-muted-foreground">{t("views.user.entitlements.empty")}</p>}
        {entitlements.data.map((entitlement) => (
          <Badge variant="outline" key={entitlement.id}>{entitlement.id}</Badge>
        ))}
      </div>
    </div>
  );
}
