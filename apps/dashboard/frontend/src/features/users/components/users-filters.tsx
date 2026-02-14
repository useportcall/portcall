import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { UsersTableToolbar } from "./users-table-toolbar";

export type UserFilterOption = "all" | "true" | "false";

export function UsersFilters(props: {
  search: string;
  subscription: UserFilterOption;
  paymentMethod: UserFilterOption;
  onSearchChange: (value: string) => void;
  onSubscriptionChange: (value: UserFilterOption) => void;
  onPaymentMethodChange: (value: UserFilterOption) => void;
  onReset: () => void;
  t: (key: string) => string;
}) {
  return (
    <div className="flex flex-wrap items-end gap-2">
      <UsersTableToolbar
        value={props.search}
        onChange={props.onSearchChange}
        placeholder={props.t("views.users.table.search")}
      />
      <div className="flex flex-wrap items-center gap-2">
        <FilterSelect
          label={props.t("views.users.table.subscription_filter_label")}
          value={props.subscription}
          options={[
            { value: "all", label: props.t("views.users.table.subscription_option_all") },
            { value: "true", label: props.t("views.users.table.subscription_option_true") },
            { value: "false", label: props.t("views.users.table.subscription_option_false") },
          ]}
          onValueChange={(v) => props.onSubscriptionChange(v as UserFilterOption)}
        />
        <FilterSelect
          label={props.t("views.users.table.payment_filter_label")}
          value={props.paymentMethod}
          options={[
            { value: "all", label: props.t("views.users.table.payment_option_all") },
            { value: "true", label: props.t("views.users.table.payment_option_true") },
            { value: "false", label: props.t("views.users.table.payment_option_false") },
          ]}
          onValueChange={(v) => props.onPaymentMethodChange(v as UserFilterOption)}
        />
        <Button variant="outline" size="sm" onClick={props.onReset}>
          {props.t("views.users.table.reset_filters")}
        </Button>
      </div>
    </div>
  );
}

function FilterSelect(props: {
  label: string;
  value: string;
  options: { value: string; label: string }[];
  onValueChange: (value: string) => void;
}) {
  return (
    <div className="space-y-1">
      <Label className="text-xs text-muted-foreground">{props.label}</Label>
      <Select value={props.value} onValueChange={props.onValueChange}>
        <SelectTrigger className="w-[220px] h-9">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          {props.options.map((option) => (
            <SelectItem key={option.value} value={option.value}>
              {option.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  );
}
