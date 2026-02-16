import { ClayPricingTable } from "../../components/clay-pricing-table";
import { AuthenticatedLayout } from "@repo/ui/layout";

export default async function Page() {
  return (
    <AuthenticatedLayout>
      <ClayPricingTable />
    </AuthenticatedLayout>
  );
}
