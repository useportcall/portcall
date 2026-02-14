import { Address } from "@/models/address";
import { useTranslation } from "react-i18next";

export function BillingAddress({ address }: { address: Address | null }) {
  const { t } = useTranslation();
  if (!address) {
    return (
      <div className="space-y-0.5">
        <p>{t("views.user.fields.billing_address_empty")}</p>
      </div>
    );
  }

  return (
    <div className="space-y-0.5">
      <p>{address.line1}</p>
      <p>{address.city}</p>
      <p>{address.country}</p>
      <p>{address.postal_code}</p>
    </div>
  );
}
