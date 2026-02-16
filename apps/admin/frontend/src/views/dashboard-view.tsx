import { Suspense } from "react";
import { useStats } from "@/hooks/use-stats";
import { useApps } from "@/hooks/use-apps";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  AppWindow,
  Users,
  CreditCard,
  FileText,
  ScrollText,
  Boxes,
  Activity,
  TrendingUp,
} from "lucide-react";
import { Link } from "react-router-dom";

function DashboardContent() {
  const { data: stats } = useStats();
  const { data: apps } = useApps();

  const statCards = [
    {
      title: "Total Apps",
      value: stats.total_apps,
      icon: AppWindow,
      color: "text-blue-600",
    },
    {
      title: "Total Users",
      value: stats.total_users,
      icon: Users,
      color: "text-green-600",
    },
    {
      title: "Subscriptions",
      value: stats.total_subscriptions,
      icon: CreditCard,
      color: "text-purple-600",
    },
    {
      title: "Active Subscriptions",
      value: stats.active_subscriptions,
      icon: Activity,
      color: "text-emerald-600",
    },
    {
      title: "Total Plans",
      value: stats.total_plans,
      icon: Boxes,
      color: "text-orange-600",
    },
    {
      title: "Quotes",
      value: stats.total_quotes,
      icon: ScrollText,
      color: "text-cyan-600",
    },
    {
      title: "Paid Invoices",
      value: stats.paid_invoices,
      icon: TrendingUp,
      color: "text-green-700",
    },
    {
      title: "Pending Invoices",
      value: stats.pending_invoices,
      icon: FileText,
      color: "text-yellow-600",
    },
  ];

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <p className="text-muted-foreground mt-1">
          Overview of all Portcall data across all apps
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {statCards.map((stat) => (
          <Card key={stat.title}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                {stat.title}
              </CardTitle>
              <stat.icon className={`h-4 w-4 ${stat.color}`} />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {stat.value.toLocaleString()}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Recent Apps */}
      <div>
        <h2 className="text-xl font-semibold mb-4">Recent Apps</h2>
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {apps.slice(0, 6).map((app) => (
            <Link key={app.id} to={`/apps/${app.id}`}>
              <Card className="hover:shadow-lg transition-shadow cursor-pointer">
                <CardHeader>
                  <CardTitle className="flex items-center justify-between">
                    <span>{app.name}</span>
                    <span
                      className={`text-xs px-2 py-1 rounded-full ${
                        app.is_live
                          ? "bg-green-100 text-green-800"
                          : "bg-gray-100 text-gray-800"
                      }`}
                    >
                      {app.is_live ? "Live" : "Draft"}
                    </span>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-2 gap-2 text-sm">
                    <div>
                      <span className="text-muted-foreground">Users:</span>{" "}
                      <strong>{app.user_count}</strong>
                    </div>
                    <div>
                      <span className="text-muted-foreground">Subs:</span>{" "}
                      <strong>{app.subscription_count}</strong>
                    </div>
                    <div>
                      <span className="text-muted-foreground">Plans:</span>{" "}
                      <strong>{app.plan_count}</strong>
                    </div>
                    <div>
                      <span className="text-muted-foreground">Quotes:</span>{" "}
                      <strong>{app.quote_count}</strong>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </Link>
          ))}
        </div>
      </div>
    </div>
  );
}

export default function DashboardView() {
  return (
    <Suspense
      fallback={
        <div className="flex items-center justify-center h-64">
          <div className="animate-pulse text-muted-foreground">Loading...</div>
        </div>
      }
    >
      <DashboardContent />
    </Suspense>
  );
}
