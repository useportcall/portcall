import { InvoiceListTable } from "./components/invoice-list-table";
import { useTranslation } from "react-i18next";

export default function InvoiceListView() {
  const { t } = useTranslation();
  return (
    <div className="w-full p-4 lg:p-10 flex flex-col gap-6">
      <div className="flex flex-row justify-between">
        <div className="flex flex-col space-y-8">
          <div className="flex flex-col justify-start space-y-2">
            <p className="text-lg md:text-xl font-semibold">{t("views.invoices.title")}</p>
            <p className="text-sm text-muted-foreground">
              {t("views.invoices.description")}
            </p>
          </div>
        </div>
      </div>
      <InvoiceListTable />
    </div>
  );
}
