import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Progress } from "@/components/ui/progress";
import { Separator } from "@/components/ui/separator";
import { Tooltip } from "@/components/ui/tooltip";
import {
  useCancelSubscription,
  useGetUser,
  useGetUserSubscription,
  useUpdateUser,
} from "@/hooks";
import {
  useListBasicEntitlements,
  useListMeteredEntitlements,
} from "@/hooks/api/entitlements";
import { cn } from "@/lib/utils";
import { Subscription } from "@/models/subscription";
import { User } from "@/models/user";
import {
  TooltipArrow,
  TooltipContent,
  TooltipTrigger,
} from "@radix-ui/react-tooltip";
import { ArrowLeft, Edit, Loader2 } from "lucide-react";
import { useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { UserInvoices } from "./components/user-invoices";
import { Address } from "@/models/address";

export default function UserView() {
  const navigate = useNavigate();

  return (
    <div className="w-full h-full p-4 lg:p-10 flex flex-col gap-6">
      <div className="flex flex-row justify-between">
        <div className="flex flex-col space-y-8">
          <Button
            variant="ghost"
            className="mb-4 p-0 self-start text-sm flex flex-row justify-start items-center gap-2"
            onClick={() => navigate(-1)}
          >
            <ArrowLeft className="w-4 h-4" /> Back to users
          </Button>
        </div>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-[20rem_auto_2fr] w-full h-full">
        <div className="h-full flex flex-col gap-6 px-2">
          <UserEmailAndName />
          <Separator />
          <UserSubscription />
        </div>
        <Separator orientation="vertical" className="hidden md:block" />
        <div className="h-full mt-5 px-10 space-y-10">
          <div className="space-y-10">
            <UserEntitlements />
            <UserMeteredFeatures />
          </div>
          <Separator />
          <div className="space-y-10">
            <UserInvoices />
          </div>
        </div>
      </div>
    </div>
  );
}

function UserEmailAndName() {
  const { id } = useParams();
  const { data: user } = useGetUser(id!);

  if (!user) {
    return <></>;
  }

  return (
    <div className="space-y-4">
      <h2 className="text-xl font-semibold">User Information</h2>
      <div className="space-y-1">
        <p className="text-sm text-slate-400">Email</p> {user.data.email}
      </div>
      <div className="space-y-1">
        <p className="text-sm text-slate-400">Name</p>{" "}
        <MutableUserName user={user.data} />
      </div>
      <div className="space-y-1">
        <p className="text-sm text-slate-400">Billing address</p>
        <BillingAddress address={user.data.billing_address} />
      </div>
    </div>
  );
}

function MutableUserName({ user }: { user: User }) {
  const { mutate } = useUpdateUser(user.id);

  return (
    <input
      className="outline-none"
      defaultValue={user.name}
      placeholder="Add name..."
      onBlur={(e) => {
        if (!e.target.value || e.target.value === user.name) return;
        mutate({ name: e.target.value });
      }}
    />
  );
}

function BillingAddress({ address }: { address: Address | null }) {
  if (!address) {
    return (
      <div className="space-y-0.5">
        <p>No billing address added.</p>
      </div>
    );
  }

  return (
    <div className="space-y-0.5">
      <p>{address.line1}</p>
      <p>{address.city}</p>
      <p>{address.country}</p>
      <p>{address.postal_code}</p>
    </div>
  );
}

function UserSubscription() {
  const { id } = useParams();
  const { data: subscription } = useGetUserSubscription(id!);

  if (!subscription?.data) {
    return (
      <div className="space-y-2">
        <h2 className="text-xl font-semibold">Current Subscription</h2>
        <p>-</p>
      </div>
    );
  }

  return (
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
      <div className="space-y-1 ">
        <p className="text-sm text-slate-400">Manage subscription</p>{" "}
        <div className="flex flex-col gap-2">
          <Button disabled variant={"outline"}>
            Apply discount (coming soon!)
          </Button>
          <Button disabled variant={"outline"}>
            Change plan (coming soon!)
          </Button>
          <CancelUserSubscription subscription={subscription.data} />
        </div>
      </div>
    </div>
  );
}

function CancelUserSubscription({
  subscription,
}: {
  subscription: Subscription;
}) {
  const [open, setOpen] = useState(false);
  const { mutateAsync, isPending } = useCancelSubscription(subscription.id);

  const onClick = async () => {
    try {
      await mutateAsync();
      setOpen(false);
    } catch (error) {
      // Handle error
    }
  };

  if (subscription.status !== "active") {
    return null;
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="destructive">Cancel Subscription</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Cancel Subscription</DialogTitle>
          <DialogDescription>
            Are you sure you want to cancel this user's subscription?
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <DialogClose>
            <Button variant="outline">No</Button>
          </DialogClose>
          <Button disabled={isPending} variant="destructive" onClick={onClick}>
            Yes, cancel
            {isPending && <Loader2 className="animate-spin size-3" />}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

function UserEntitlements() {
  const { id } = useParams();
  const { data: entitlements } = useListBasicEntitlements({ userId: id! });

  console.log(entitlements);

  if (!entitlements) {
    return <></>;
  }

  return (
    <div className="space-y-4">
      <div className="flex flex-row justify-between items-center">
        <h2 className="text-xl font-semibold">Entitlements</h2>
        <EditUserEntitlements />
      </div>
      <p className="text-sm text-slate-400">
        See and manage features this user is entitled to based on their
        subscription.
      </p>
      <div className="flex flex-wrap gap-2">
        {!entitlements.data.length && <p>-</p>}
        {entitlements.data.map((entitlement) => (
          <Badge variant={"outline"} key={entitlement.id}>
            {entitlement.id}
          </Badge>
        ))}
      </div>
    </div>
  );
}

function EditUserEntitlements() {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button
          size="icon"
          variant="ghost"
          className="size-[22px] rounded-md hover:bg-slate-100"
        >
          <Edit className="" />
        </Button>
      </DialogTrigger>
      <DialogContent className="lg:min-w-2xl h-96">
        <DialogHeader>
          <DialogTitle>Edit Entitlements</DialogTitle>
          <DialogDescription>
            Modify the entitlements for this user.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <DialogClose>
            <Button variant="outline">Cancel</Button>
          </DialogClose>
          <Button>Save Changes</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

function UserMeteredFeatures() {
  const { id } = useParams();
  const { data: entitlements } = useListMeteredEntitlements({ userId: id! });

  if (!entitlements) {
    return <></>;
  }

  return (
    <div className="space-y-4">
      <h2 className="text-xl font-semibold">Metered Features</h2>
      <p className="text-sm text-slate-400">
        See and manage metered features this user has based on their
        subscription.
      </p>
      <div className="flex flex-col gap-2">
        {!entitlements.data.length && <p>-</p>}
        {entitlements.data.map((entitlement) => (
          <Card key={entitlement.id} className="w-1/3">
            <CardHeader>
              <CardTitle>
                <Badge variant={"outline"}>{entitlement.feature}</Badge>
              </CardTitle>
            </CardHeader>
            <CardContent className="flex flex-col gap-2">
              <Tooltip>
                <TooltipTrigger asChild>
                  <Progress value={entitlement.usage} max={entitlement.quota} />
                </TooltipTrigger>
                <TooltipContent>
                  <TooltipArrow />
                  <p className="bg-slate-800 text-white px-4 text-sm rounded-md">
                    {entitlement.usage}
                  </p>
                </TooltipContent>
              </Tooltip>

              <div className="flex flex-col gap-0">
                <p className="text-sm text-slate-400">
                  Limit: {entitlement.quota}
                </p>
                <p className="text-sm text-slate-400">
                  Resets on: {entitlement.next_reset_at}
                </p>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
