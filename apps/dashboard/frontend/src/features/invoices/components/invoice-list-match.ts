import { Invoice } from "@/models/invoice";

export function matchesInvoice(invoice: Invoice, status: string) {
  return status === "all" || invoice.status.toLowerCase() === status;
}
