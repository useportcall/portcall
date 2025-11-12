import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  CLAY_EXPLORER_PLAN_ID,
  CLAY_FREE_PLAN_ID,
  CLAY_PRO_PLAN_ID,
  CLAY_STARTER_PLAN_ID,
} from "@/constants";
import { getUserSubscription } from "@/lib/api";
import { cn } from "@/lib/utils";
import { ArrowRight, Check, Dot, X } from "lucide-react";
import { SubscribeButton } from "./subscripe-button";
import { UnsubscribeButton } from "./unsubscribe-button";

const features = [
  {
    feature: "Sculptor",
    items: [true, true, true, true],
  },
  {
    feature: "Sequencer",
    items: [true, true, true, true],
  },
  {
    feature: "Exporting",
    items: [true, true, true, true],
  },
  {
    feature: "AI / Claygent",
    items: [true, true, true, true],
  },
  {
    feature: "Rollover credits",
    items: [true, true, true, true],
  },
  {
    feature: "100+ integration providers",
    items: [true, true, true, true],
  },
  {
    feature: "Phone number enrichments",
    items: [false, true, true, true],
  },
  {
    feature: "Use your own API key",
    items: [false, true, true, true],
  },
  {
    feature: "Integrate with any HTTP API",
    items: [false, false, true, true],
  },
  {
    feature: "Webhooks",
    items: [false, false, true, true],
  },
  {
    feature: "Email sequencing integrations",
    items: [false, false, true, true],
  },
  {
    feature: "CRM integrations",
    items: [false, false, false, true],
  },
];

export function ClayPricingTable() {
  return (
    <Table>
      <TableHeader className="[&_tr]:border-b-0">
        <TableRow className="border-0">
          <TableHead className="bg-[#f3f2ed]"></TableHead>
          <TableHead className="bg-[#f3f2ed] text-center">
            <div className="flex flex-col gap-4 py-4 justify-center items-center">
              <span>Free</span>
              <span>$0 / mo</span>
              <Select>
                <SelectTrigger className="bg-white text-slate-800">
                  <SelectValue
                    defaultValue={"100"}
                    placeholder="100 credits/month"
                  />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="100">100 credits/month</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </TableHead>
          <TableHead className="bg-[#d8d0ff] text-center">
            <div className="flex flex-col gap-4 py-4 justify-center items-center">
              <Badge className="bg-[#7934f0] text-white rounded-full">
                Starter
              </Badge>
              <span>$229 / mo</span>
              <Select>
                <SelectTrigger className="bg-white text-slate-800">
                  <SelectValue
                    defaultValue={"2000"}
                    placeholder="2,000 credits/month"
                  />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="2000">2,000 credits/month</SelectItem>
                  <SelectItem value="3000">3,000 credits/month</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </TableHead>
          <TableHead className="bg-[#fef0a4] text-center">
            <div className="flex flex-col gap-4 py-4 justify-center items-center">
              <Badge className="bg-[#fdad15] text-white rounded-full">
                Explorer
              </Badge>
              <span>$349 / mo</span>
              <Select>
                <SelectTrigger className="bg-white text-slate-800">
                  <SelectValue
                    defaultValue={"1000"}
                    placeholder="10k credits/month"
                  />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="10000">10k credits/month</SelectItem>
                  <SelectItem value="14000">14k credits/month</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </TableHead>
          <TableHead className="bg-[#ffd1f1] text-center py-4">
            <div className="flex flex-col gap-4 justify-center items-center">
              <div className="flex flex-col gap-4 py-4 justify-center items-center">
                <Badge className="bg-[#cc0687] text-white rounded-full">
                  Pro
                </Badge>
                <span>$800 / mo</span>
                <Select>
                  <SelectTrigger className="bg-white text-slate-800">
                    <SelectValue
                      defaultValue={"50000"}
                      placeholder="50k credits/month"
                    />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="50000">50k credits/month</SelectItem>
                    <SelectItem value="70000">70k credits/month</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
          </TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow className="border-0">
          <TableCell>Users</TableCell>
          <TableCell className="text-center">Unlimited</TableCell>
          <TableCell className="bg-[#f5f3ff] text-center">Unlimited</TableCell>
          <TableCell className="text-center bg-[#fff8d2]">Unlimited</TableCell>
          <TableCell className="text-center bg-[#ffe5f7]">Unlimited</TableCell>
        </TableRow>
        <TableRow className="border-0">
          <TableCell>People/Company searches</TableCell>
          <TableCell className="text-center">Up to 100/search</TableCell>
          <TableCell className="bg-[#f5f3ff] text-center">
            Up to 5,000/search
          </TableCell>
          <TableCell className="text-center bg-[#fff8d2]">
            Up to 10,000/search
          </TableCell>
          <TableCell className="text-center bg-[#ffe5f7]">
            Up to 25,000/search
          </TableCell>
        </TableRow>
        {features.map((feature) => (
          <FeatureRow
            key={feature.feature}
            feature={feature.feature}
            items={feature.items}
          />
        ))}
        <TableRow className="border-0">
          <TableCell></TableCell>
          <TableCell className="text-center">
            <PlanSubmitButton planId={CLAY_FREE_PLAN_ID} />
          </TableCell>
          <TableCell className="text-center">
            <PlanSubmitButton planId={CLAY_STARTER_PLAN_ID} />
          </TableCell>
          <TableCell className="text-center">
            <PlanSubmitButton planId={CLAY_EXPLORER_PLAN_ID} />
          </TableCell>
          <TableCell className="text-center">
            <PlanSubmitButton planId={CLAY_PRO_PLAN_ID} />
          </TableCell>
        </TableRow>
      </TableBody>
    </Table>
  );
}

