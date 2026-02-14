import { BaseSection, BaseView } from "@/components/base-view";
import { useTranslation } from "react-i18next";
import { ConnectionsGrid } from "./components/connections-grid";
import { CreateConnectionDialog } from "./components/create-connection-dialog";

export default function PaymentIntegrationsView() {
  const { t } = useTranslation();
  return (
    <BaseView
      title={t("views.integrations.title")}
      description={t("views.integrations.description")}
    >
      <BaseSection title={t("views.integrations.connected_providers")} actions={<CreateConnectionDialog />}>
        <ConnectionsGrid />
      </BaseSection>
    </BaseView>
  );
}
