import { useListMeteredEntitlements } from "@/hooks/api/entitlements";
import { useTranslation } from "react-i18next";
import { useParams } from "react-router-dom";
import { EditMeteredFeaturesDialog } from "../edit-metered-features-dialog";
import { MeteredEntitlementCard } from "./metered-entitlement-card";

export function UserMeteredFeaturesSection() {
  const { t } = useTranslation();
  const { id } = useParams();
  const { data: entitlements } = useListMeteredEntitlements({ userId: id! });

  if (!entitlements) return null;

  return (
    <div className="space-y-2">
      <div className="flex items-center justify-between">
        <h2 className="text-sm text-muted-foreground">{t("views.user.metered_features.title")}</h2>
        <EditMeteredFeaturesDialog userId={id!} />
      </div>
      <p className="text-xs text-muted-foreground">{t("views.user.metered_features.description")}</p>
      <div className="grid gap-3 md:grid-cols-2 lg:grid-cols-3">
        {!entitlements.data.length && (
          <p className="text-sm text-muted-foreground col-span-full">{t("views.user.metered_features.empty")}</p>
        )}
        {entitlements.data.map((entitlement) => (
          <MeteredEntitlementCard key={entitlement.id} entitlement={entitlement} />
        ))}
      </div>
    </div>
  );
}
