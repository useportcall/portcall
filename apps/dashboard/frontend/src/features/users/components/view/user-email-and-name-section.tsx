import { useGetUser } from "@/hooks";
import { useTranslation } from "react-i18next";
import { useParams } from "react-router-dom";
import { BillingAddress } from "./billing-address";
import { MutableUserName } from "./mutable-user-name";

export function UserEmailAndNameSection() {
  const { t } = useTranslation();
  const { id } = useParams();
  const { data: user } = useGetUser(id!);

  if (!user) return null;

  return (
    <div className="space-y-4">
      <h2 className="text-xl font-semibold">{t("views.user.sections.info")}</h2>
      <div className="space-y-1">
        <p className="text-sm text-muted-foreground">{t("views.user.fields.email")}</p>
        {user.data.email}
      </div>
      <div className="space-y-1">
        <p className="text-sm text-muted-foreground">{t("views.user.fields.name")}</p>
        <MutableUserName user={user.data} />
      </div>
      <div className="space-y-1">
        <p className="text-sm text-muted-foreground">{t("views.user.fields.billing_address")}</p>
        <BillingAddress address={user.data.billing_address} />
      </div>
    </div>
  );
}
