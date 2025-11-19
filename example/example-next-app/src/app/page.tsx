import { Button } from "@/components/ui/button";
import { getUserSubscription } from "@/lib/api";
import { logout } from "./actions/logout";
import { ClayPricingTable } from "../components/clay-pricing-table";
import Features from "../components/enabled-features";

export default async function Page() {
  const subscription = await getUserSubscription();

  return (
    <div className="grid grid-cols-1 sm:grid-cols-[8rem_1fr] h-screen w-screen">
      <div className="w-full h-screen border-r hidden sm:block"></div>
      <div className="flex flex-col h-screen">
        <div className="h-15 flex items-center justify-end gap-2 p-4 border-b">
          <LogoutButton />
        </div>
        <div className="overflow-auto grid grid-cols-1  md:grid-cols-[3fr_1fr] h-full">
          <div className="border-r overflow-auto">
            <div className="text-xs  h-full">
              <NextReset subscription={subscription} />
              <ClayPricingTable />
            </div>
          </div>
          <Features />
        </div>
      </div>
    </div>
  );
}

function LogoutButton() {
  return (
    <form action={logout}>
      <Button type="submit">Logout</Button>
    </form>
  );
}

function NextReset(props: { subscription: any }) {
  if (!props.subscription) {
    return null;
  }

  return (
    <div className="p-4 flex gap-2 items-center border-b">
      <h4 className="text-sm font-semibold">Next reset:</h4>
      <span>{props.subscription.next_reset_at}</span>
    </div>
  );
}