function FeatureRow({ feature, items }: { feature: string; items: boolean[] }) {
  return (
    <TableRow className="border-0">
      <TableCell>{feature}</TableCell>
      {items.map((item, index) => (
        <TableCell
          key={index}
          className={cn("text-center", {
            "bg-[#f5f3ff]": index === 1,
            "bg-[#fff8d2]": index === 2,
            "bg-[#ffe5f7]": index === 3,
          })}
        >
          <div className="w-full flex justify-center">
            {item ? (
              <div
                className={cn(
                  "bg-black rounded-full size-4 flex justify-center items-center"
                )}
              >
                <Check className="text-white size-2" />
              </div>
            ) : (
              <div className="size-4 flex justify-center items-center">
                <Dot className="size-2" />
              </div>
            )}
          </div>
        </TableCell>
      ))}
    </TableRow>
  );
}

async function PlanSubmitButton({ planId }: { planId: string }) {
  const subscription = await getUserSubscription();

  if (!subscription) {
    return <Subscribe planId={planId} />;
  }

  return <Unsubscribe id={subscription.id} />;
}

function Subscribe({ planId }: { planId: string }) {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button type="submit">
          Subscribe <ArrowRight />
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogTitle>Subscribe</DialogTitle>
        <DialogDescription>
          You will be redirected to the checkout page.
        </DialogDescription>
        <DialogFooter>
          <SubscribeButton planId={planId} />
          <DialogClose asChild>
            <Button variant="outline">No</Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

function Unsubscribe({ id }: { id: string }) {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button type="submit">
          Unsubscribe <X />
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogTitle>Subscribe</DialogTitle>
        <DialogDescription>
          Are you sure you want to unsubscribe?
        </DialogDescription>
        <DialogFooter>
          <UnsubscribeButton subscriptionId={id} />
          <DialogClose asChild>
            <Button variant="outline">No</Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

// function UpdateButton({
//   subscriptionId,
//   planId,
// }: {
//   subscriptionId: string;
//   planId: string;
// }) {
//   const { mutateAsync, isPending } = useUpdateSubscription(subscriptionId);

//   return (
//     <ConfirmationDialog
//       isPending={isPending}
//       title="Switch Plan"
//       description="Are you sure you want to switch plan?"
//     >
//       <Button>
//         Switch <Repeat />
//       </Button>
//     </ConfirmationDialog>
//   );
// }
