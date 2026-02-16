import { Suspense } from "react";
import { useParams, Link } from "react-router-dom";
import { useApp } from "@/hooks/use-apps";
import { useUsers } from "@/hooks/use-users";
import { useSubscriptions } from "@/hooks/use-subscriptions";
import { usePlans } from "@/hooks/use-plans";
import { useQuotes } from "@/hooks/use-quotes";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { formatDate } from "@/lib/utils";
import {
  ArrowLeft,
  Users,
  CreditCard,
  Boxes,
  ScrollText,
  FileText,
  Link2,
  Sparkles,
} from "lucide-react";
import { Button } from "@/components/ui/button";

function AppDetailContent() {
  const { appId } = useParams<{ appId: string }>();
  const { data: app } = useApp(appId!);
  const { data: users } = useUsers(appId!);
  const { data: subscriptions } = useSubscriptions(appId!);
  const { data: plans } = usePlans(appId!);
  const { data: quotes } = useQuotes(appId!);

  const statCards = [
    { title: "Users", value: app.user_count, icon: Users },
    { title: "Subscriptions", value: app.subscription_count, icon: CreditCard },
    { title: "Plans", value: app.plan_count, icon: Boxes },
    { title: "Quotes", value: app.quote_count, icon: ScrollText },
    { title: "Invoices", value: app.invoice_count, icon: FileText },
    { title: "Features", value: app.feature_count, icon: Sparkles },
    { title: "Connections", value: app.connection_count, icon: Link2 },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Link to="/apps">
          <Button variant="outline" size="icon">
            <ArrowLeft className="h-4 w-4" />
          </Button>
        </Link>
        <div className="flex-1">
          <div className="flex items-center gap-3">
            <h1 className="text-3xl font-bold">{app.name}</h1>
            <Badge variant={app.is_live ? "success" : "secondary"}>
              {app.is_live ? "Live" : "Draft"}
            </Badge>
          </div>
          <p className="text-muted-foreground mt-1">
            ID: {app.id} • Owner: {app.account_email}
          </p>
        </div>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-4 lg:grid-cols-7">
        {statCards.map((stat) => (
          <Card key={stat.title}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-xs font-medium">
                {stat.title}
              </CardTitle>
              <stat.icon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-xl font-bold">{stat.value}</div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Tabs */}
      <Tabs defaultValue="users" className="space-y-4">
        <TabsList>
          <TabsTrigger value="users">Users ({users.length})</TabsTrigger>
          <TabsTrigger value="subscriptions">
            Subscriptions ({subscriptions.length})
          </TabsTrigger>
          <TabsTrigger value="plans">Plans ({plans.length})</TabsTrigger>
          <TabsTrigger value="quotes">Quotes ({quotes.length})</TabsTrigger>
        </TabsList>

        <TabsContent value="users">
          <div className="border rounded-lg">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Email</TableHead>
                  <TableHead>Has Subscription</TableHead>
                  <TableHead>Has Payment</TableHead>
                  <TableHead className="text-right">Entitlements</TableHead>
                  <TableHead>Created</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {users.map((user) => (
                  <TableRow key={user.id}>
                    <TableCell className="font-medium">{user.name}</TableCell>
                    <TableCell>{user.email}</TableCell>
                    <TableCell>
                      <Badge
                        variant={
                          user.has_subscription ? "success" : "secondary"
                        }
                      >
                        {user.has_subscription ? "Yes" : "No"}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant={
                          user.has_payment_method ? "success" : "secondary"
                        }
                      >
                        {user.has_payment_method ? "Yes" : "No"}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-right">
                      {user.entitlement_count}
                    </TableCell>
                    <TableCell>{formatDate(user.created_at)}</TableCell>
                  </TableRow>
                ))}
                {users.length === 0 && (
                  <TableRow>
                    <TableCell
                      colSpan={6}
                      className="text-center text-muted-foreground"
                    >
                      No users found
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>
        </TabsContent>

        <TabsContent value="subscriptions">
          <div className="border rounded-lg">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>User</TableHead>
                  <TableHead>Plan</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Billing</TableHead>
                  <TableHead>Next Reset</TableHead>
                  <TableHead className="text-right">Items</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {subscriptions.map((sub) => (
                  <TableRow key={sub.id}>
                    <TableCell>
                      <div>
                        <div className="font-medium">{sub.user_name}</div>
                        <div className="text-xs text-muted-foreground">
                          {sub.user_email}
                        </div>
                      </div>
                    </TableCell>
                    <TableCell>{sub.plan_name || "-"}</TableCell>
                    <TableCell>
                      <Badge
                        variant={
                          sub.status === "active"
                            ? "success"
                            : sub.status === "canceled"
                              ? "destructive"
                              : "secondary"
                        }
                      >
                        {sub.status}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      {sub.billing_interval} / {sub.currency}
                    </TableCell>
                    <TableCell>{formatDate(sub.next_reset_at)}</TableCell>
                    <TableCell className="text-right">
                      {sub.item_count}
                    </TableCell>
                  </TableRow>
                ))}
                {subscriptions.length === 0 && (
                  <TableRow>
                    <TableCell
                      colSpan={6}
                      className="text-center text-muted-foreground"
                    >
                      No subscriptions found
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>
        </TabsContent>

        <TabsContent value="plans">
          <div className="border rounded-lg">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Pricing</TableHead>
                  <TableHead>Free</TableHead>
                  <TableHead className="text-right">Subscriptions</TableHead>
                  <TableHead className="text-right">Items</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {plans.map((plan) => (
                  <TableRow key={plan.id}>
                    <TableCell className="font-medium">{plan.name}</TableCell>
                    <TableCell>
                      <Badge
                        variant={
                          plan.status === "active"
                            ? "success"
                            : plan.status === "draft"
                              ? "secondary"
                              : "outline"
                        }
                      >
                        {plan.status}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      {plan.interval_count}x {plan.interval} / {plan.currency}
                    </TableCell>
                    <TableCell>
                      <Badge variant={plan.is_free ? "success" : "secondary"}>
                        {plan.is_free ? "Yes" : "No"}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-right">
                      {plan.subscription_count}
                    </TableCell>
                    <TableCell className="text-right">
                      {plan.item_count}
                    </TableCell>
                  </TableRow>
                ))}
                {plans.length === 0 && (
                  <TableRow>
                    <TableCell
                      colSpan={6}
                      className="text-center text-muted-foreground"
                    >
                      No plans found
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>
        </TabsContent>

        <TabsContent value="quotes">
          <div className="border rounded-lg">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Title</TableHead>
                  <TableHead>Recipient</TableHead>
                  <TableHead>Plan</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Valid</TableHead>
                  <TableHead>Created</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {quotes.map((quote) => (
                  <TableRow key={quote.id}>
                    <TableCell className="font-medium">
                      {quote.public_title}
                    </TableCell>
                    <TableCell>
                      <div>
                        <div>{quote.user_name}</div>
                        <div className="text-xs text-muted-foreground">
                          {quote.recipient_email}
                        </div>
                      </div>
                    </TableCell>
                    <TableCell>{quote.plan_name || "-"}</TableCell>
                    <TableCell>
                      <Badge
                        variant={
                          quote.status === "accepted"
                            ? "success"
                            : quote.status === "draft"
                              ? "secondary"
                              : quote.status === "sent"
                                ? "outline"
                                : "destructive"
                        }
                      >
                        {quote.status}
                      </Badge>
                    </TableCell>
                    <TableCell>{quote.days_valid} days</TableCell>
                    <TableCell>{formatDate(quote.created_at)}</TableCell>
                  </TableRow>
                ))}
                {quotes.length === 0 && (
                  <TableRow>
                    <TableCell
                      colSpan={6}
                      className="text-center text-muted-foreground"
                    >
                      No quotes found
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

export default function AppDetailView() {
  return (
    <Suspense
      fallback={
        <div className="flex items-center justify-center h-64">
          <div className="animate-pulse text-muted-foreground">Loading...</div>
        </div>
      }
    >
      <AppDetailContent />
    </Suspense>
  );
}
