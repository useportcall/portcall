import { StatusFilter } from "@/components/table/status-filter";

export function InvoiceListFilters(props: {
  status: string;
  onStatusChange: (value: string) => void;
  t: (key: string) => string;
}) {
  return (
    <div className="flex flex-wrap items-end gap-2">
      <StatusFilter
        label={props.t("views.invoices.table.status_filter_label")}
        value={props.status}
        options={[
          { value: "paid", label: "Paid" },
          { value: "pending", label: "Pending" },
          { value: "draft", label: "Draft" },
          { value: "void", label: "Void" },
        ]}
        allLabel={props.t("views.invoices.table.filter_all")}
        onValueChange={props.onStatusChange}
      />
    </div>
  );
}
