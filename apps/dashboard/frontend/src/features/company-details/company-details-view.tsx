import { BaseSection, BaseView } from "@/components/base-view";
import { useGetCompany } from "@/hooks";
import { MutableBillingAddressInput } from "./components/mutable-billing-address-input";
import { MutableCompanyField } from "./components/mutable-company-field";
import { CompanyLogoInput } from "./components/company-logo-input";
import { useTranslation } from "react-i18next";

export default function CompanyDetailsView() {
  const { t } = useTranslation();
  const { data: company } = useGetCompany();
  const c = company.data;

  return (
    <BaseView title={t("views.company.title")} description={t("views.company.description")}>
      <BaseSection title={t("views.company.sections.branding")}><CompanyLogoInput /></BaseSection>
      <BaseSection title={t("views.company.sections.general")}><MutableCompanyField id="company_name" label={t("views.company.fields.registered_name")} defaultValue={c.name || ""} placeholder={t("views.company.placeholders.company_name")} payloadKey="name" /><MutableCompanyField id="business_alias" label={t("views.company.fields.dba")} defaultValue={c.alias || ""} placeholder={t("views.company.placeholders.optional")} payloadKey="alias" /></BaseSection>
      <BaseSection title={t("views.company.sections.address")}>{(["line1", "line2", "city", "state", "postal_code", "country"] as const).map((prop) => <MutableBillingAddressInput key={prop} prop={prop} address={c.billing_address} placeholder={placeholderForAddress(prop, t)} />)}</BaseSection>
      <BaseSection title={t("views.company.sections.tax")}><MutableCompanyField id="vat_number" label={t("views.company.fields.vat")} defaultValue={c.vat_number || ""} placeholder="GB1234567" payloadKey="vat_number" /></BaseSection>
      <BaseSection title={t("views.company.sections.contact")}><MutableCompanyField id="first_name" label={t("views.company.fields.first_name")} defaultValue={c.first_name || ""} placeholder={t("views.company.fields.first_name")} payloadKey="first_name" /><MutableCompanyField id="last_name" label={t("views.company.fields.last_name")} defaultValue={c.last_name || ""} placeholder={t("views.company.fields.last_name")} payloadKey="last_name" /><MutableCompanyField id="phone_number" label={t("views.company.fields.phone")} defaultValue={c.phone || ""} placeholder="+44 1234 567890" payloadKey="phone_number" /></BaseSection>
    </BaseView>
  );
}

function placeholderForAddress(prop: "line1" | "line2" | "city" | "state" | "postal_code" | "country", t: (k: string) => string) {
  switch (prop) {
    case "line1": return t("views.company.address.line1");
    case "line2": return t("views.company.address.line2");
    case "city": return t("views.company.address.city");
    case "state": return t("views.company.address.state");
    case "postal_code": return t("views.company.address.postal_code");
    case "country": return t("views.company.address.country");
  }
}
