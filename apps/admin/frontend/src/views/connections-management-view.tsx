import { Suspense, useState } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  useCheckMockConnections,
  useClearMockConnections,
} from "@/hooks/use-connections";
import {
  AlertTriangle,
  CheckCircle2,
  Loader2,
  RefreshCw,
  Trash2,
  Zap,
} from "lucide-react";
import { toast } from "@/lib/toast";

function ConnectionsManagementContent() {
  const { data, isLoading, refetch, isRefetching } = useCheckMockConnections();
  const { mutateAsync: clearMock, isPending: isClearing } =
    useClearMockConnections();
  const [showConfirmDialog, setShowConfirmDialog] = useState(false);

  const handleClearMockConnections = async () => {
    try {
      const result = await clearMock();
      toast.success(
        `Cleared ${result.cleared} mock connection(s) from live apps`,
      );
      if (result.errors && result.errors.length > 0) {
        toast.error(`Some errors occurred: ${result.errors.join(", ")}`);
      }
      setShowConfirmDialog(false);
    } catch (error) {
      toast.error("Failed to clear mock connections");
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Connections Management</h1>
          <p className="text-muted-foreground mt-1">
            Monitor and manage payment connections across all apps
          </p>
        </div>
        <Button
          variant="outline"
          size="sm"
          onClick={() => refetch()}
          disabled={isRefetching}
        >
          {isRefetching ? (
            <Loader2 className="w-4 h-4 animate-spin mr-2" />
          ) : (
            <RefreshCw className="w-4 h-4 mr-2" />
          )}
          Refresh
        </Button>
      </div>

      {/* Status Card */}
      <Card
        className={
          data?.has_mock_connections ? "border-amber-500" : "border-green-500"
        }
      >
        <CardHeader className="pb-3">
          <div className="flex items-center gap-3">
            {data?.has_mock_connections ? (
              <AlertTriangle className="w-6 h-6 text-amber-500" />
            ) : (
              <CheckCircle2 className="w-6 h-6 text-green-500" />
            )}
            <div>
              <CardTitle className="text-lg">
                {data?.has_mock_connections
                  ? "Mock Connections Found in Live Apps"
                  : "No Mock Connections in Live Apps"}
              </CardTitle>
              <CardDescription>
                {data?.total_live_apps} live app(s) total,{" "}
                {data?.total_with_mock} with mock connections
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        {data?.has_mock_connections && (
          <CardContent>
            <Button
              variant="destructive"
              size="sm"
              onClick={() => setShowConfirmDialog(true)}
              disabled={isClearing}
            >
              {isClearing ? (
                <Loader2 className="w-4 h-4 animate-spin mr-2" />
              ) : (
                <Trash2 className="w-4 h-4 mr-2" />
              )}
              Clear All Mock Connections from Live Apps
            </Button>
          </CardContent>
        )}
      </Card>

      {/* Affected Apps Table */}
      {data?.live_apps_with_mock && data.live_apps_with_mock.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Live Apps with Mock Connections</CardTitle>
            <CardDescription>
              These live apps have mock/local payment connections that should be
              replaced with real providers
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>App Name</TableHead>
                  <TableHead>App ID</TableHead>
                  <TableHead>Connection Name</TableHead>
                  <TableHead>Status</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {data.live_apps_with_mock.map((item, index) => (
                  <TableRow key={index}>
                    <TableCell className="font-medium">
                      {item.app_name}
                    </TableCell>
                    <TableCell className="font-mono text-sm">
                      {item.app_public_id}
                    </TableCell>
                    <TableCell>{item.connection_name || "N/A"}</TableCell>
                    <TableCell>
                      <Badge variant="destructive" className="gap-1">
                        <Zap className="w-3 h-3" />
                        Mock
                      </Badge>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      )}

      {/* Confirmation Dialog */}
      <Dialog open={showConfirmDialog} onOpenChange={setShowConfirmDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Clear Mock Connections</DialogTitle>
            <DialogDescription>
              This will permanently remove all mock/local payment connections
              from live apps. These apps will need to set up real payment
              providers.
            </DialogDescription>
          </DialogHeader>
          <div className="py-4">
            <p className="text-sm text-muted-foreground">
              <strong>{data?.total_with_mock}</strong> mock connection(s) will
              be removed from{" "}
              <strong>{data?.live_apps_with_mock?.length}</strong> live app(s).
            </p>
          </div>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setShowConfirmDialog(false)}
              disabled={isClearing}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={handleClearMockConnections}
              disabled={isClearing}
            >
              {isClearing ? (
                <>
                  <Loader2 className="w-4 h-4 animate-spin mr-2" />
                  Clearing...
                </>
              ) : (
                <>
                  <Trash2 className="w-4 h-4 mr-2" />
                  Clear All
                </>
              )}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}

export default function ConnectionsManagementView() {
  return (
    <Suspense
      fallback={
        <div className="flex items-center justify-center h-64">
          <div className="animate-pulse text-muted-foreground">Loading...</div>
        </div>
      }
    >
      <ConnectionsManagementContent />
    </Suspense>
  );
}
