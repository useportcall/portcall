import EmptyTable from "@/components/empty-table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useListSubscriptions } from "@/hooks";
import { cn } from "@/lib/utils";
import { Plus } from "lucide-react";
import { useNavigate } from "react-router-dom";

export function SubscriptionListTable() {
  const navigate = useNavigate();
  const { data: subscriptions } = useListSubscriptions();

  if (subscriptions && !subscriptions.data.length) {
    return (
      <EmptyTable message="No subscriptions created yet." button={<></>} />
    );
  }

  return (
    <>
      <div className="flex w-full items-start justify-between">
        <div className="relative flex items-center gap-2 w-full max-w-xs">
          {/* TODO: [MAYBE FOR LATER] Add function to search for subscription using user email */}
          <Input
            type="text"
            placeholder="Find subscription by user email..."
            value=""
            onChange={() => {}}
            className=""
          />
        </div>
        <Button size={"sm"} variant={"outline"} disabled>
          <Plus />
          Create subscription (coming soon!)
        </Button>
      </div>

      <Card className="p-2 animate-fade-in">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="">User</TableHead>
              <TableHead className="">Plan</TableHead>
              <TableHead className="">Status</TableHead>
              <TableHead className="">Created</TableHead>
              <TableHead className="">Next reset</TableHead>
              <TableHead className=""></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody className="">
            {subscriptions.data.map((subscription) => (
              <TableRow
                className=""
                // onClick={() => navigate("/subscriptions/" + subscription.id)}
                key={subscription.id}
              >
                <TableCell className="w-20 py-1">
                  <div>
                    <h4 className="font-semibold overflow-ellipsis">
                      {subscription.user.email}
                    </h4>
                    <span className="text-xs text-slate-500 italic">
                      {subscription.user.id}
                    </span>
                  </div>
                </TableCell>
                <TableCell className="max-w-20 truncate whitespace-nowrap overflow-ellipsis">
                  <div className="flex flex-col">
                    <Badge className="w-fit" variant={"outline"}>
                      {subscription.plan.name}
                    </Badge>
                    {/* <span className="text-xs text-slate-500 italic">
                    {subscription.plan.id}
                  </span> */}
                  </div>
                </TableCell>
                <TableCell className="max-w-20 truncate whitespace-nowrap overflow-ellipsis">
                  <Badge
                    variant={"outline"}
                    className={cn(
                      subscription.status == "active"
                        ? "bg-emerald-400/50 text-emerald-800 border-none"
                        : ""
                    )}
                  >
                    {subscription.status}
                  </Badge>
                </TableCell>
                <TableCell className="max-w-20 truncate whitespace-nowrap overflow-ellipsis">
                  {new Date(subscription.created_at).toLocaleDateString(
                    "en-US",
                    {
                      year: "2-digit",
                      month: "2-digit",
                      day: "2-digit",
                      hour: "2-digit",
                      minute: "2-digit",
                      hour12: false,
                    }
                  )}
                </TableCell>
                <TableCell className="max-w-20 truncate whitespace-nowrap overflow-ellipsis">
                  {subscription.next_reset_at
                    ? new Date(subscription.next_reset_at).toLocaleDateString(
                        "en-US",
                        {
                          year: "2-digit",
                          month: "2-digit",
                          day: "2-digit",
                          hour: "2-digit",
                          minute: "2-digit",
                          hour12: false,
                        }
                      )
                    : ""}
                </TableCell>
                <TableCell align="right" className="right py-1">
                  <Button
                    onClick={() =>
                      navigate(`/subscriptions/${subscription.id}`)
                    }
                    className="text-xs px-3 py-1"
                    variant="outline"
                  >
                    View details
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        {subscriptions?.data.length === 0 && (
          <div className="w-full flex flex-col justify-center items-center p-10 gap-4">
            <p className="text-sm text-slate-400">No subscriptions found.</p>
          </div>
        )}
      </Card>
    </>
  );
}
