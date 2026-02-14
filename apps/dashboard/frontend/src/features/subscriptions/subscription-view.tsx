import { Separator } from "@/components/ui/separator";
import { useGetSubscription } from "@/hooks";
import { useParams } from "react-router-dom";
import { SubscriptionBillingMetersSection } from "./components/view/subscription-billing-meters-section";
import { SubscriptionHeader } from "./components/view/subscription-header";
import { SubscriptionInvoiceList } from "./components/view/subscription-invoice-list";
import { SubscriptionSummary } from "./components/view/subscription-summary";

export default function SubscriptionView() {
  const { id } = useParams();
  const { data: subscription } = useGetSubscription(id!);

  return (
    <div className="w-full h-full p-4 lg:p-10 flex flex-col gap-6">
      <SubscriptionHeader email={subscription.data.user.email} />
      <Separator />
      <div className="grid grid-cols-1 md:grid-cols-[20rem_auto_2fr] w-full h-full">
        <div className="h-full flex flex-col gap-6 px-2">
          <SubscriptionSummary subscription={subscription.data} />
        </div>
        <Separator orientation="vertical" className="hidden md:block" />
        <div className="h-full px-10 space-y-10">
          <SubscriptionBillingMetersSection subscription={subscription.data} />
          <Separator />
          <SubscriptionInvoiceList subscription={subscription.data} />
        </div>
      </div>
    </div>
  );
}
