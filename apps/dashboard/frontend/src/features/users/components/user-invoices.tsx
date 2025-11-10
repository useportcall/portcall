import EmptyTable from "@/components/empty-table";
import { Badge } from "@/components/ui/badge";
import { useListInvoices } from "@/hooks";
import { useParams } from "react-router-dom";

export function UserInvoices() {
  const { id } = useParams();
  const { data: invoices } = useListInvoices({ userId: id! });

  if (!invoices) {
    return null;
  }

  return (
    <div className="space-y-4 ">
      <h2 className="text-xl font-semibold">User Invoices</h2>
      <div className="space-y-4">
        {!invoices.data.length && (
          <EmptyTable
            message="No invoices issued for this user yet."
            button={""}
          />
        )}
        {invoices.data.map((invoice) => (
          <div key={invoice.id} className="rounded-md border p-4">
            <div className="flex justify-between items-center">
              <div className="flex flex-col items-start">
                <h3 className="font-semibold text-sm">
                  {invoice.invoice_number}
                </h3>
                {!!invoice.total && (
                  <span className="text-sm">
                    {(invoice.total / 100).toLocaleString("en-US", {
                      currency: "USD",
                      style: "currency",
                    })}
                  </span>
                )}
              </div>
              <Badge variant={"outline"}> {invoice.status}</Badge>
            </div>
            <p></p>
          </div>
        ))}
      </div>
    </div>
  );
}
