import { SubscriptionListTable } from "./components/subscription-list-table";

export default function SubscriptionListView() {
  return (
    <div className="w-full p-4 lg:p-10 flex flex-col gap-6">
      <div className="flex flex-row justify-between">
        <div className="flex flex-col space-y-8">
          <div className="flex flex-col justify-start space-y-2">
            <p className="text-lg md:text-xl font-semibold">Subscriptions</p>
            <p className="text-sm text-slate-400">
              Review and manage subscriptions here.
            </p>
          </div>
        </div>
      </div>
      <SubscriptionListTable />
    </div>
  );
}
