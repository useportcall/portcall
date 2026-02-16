import { Suspense } from "react";
import { useParams, Link, useSearchParams } from "react-router-dom";
import { useDogfoodStatus, useDogfoodUser } from "@/hooks/use-dogfood";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { formatDate } from "@/lib/utils";
import {
  ArrowLeft,
  Loader2,
  User,
  CreditCard,
  Sparkles,
  Receipt,
  Mail,
  Calendar,
  Hash,
} from "lucide-react";

function UserDetailContent() {
  const { userId } = useParams<{ userId: string }>();
  const [searchParams] = useSearchParams();
  const env = searchParams.get("env");
  const { data: status } = useDogfoodStatus();

  const appId = env === "test" ? status?.test_app?.id : status?.live_app?.id;
  const { data: user } = useDogfoodUser(appId, userId);

  if (!user) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-muted-foreground">User not found</p>
      </div>
    );
  }

  const infoItems = [
    { label: "Public ID", value: user.public_id, icon: Hash },
    { label: "Email", value: user.email, icon: Mail },
    { label: "Created", value: formatDate(user.created_at), icon: Calendar },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Link to="/dogfood">
          <Button variant="outline" size="icon">
            <ArrowLeft className="h-4 w-4" />
          </Button>
        </Link>
        <div className="flex-1">
          <div className="flex items-center gap-3">
            <User className="h-6 w-6 text-muted-foreground" />
            <h1 className="text-3xl font-bold">{user.name}</h1>
            {user.subscription && (
              <Badge
                variant={
                  user.subscription.status === "active"
                    ? "success"
                    : user.subscription.status === "canceled"
                      ? "destructive"
                      : "secondary"
                }
              >
                {user.subscription.status}
              </Badge>
            )}
          </div>
          <p className="text-muted-foreground mt-1">{user.email}</p>
        </div>
      </div>

      {/* Info Cards */}
      <div className="grid gap-4 md:grid-cols-3">
        {infoItems.map((item) => (
          <Card key={item.label}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-xs font-medium">
                {item.label}
              </CardTitle>
              <item.icon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-sm font-medium truncate" title={item.value}>
                {item.value}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Subscription Card */}
      {user.subscription && (
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <CreditCard className="h-5 w-5 text-muted-foreground" />
              <CardTitle>Subscription</CardTitle>
            </div>
          </CardHeader>
          <CardContent>
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
              <div>
                <p className="text-sm text-muted-foreground">Plan</p>
                <p className="font-medium">
                  {user.subscription.plan_name || "N/A"}
                </p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Status</p>
                <Badge
                  variant={
                    user.subscription.status === "active"
                      ? "success"
                      : user.subscription.status === "canceled"
                        ? "destructive"
                        : "secondary"
                  }
                >
                  {user.subscription.status}
                </Badge>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Billing</p>
                <p className="font-medium">
                  {user.subscription.billing_interval} (
                  {user.subscription.currency})
                </p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Next Reset</p>
                <p className="font-medium">
                  {formatDate(user.subscription.next_reset_at)}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Tabs */}
      <Tabs defaultValue="entitlements" className="space-y-4">
        <TabsList>
          <TabsTrigger value="entitlements">
            <Sparkles className="mr-2 h-4 w-4" />
            Entitlements ({user.entitlements.length})
          </TabsTrigger>
          <TabsTrigger value="invoices">
            <Receipt className="mr-2 h-4 w-4" />
            Invoices ({user.invoices.length})
          </TabsTrigger>
        </TabsList>

        <TabsContent value="entitlements">
          <div className="border rounded-lg">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Feature</TableHead>
                  <TableHead className="text-right">Usage</TableHead>
                  <TableHead className="text-right">Quota</TableHead>
                  <TableHead>Interval</TableHead>
                  <TableHead>Type</TableHead>
                  <TableHead>Next Reset</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {user.entitlements.map((ent) => {
                  const usagePercent =
                    ent.quota > 0
                      ? Math.round((ent.usage / ent.quota) * 100)
                      : 0;
                  const isExceeded = ent.quota > 0 && ent.usage >= ent.quota;

                  return (
                    <TableRow key={ent.id}>
                      <TableCell className="font-mono font-medium">
                        {ent.feature_public_id}
                      </TableCell>
                      <TableCell className="text-right">
                        <span
                          className={
                            isExceeded ? "text-red-600 font-medium" : ""
                          }
                        >
                          {ent.usage}
                        </span>
                      </TableCell>
                      <TableCell className="text-right">
                        {ent.quota < 0 ? (
                          <span className="text-muted-foreground">
                            Unlimited
                          </span>
                        ) : (
                          <span>
                            {ent.quota}
                            {usagePercent > 0 && (
                              <span className="text-muted-foreground ml-1">
                                ({usagePercent}%)
                              </span>
                            )}
                          </span>
                        )}
                      </TableCell>
                      <TableCell>
                        <Badge variant="outline">{ent.interval}</Badge>
                      </TableCell>
                      <TableCell>
                        <Badge
                          variant={ent.is_metered ? "default" : "secondary"}
                        >
                          {ent.is_metered ? "Metered" : "Boolean"}
                        </Badge>
                      </TableCell>
                      <TableCell>{formatDate(ent.next_reset_at)}</TableCell>
                    </TableRow>
                  );
                })}
                {user.entitlements.length === 0 && (
                  <TableRow>
                    <TableCell
                      colSpan={6}
                      className="text-center text-muted-foreground py-8"
                    >
                      No entitlements found
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>
        </TabsContent>

        <TabsContent value="invoices">
          <div className="border rounded-lg">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Invoice ID</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Total</TableHead>
                  <TableHead>Issued</TableHead>
                  <TableHead>Due</TableHead>
                  <TableHead>Paid</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {user.invoices.map((invoice) => (
                  <TableRow key={invoice.id}>
                    <TableCell className="font-mono text-xs">
                      {invoice.public_id}
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant={
                          invoice.status === "paid"
                            ? "success"
                            : invoice.status === "overdue"
                              ? "destructive"
                              : "secondary"
                        }
                      >
                        {invoice.status}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-right font-medium">
                      {(invoice.total / 100).toFixed(2)} {invoice.currency}
                    </TableCell>
                    <TableCell>{formatDate(invoice.issued_at)}</TableCell>
                    <TableCell>{formatDate(invoice.due_at)}</TableCell>
                    <TableCell>
                      {invoice.paid_at ? formatDate(invoice.paid_at) : "-"}
                    </TableCell>
                  </TableRow>
                ))}
                {user.invoices.length === 0 && (
                  <TableRow>
                    <TableCell
                      colSpan={6}
                      className="text-center text-muted-foreground py-8"
                    >
                      No invoices found
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}

export default function DogfoodUserDetailView() {
  return (
    <Suspense
      fallback={
        <div className="flex items-center justify-center h-64">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      }
    >
      <UserDetailContent />
    </Suspense>
  );
}
