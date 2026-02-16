import Features from "@/components/enabled-features";
import { createPortcallClient } from "@/portcall";
import { AuthenticatedLayout } from "@repo/ui/layout";

export default async function Page() {
  return (
    <AuthenticatedLayout>
      <Features />
    </AuthenticatedLayout>
  );
}

const portcall = createPortcallClient({
  apiKey: process.env.PC_API_SECRET!,
});
