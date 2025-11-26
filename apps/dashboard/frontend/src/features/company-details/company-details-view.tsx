import { BaseSection, BaseView, LabeledInput } from "@/components/base-view";
import { SaveIndicator } from "@/components/save-indicator";
import { useGetCompany, useUpdateCompany } from "@/hooks";
import { useUpdateAddress } from "@/hooks/api/addresses";
import { cn } from "@/lib/utils";

export default function CompanyDetailsView() {
  const { data: company } = useGetCompany();

  return (
    <BaseView
      title="Company details"
      description="Review company details here."
      actions={<SaveIndicator />}
    >
      <BaseSection title="General information">
        <MutableCompanyName company={company.data} />
        <MutableBusinessAlias company={company.data} />
      </BaseSection>
      <BaseSection title="Registered business address">
        <MutableBillingAddressInput
          prop="line1"
          address={company.data.billing_address}
          placeholder="Address line 1"
        />
        <MutableBillingAddressInput
          prop="line2"
          address={company.data.billing_address}
          placeholder="Address line 2"
        />
        <MutableBillingAddressInput
          prop="city"
          address={company.data.billing_address}
          placeholder="City"
        />
        <MutableBillingAddressInput
          prop="state"
          address={company.data.billing_address}
          placeholder="State"
        />
        <MutableBillingAddressInput
          prop="postal_code"
          address={company.data.billing_address}
          placeholder="Postal code"
        />
        <MutableBillingAddressInput
          prop="country"
          address={company.data.billing_address}
          placeholder="Country"
        />
      </BaseSection>
      <BaseSection title="VAT / sales tax information">
        <MutableVatNumber company={company.data} />
      </BaseSection>
      <BaseSection title="Contact details">
        <MutableFirstName company={company.data} />
        <MutableLastName company={company.data} />
        <MutablePhoneNumber company={company.data} />
      </BaseSection>
    </BaseView>
  );
}

function MutableCompanyName({ company }: { company: any }) {
  const { mutateAsync: updateCompany } = useUpdateCompany();

  return (
    <LabeledInput
      id="company_name"
      label="Registered business name"
      defaultValue={company?.name || ""}
      placeholder="Company name"
      helperText=""
      className={cn(" text-lg outline-none w-96")}
      onBlur={(e) => {
        if (e.target.value.length > 0) {
          updateCompany({ name: e.target.value });
        }
      }}
    />
  );
}

function MutableBusinessAlias({ company }: { company: any }) {
  const { mutateAsync: updateCompany } = useUpdateCompany();

  return (
    <LabeledInput
      id="business_alias"
      helperText=""
      label="Doing business as (DBA)"
      defaultValue={company?.alias || ""}
      placeholder="Optional"
      className={cn("text-lg outline-none w-96")}
      onBlur={(e) => {
        if (e.target.value.length > 0) {
          updateCompany({ alias: e.target.value });
        }
      }}
    />
  );
}

function MutableVatNumber({ company }: { company: any }) {
  const { mutateAsync: updateCompany } = useUpdateCompany();

  return (
    <LabeledInput
      id="vat_number"
      helperText=""
      label="VAT number"
      defaultValue={company?.vat_number || ""}
      placeholder="GB1234567"
      className={cn("text-lg outline-none w-96")}
      onBlur={(e) => {
        if (e.target.value.length > 0) {
          updateCompany({ vat_number: e.target.value });
        }
      }}
    />
  );
}

function MutableFirstName({ company }: { company: any }) {
  const { mutateAsync: updateCompany } = useUpdateCompany();

  return (
    <LabeledInput
      id="first_name"
      helperText=""
      label="First name"
      defaultValue={company?.first_name || ""}
      placeholder="First name"
      className={cn(" text-lg outline-none w-96")}
      onBlur={(e) => {
        if (e.target.value.length > 0) {
          updateCompany({ first_name: e.target.value });
        }
      }}
    />
  );
}

function MutableLastName({ company }: { company: any }) {
  const { mutateAsync: updateCompany } = useUpdateCompany();

  return (
    <LabeledInput
      id="last_name"
      label="Last name"
      helperText=""
      defaultValue={company?.last_name || ""}
      placeholder="Last name"
      className={cn(" text-lg outline-none w-96")}
      onBlur={(e) => {
        if (e.target.value.length > 0) {
          updateCompany({ last_name: e.target.value });
        }
      }}
    />
  );
}
function MutablePhoneNumber({ company }: { company: any }) {
  const { mutateAsync: updateCompany } = useUpdateCompany();

  return (
    <LabeledInput
      id="phone_number"
      helperText=""
      label="Contact phone number"
      defaultValue={company?.phone || ""}
      placeholder="+44 1234 567890"
      className={cn(" text-lg outline-none w-96")}
      onBlur={(e) => {
        if (e.target.value.length > 0) {
          updateCompany({ phone: e.target.value });
        }
      }}
    />
  );
}

function MutableBillingAddressInput({
  prop,
  address,
  placeholder,
}: {
  prop: string;
  address: any;
  placeholder: string;
}) {
  const { mutateAsync } = useUpdateAddress(address.id);

  return (
    <LabeledInput
      id={prop}
      helperText=""
      label={placeholder}
      defaultValue={address[prop] || ""}
      placeholder={placeholder}
      className={cn("text-lg outline-none w-96")}
      onBlur={(e) => {
        if (e.target.value === address[prop]) return;
        mutateAsync({ [prop]: e.target.value });
      }}
    />
  );
}
