import { Suspense } from "react";
import { useApps } from "@/hooks/use-apps";
import { Link } from "react-router-dom";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { formatDate } from "@/lib/utils";

function AppsListContent() {
  const { data: apps } = useApps();

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Apps</h1>
        <p className="text-muted-foreground mt-1">
          All registered apps in the platform
        </p>
      </div>

      <div className="border rounded-lg">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="text-right">Users</TableHead>
              <TableHead className="text-right">Subscriptions</TableHead>
              <TableHead className="text-right">Plans</TableHead>
              <TableHead className="text-right">Quotes</TableHead>
              <TableHead className="text-right">Invoices</TableHead>
              <TableHead>Created</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {apps.map((app) => (
              <TableRow key={app.id}>
                <TableCell>
                  <Link
                    to={`/apps/${app.id}`}
                    className="font-medium hover:underline text-primary"
                  >
                    {app.name}
                  </Link>
                </TableCell>
                <TableCell>
                  <Badge variant={app.is_live ? "success" : "secondary"}>
                    {app.is_live ? "Live" : "Draft"}
                  </Badge>
                </TableCell>
                <TableCell className="text-right">{app.user_count}</TableCell>
                <TableCell className="text-right">
                  {app.subscription_count}
                </TableCell>
                <TableCell className="text-right">{app.plan_count}</TableCell>
                <TableCell className="text-right">{app.quote_count}</TableCell>
                <TableCell className="text-right">
                  {app.invoice_count}
                </TableCell>
                <TableCell>{formatDate(app.created_at)}</TableCell>
              </TableRow>
            ))}
            {apps.length === 0 && (
              <TableRow>
                <TableCell
                  colSpan={8}
                  className="text-center text-muted-foreground"
                >
                  No apps found
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}

export default function AppsListView() {
  return (
    <Suspense
      fallback={
        <div className="flex items-center justify-center h-64">
          <div className="animate-pulse text-muted-foreground">Loading...</div>
        </div>
      }
    >
      <AppsListContent />
    </Suspense>
  );
}
