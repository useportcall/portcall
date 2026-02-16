import { AuthenticatedLayout } from "@repo/ui/layout";
import PricingPage from "../../components/pricing-page";

export default async function Page() {
  return (
    <AuthenticatedLayout>
      <PricingPage />
    </AuthenticatedLayout>
  );
}
