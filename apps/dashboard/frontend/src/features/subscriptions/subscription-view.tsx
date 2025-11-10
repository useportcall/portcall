import EmptyTable from "@/components/empty-table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { useGetSubscription, useListSubscriptionInvoices } from "@/hooks";
import { cn } from "@/lib/utils";
import { Subscription } from "@/models/subscription";
import { ArrowLeft } from "lucide-react";
import { Link, useNavigate, useParams } from "react-router-dom";

export default function SubscriptionView() {
  const navigate = useNavigate();

  const { id } = useParams();

  const { data: subscription } = useGetSubscription(id!);

  if (!subscription) {
    return <></>;
  }

  return (
    <div className="w-full h-full p-4 lg:p-10 flex flex-col gap-6">
      <div className="flex flex-row justify-between">
        <div className="flex flex-col space-y-8">
          <Button
            variant="ghost"
            className="mb-4 p-0 self-start text-sm flex flex-row justify-start items-center gap-2"
            onClick={() => navigate(-1)}
          >
            <ArrowLeft className="w-4 h-4" /> Back to subscriptions
          </Button>
          <div className="flex flex-col justify-start space-y-2">
            <p className="text-lg md:text-xl font-semibold space-mono-bold">
              Subscription for {subscription.data.user.email}
            </p>
            <p className="text-sm text-slate-400">
              Find more details about this subscription here.
            </p>
          </div>
        </div>
      </div>
      <Separator />
      <div className="grid grid-cols-1 md:grid-cols-[20rem_auto_2fr] w-full h-full">
        <div className="h-full flex flex-col gap-6 px-2">
          <div className="space-y-4">
            <h2 className="text-xl font-semibold">Current Subscription</h2>
            <div className="space-y-1">
              <p className="text-sm text-slate-400">Subscription ID</p>{" "}
              <span className="flex w-full gap-2">
                <Link
                  to={`/subscriptions/${subscription.data.id}`}
                  target="_blank"
                  className="text-xs truncate whitespace-nowrap overflow-ellipsis hover:underline"
                >
                  {subscription.data.id}
                </Link>
              </span>
            </div>
            <div className="space-y-1">
              <p className="text-sm text-slate-400">Related plan</p>{" "}
              <span className="flex w-full gap-2">
                <Link
                  to={`/plans/${subscription.data.plan.id}`}
                  target="_blank"
                  className="text-xs truncate whitespace-nowrap overflow-ellipsis hover:underline"
                >
                  {subscription.data.plan.id}
                </Link>
              </span>
            </div>
            <div className="space-y-1">
              <p className="text-sm text-slate-400">Subscription status</p>{" "}
              <Badge
                className={cn(
                  subscription.data.status === "active"
                    ? "bg-emerald-400/50 text-emerald-800 border-none"
                    : ""
                )}
                variant={"outline"}
              >
                {subscription.data.status}
              </Badge>
            </div>
            <div className="space-y-1">
              <p className="text-sm text-slate-400">Reset date</p>{" "}
              <p className="text-xs">{subscription.data.next_reset_at}</p>
            </div>
          </div>
        </div>
        <Separator orientation="vertical" className="hidden md:block" />
        <div className="h-full px-10 space-y-10">
          <SubscriptionInvoiceList subscription={subscription.data} />
        </div>
      </div>
    </div>
  );
}

function SubscriptionInvoiceList({
  subscription,
}: {
  subscription: Subscription;
}) {
  const { data: invoices } = useListSubscriptionInvoices(subscription.id);

  if (!invoices) {
    return <></>;
  }

  return (
    <div className="space-y-4 ">
      <h2 className="text-xl font-semibold">Invoices for this subscription</h2>
      <div className="space-y-4">
        {!invoices?.data?.length && (
          <EmptyTable
            message="No invoices tied to this subscription issued yet."
            button={""}
          />
        )}
        {/* TODO: investigate whether querying the right thing here */}
        {invoices?.data?.map((invoice: any) => (
          <div key={invoice.id} className="rounded border p-2">
            <div className="flex justify-between items-start">
              <div className="flex flex-col items-start">
                <h3 className="font-semibold text-sm">
                  {invoice.invoice_number}
                </h3>
                <span className="text-sm">
                  {(invoice.total / 100).toLocaleString("en-US", {
                    currency: "USD",
                    style: "currency",
                  })}
                </span>
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
