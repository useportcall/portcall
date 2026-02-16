import { AuthenticatedLayout } from "@repo/ui/layout";
import HomePage from "../components/home-page";

export default async function Page() {
  return (
    <AuthenticatedLayout>
      <HomePage />
    </AuthenticatedLayout>
  );
}
