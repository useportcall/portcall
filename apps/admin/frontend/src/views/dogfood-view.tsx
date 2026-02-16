import { Suspense, useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import {
  useDogfoodStatus,
  useDogfoodUsers,
  useDogfoodFeatures,
  useDogfoodPlans,
  useSetupDogfood,
  useCreateDogfoodFeature,
  useRefreshK8sSecrets,
  useFixDogfoodUsers,
  useResetDogfoodPassword,
  useValidateDogfoodUsers,
  DogfoodUser,
  DogfoodFeature,
  DogfoodPlanDetail,
  ResetPasswordResponse,
  ValidateUsersResponse,
} from "@/hooks/use-dogfood";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
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
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { formatDate } from "@/lib/utils";
import {
  Dog,
  Loader2,
  CheckCircle2,
  Users,
  Boxes,
  Key,
  Plus,
  ChevronRight,
  Sparkles,
  Cloud,
  ToggleLeft,
  ToggleRight,
  Crown,
  Zap,
  Wrench,
  KeyRound,
  Copy,
  CheckCheck,
  ShieldCheck,
} from "lucide-react";

const ENV_STORAGE_KEY = "dogfood_environment";

function getStoredEnvironment(): "live" | "test" {
  if (typeof window !== "undefined") {
    const stored = localStorage.getItem(ENV_STORAGE_KEY);
    if (stored === "test" || stored === "live") return stored;
  }
  return "live";
}

function setStoredEnvironment(env: "live" | "test") {
  if (typeof window !== "undefined") {
    localStorage.setItem(ENV_STORAGE_KEY, env);
  }
}

function SetupCard() {
  const setupMutation = useSetupDogfood();
  const [updateK8s, setUpdateK8s] = useState(true);

  return (
    <Card className="max-w-2xl mx-auto mt-12">
      <CardHeader className="text-center">
        <div className="flex justify-center mb-4">
          <Dog className="h-16 w-16 text-amber-500" />
        </div>
        <CardTitle className="text-2xl">Dogfood Account Setup</CardTitle>
        <CardDescription className="text-base">
          Create the dogfood account to enable billing and entitlement
          management for the Portcall dashboard itself.
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <ul className="space-y-2 text-sm text-muted-foreground">
          <li className="flex items-center gap-2">
            <CheckCircle2 className="h-4 w-4 text-green-500" />
            <span>Dogfood account (dogfood@useportcall.com)</span>
          </li>
          <li className="flex items-center gap-2">
            <CheckCircle2 className="h-4 w-4 text-green-500" />
            <span>Live and Test apps with billing-exempt status</span>
          </li>
          <li className="flex items-center gap-2">
            <CheckCircle2 className="h-4 w-4 text-green-500" />
            <span>Free tier and Pro tier plans</span>
          </li>
          <li className="flex items-center gap-2">
            <CheckCircle2 className="h-4 w-4 text-green-500" />
            <span>max_subscriptions and number_of_users features</span>
          </li>
        </ul>

        <div className="flex items-center gap-2 p-3 bg-amber-50 border border-amber-200 rounded-lg">
          <input
            type="checkbox"
            id="updateK8s"
            checked={updateK8s}
            onChange={(e) => setUpdateK8s(e.target.checked)}
            className="h-4 w-4"
          />
          <Label
            htmlFor="updateK8s"
            className="text-sm text-amber-800 cursor-pointer"
          >
            Update Kubernetes secrets in production
          </Label>
        </div>

        {setupMutation.isError && (
          <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-red-800 text-sm">
            Failed to setup: {setupMutation.error?.message}
          </div>
        )}

        {setupMutation.isSuccess && (
          <div className="p-3 bg-green-50 border border-green-200 rounded-lg text-green-800 text-sm space-y-2">
            <p className="font-medium">✅ Dogfood account created!</p>
            {setupMutation.data.k8s_updated && <p>✅ K8s secrets updated</p>}
            <div className="mt-2 p-2 bg-white rounded text-xs font-mono break-all">
              <p>
                <strong>Live:</strong> {setupMutation.data.live_secret}
              </p>
              <p>
                <strong>Test:</strong> {setupMutation.data.test_secret}
              </p>
            </div>
          </div>
        )}

        <Button
          onClick={() => setupMutation.mutate(updateK8s)}
          disabled={setupMutation.isPending}
          className="w-full"
          size="lg"
        >
          {setupMutation.isPending ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Creating...
            </>
          ) : (
            <>
              <Dog className="mr-2 h-4 w-4" />
              Create Dogfood Account
            </>
          )}
        </Button>
      </CardContent>
    </Card>
  );
}

function EnvironmentToggle({
  environment,
  onToggle,
}: {
  environment: "live" | "test";
  onToggle: (env: "live" | "test") => void;
}) {
  return (
    <div className="flex items-center gap-2 p-1 bg-muted rounded-lg">
      <button
        onClick={() => onToggle("test")}
        className={`flex items-center gap-2 px-3 py-1.5 rounded-md text-sm transition-all ${
          environment === "test"
            ? "bg-amber-100 text-amber-800 border border-amber-300"
            : "text-muted-foreground hover:text-foreground"
        }`}
      >
        {environment === "test" ? (
          <ToggleRight className="h-4 w-4" />
        ) : (
          <ToggleLeft className="h-4 w-4" />
        )}
        Test
      </button>
      <button
        onClick={() => onToggle("live")}
        className={`flex items-center gap-2 px-3 py-1.5 rounded-md text-sm transition-all ${
          environment === "live"
            ? "bg-green-100 text-green-800 border border-green-300"
            : "text-muted-foreground hover:text-foreground"
        }`}
      >
        {environment === "live" ? (
          <ToggleRight className="h-4 w-4" />
        ) : (
          <ToggleLeft className="h-4 w-4" />
        )}
        Live
      </button>
    </div>
  );
}

function PlanCard({ plan }: { plan: DogfoodPlanDetail }) {
  const isPro = plan.name.toLowerCase().includes("pro");

  return (
    <Card className={`relative ${isPro ? "border-amber-300 border-2" : ""}`}>
      {isPro && (
        <div className="absolute -top-3 left-1/2 -translate-x-1/2">
          <Badge className="bg-gradient-to-r from-amber-500 to-orange-500 text-white">
            <Crown className="h-3 w-3 mr-1" />
            Pro
          </Badge>
        </div>
      )}
      <CardHeader className="pt-6">
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg flex items-center gap-2">
            {isPro ? (
              <Zap className="h-5 w-5 text-amber-500" />
            ) : (
              <Boxes className="h-5 w-5 text-blue-500" />
            )}
            {plan.name}
          </CardTitle>
          <Badge variant={plan.is_free ? "secondary" : "default"}>
            {plan.is_free ? "Free" : "Paid"}
          </Badge>
        </div>
        <CardDescription>
          {plan.description || "Dashboard subscription tier"}
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-3">
        <div className="flex justify-between text-sm">
          <span className="text-muted-foreground">Subscribers</span>
          <span className="font-medium">{plan.subscriber_count}</span>
        </div>
        {plan.features && plan.features.length > 0 && (
          <div className="space-y-2 pt-2 border-t">
            <h4 className="text-sm font-medium">Quotas</h4>
            {plan.features.map((f) => (
              <div key={f.feature_id} className="flex justify-between text-sm">
                <span className="text-muted-foreground">{f.feature_id}</span>
                <span className="font-medium">
                  {f.quota < 0 ? "Unlimited" : f.quota}
                </span>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}

function UsersTable({
  users,
  plans,
  subscriptionFilter,
  onFilterChange,
  appId,
  environment,
}: {
  users: DogfoodUser[];
  plans: DogfoodPlanDetail[];
  subscriptionFilter: string;
  onFilterChange: (filter: string) => void;
  appId?: number;
  environment: "live" | "test";
}) {
  const navigate = useNavigate();
  const fixUsersMutation = useFixDogfoodUsers(appId);
  const [showFixDialog, setShowFixDialog] = useState(false);
  const [fixResult, setFixResult] = useState<any>(null);

  const filteredUsers = users.filter((user) => {
    if (subscriptionFilter === "all") return true;
    if (subscriptionFilter === "none") return !user.subscription;
    if (subscriptionFilter === "subscribed") return !!user.subscription;
    return (
      user.subscription?.plan_name?.toLowerCase() ===
      subscriptionFilter.toLowerCase()
    );
  });

  // Count users that need fixing
  const usersNeedingFix = users.filter(
    (u) => !u.email.includes("@") || !u.subscription,
  ).length;

  const handleDryRun = async () => {
    try {
      const result = await fixUsersMutation.mutateAsync(true);
      setFixResult(result);
    } catch (e) {
      console.error("Dry run failed:", e);
    }
  };

  const handleApplyFix = async () => {
    try {
      const result = await fixUsersMutation.mutateAsync(false);
      setFixResult(result);
    } catch (e) {
      console.error("Apply fix failed:", e);
    }
  };

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <Label className="text-sm">Filter:</Label>
            <Select value={subscriptionFilter} onValueChange={onFilterChange}>
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="All users" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All users</SelectItem>
                <SelectItem value="subscribed">Has subscription</SelectItem>
                <SelectItem value="none">No subscription</SelectItem>
                {plans.map((plan) => (
                  <SelectItem
                    key={plan.public_id}
                    value={plan.name.toLowerCase()}
                  >
                    {plan.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <span className="text-sm text-muted-foreground">
            {filteredUsers.length} of {users.length}
          </span>
        </div>

        {usersNeedingFix > 0 && (
          <Dialog open={showFixDialog} onOpenChange={setShowFixDialog}>
            <DialogTrigger asChild>
              <Button variant="outline" size="sm" className="gap-2">
                <Wrench className="h-4 w-4" />
                Fix {usersNeedingFix} users
              </Button>
            </DialogTrigger>
            <DialogContent className="max-w-2xl max-h-[80vh] overflow-y-auto">
              <DialogHeader>
                <DialogTitle>Fix User Data</DialogTitle>
                <DialogDescription>
                  Fix malformed emails and create missing subscriptions for
                  users.
                </DialogDescription>
              </DialogHeader>
              {!fixResult ? (
                <div className="space-y-4 py-4">
                  <p className="text-sm text-muted-foreground">
                    Found <strong>{usersNeedingFix}</strong> users that need
                    attention:
                  </p>
                  <ul className="text-sm space-y-1 list-disc pl-4">
                    <li>Users with malformed emails (missing @domain)</li>
                    <li>Users without an active subscription</li>
                  </ul>
                  <p className="text-sm text-muted-foreground">
                    Run a dry run first to preview changes, then apply.
                  </p>
                </div>
              ) : (
                <div className="space-y-4 py-4">
                  <div className="grid grid-cols-2 gap-4 text-sm">
                    <div>
                      Total users: <strong>{fixResult.total_users}</strong>
                    </div>
                    <div>
                      Emails fixed: <strong>{fixResult.emails_fixed}</strong>
                    </div>
                    <div>
                      Subscribed: <strong>{fixResult.subscribed}</strong>
                    </div>
                    <div>
                      Already OK: <strong>{fixResult.already_ok}</strong>
                    </div>
                    {fixResult.failed > 0 && (
                      <div className="text-destructive">
                        Failed: {fixResult.failed}
                      </div>
                    )}
                  </div>
                  {fixResult.dry_run && (
                    <Badge variant="secondary">Dry Run - No changes made</Badge>
                  )}
                  {fixResult.results?.length > 0 && (
                    <div className="border rounded-lg max-h-48 overflow-y-auto">
                      <Table>
                        <TableHeader>
                          <TableRow>
                            <TableHead>User</TableHead>
                            <TableHead>Action</TableHead>
                            <TableHead>Details</TableHead>
                          </TableRow>
                        </TableHeader>
                        <TableBody>
                          {fixResult.results.map((r: any) => (
                            <TableRow key={r.user_id}>
                              <TableCell className="font-mono text-xs">
                                {r.name || r.user_id}
                              </TableCell>
                              <TableCell>
                                <Badge
                                  variant={
                                    r.action === "ok" ? "secondary" : "default"
                                  }
                                >
                                  {r.action}
                                </Badge>
                              </TableCell>
                              <TableCell className="text-xs text-muted-foreground">
                                {r.old_email && (
                                  <span>
                                    Email: {r.old_email} → {r.new_email}
                                  </span>
                                )}
                                {r.subscribed && (
                                  <span>Subscribed to {r.plan_name}</span>
                                )}
                                {r.error && (
                                  <span className="text-destructive">
                                    {r.error}
                                  </span>
                                )}
                              </TableCell>
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    </div>
                  )}
                </div>
              )}
              <DialogFooter className="gap-2">
                {fixResult?.dry_run && (
                  <Button
                    onClick={handleApplyFix}
                    disabled={fixUsersMutation.isPending}
                  >
                    {fixUsersMutation.isPending && (
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    )}
                    Apply Changes
                  </Button>
                )}
                {!fixResult?.dry_run && fixResult && (
                  <Button
                    variant="outline"
                    onClick={() => {
                      setShowFixDialog(false);
                      setFixResult(null);
                    }}
                  >
                    Done
                  </Button>
                )}
                {(!fixResult || fixResult.dry_run) && (
                  <Button
                    variant="outline"
                    onClick={handleDryRun}
                    disabled={fixUsersMutation.isPending}
                  >
                    {fixUsersMutation.isPending && (
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    )}
                    {fixResult ? "Run Again" : "Preview Changes"}
                  </Button>
                )}
              </DialogFooter>
            </DialogContent>
          </Dialog>
        )}
      </div>

      <div className="border rounded-lg">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Email</TableHead>
              <TableHead>Subscription</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Created</TableHead>
              <TableHead className="w-10"></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filteredUsers.map((user) => (
              <TableRow
                key={user.id}
                className="cursor-pointer hover:bg-muted/50"
                onClick={() =>
                  navigate(
                    `/dogfood/users/${user.public_id}?env=${environment}`,
                  )
                }
              >
                <TableCell className="font-medium">{user.name}</TableCell>
                <TableCell className="text-muted-foreground">
                  {user.email}
                </TableCell>
                <TableCell>
                  {user.subscription ? (
                    <Badge
                      variant={
                        user.subscription.plan_name
                          ?.toLowerCase()
                          .includes("pro")
                          ? "default"
                          : "secondary"
                      }
                    >
                      {user.subscription.plan_name || "Subscribed"}
                    </Badge>
                  ) : (
                    <span className="text-muted-foreground">None</span>
                  )}
                </TableCell>
                <TableCell>
                  {user.subscription ? (
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
                  ) : (
                    <Badge variant="outline">No Subscription</Badge>
                  )}
                </TableCell>
                <TableCell>{formatDate(user.created_at)}</TableCell>
                <TableCell>
                  <ChevronRight className="h-4 w-4 text-muted-foreground" />
                </TableCell>
              </TableRow>
            ))}
            {filteredUsers.length === 0 && (
              <TableRow>
                <TableCell
                  colSpan={6}
                  className="text-center text-muted-foreground py-8"
                >
                  {subscriptionFilter === "all"
                    ? "No users found."
                    : "No users match filter."}
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}

function FeaturesTable({
  features,
  appId,
}: {
  features: DogfoodFeature[];
  appId: number;
}) {
  const [isOpen, setIsOpen] = useState(false);
  const [newFeatureId, setNewFeatureId] = useState("");
  const [isMetered, setIsMetered] = useState(true);
  const createMutation = useCreateDogfoodFeature(appId);

  const handleCreate = async () => {
    if (!newFeatureId.trim()) return;
    await createMutation.mutateAsync({
      public_id: newFeatureId.trim(),
      is_metered: isMetered,
    });
    setIsOpen(false);
    setNewFeatureId("");
  };

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <div>
          <h3 className="text-lg font-medium">Features</h3>
          <p className="text-sm text-muted-foreground">
            Entitlement namespaces
          </p>
        </div>
        <Dialog open={isOpen} onOpenChange={setIsOpen}>
          <DialogTrigger asChild>
            <Button size="sm">
              <Plus className="mr-2 h-4 w-4" />
              Add Feature
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Create Feature</DialogTitle>
              <DialogDescription>
                Add a feature for entitlement tracking.
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4 py-4">
              <div className="space-y-2">
                <Label htmlFor="featureId">Feature ID</Label>
                <Input
                  id="featureId"
                  value={newFeatureId}
                  onChange={(e) => setNewFeatureId(e.target.value)}
                  placeholder="e.g., api_calls"
                />
              </div>
              <div className="flex items-center gap-2">
                <input
                  type="checkbox"
                  id="isMetered"
                  checked={isMetered}
                  onChange={(e) => setIsMetered(e.target.checked)}
                  className="h-4 w-4"
                />
                <Label htmlFor="isMetered">Metered feature</Label>
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsOpen(false)}>
                Cancel
              </Button>
              <Button
                onClick={handleCreate}
                disabled={createMutation.isPending || !newFeatureId.trim()}
              >
                {createMutation.isPending ? (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                ) : (
                  <Plus className="mr-2 h-4 w-4" />
                )}
                Create
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      <div className="border rounded-lg">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Feature ID</TableHead>
              <TableHead>Type</TableHead>
              <TableHead>ID</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {features.map((feature) => (
              <TableRow key={feature.id}>
                <TableCell className="font-mono font-medium">
                  {feature.public_id}
                </TableCell>
                <TableCell>
                  <Badge variant={feature.is_metered ? "default" : "secondary"}>
                    {feature.is_metered ? "Metered" : "Boolean"}
                  </Badge>
                </TableCell>
                <TableCell className="text-muted-foreground">
                  {feature.id}
                </TableCell>
              </TableRow>
            ))}
            {features.length === 0 && (
              <TableRow>
                <TableCell
                  colSpan={3}
                  className="text-center text-muted-foreground py-8"
                >
                  No features found.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}

function DogfoodDashboard() {
  const { data: status } = useDogfoodStatus();
  const [environment, setEnvironment] = useState<"live" | "test">(
    getStoredEnvironment,
  );
  const [subscriptionFilter, setSubscriptionFilter] = useState("all");
  const [showPasswordDialog, setShowPasswordDialog] = useState(false);
  const [passwordResult, setPasswordResult] =
    useState<ResetPasswordResponse | null>(null);
  const [copiedField, setCopiedField] = useState<string | null>(null);
  const [showValidateDialog, setShowValidateDialog] = useState(false);
  const [validateResult, setValidateResult] =
    useState<ValidateUsersResponse | null>(null);

  const currentApp =
    environment === "live" ? status?.live_app : status?.test_app;
  const { data: users } = useDogfoodUsers(currentApp?.id);
  const { data: features } = useDogfoodFeatures(currentApp?.id);
  const { data: plans } = useDogfoodPlans(currentApp?.id);
  const refreshK8sMutation = useRefreshK8sSecrets();
  const resetPasswordMutation = useResetDogfoodPassword();
  const validateUsersMutation = useValidateDogfoodUsers();

  const handleResetPassword = async () => {
    try {
      const result = await resetPasswordMutation.mutateAsync();
      setPasswordResult(result);
      setShowPasswordDialog(true);
    } catch (error) {
      console.error("Failed to reset password:", error);
    }
  };

  const handleValidateUsers = async (dryRun: boolean) => {
    try {
      const result = await validateUsersMutation.mutateAsync(dryRun);
      setValidateResult(result);
      setShowValidateDialog(true);
    } catch (error) {
      console.error("Failed to validate users:", error);
    }
  };

  const handleCopy = async (value: string, field: string) => {
    await navigator.clipboard.writeText(value);
    setCopiedField(field);
    setTimeout(() => setCopiedField(null), 2000);
  };

  useEffect(() => {
    setStoredEnvironment(environment);
  }, [environment]);

  if (!status?.configured || !status.live_app) {
    return <SetupCard />;
  }

  if (environment === "test" && !status.test_app) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <Dog className="h-8 w-8 text-amber-500" />
            <h1 className="text-3xl font-bold">Dogfood Account</h1>
          </div>
          <EnvironmentToggle
            environment={environment}
            onToggle={setEnvironment}
          />
        </div>
        <Card className="max-w-lg mx-auto mt-12">
          <CardContent className="pt-6 text-center">
            <p className="text-muted-foreground">No test app configured.</p>
          </CardContent>
        </Card>
      </div>
    );
  }

  const statCards = [
    { title: "Users", value: users?.length || 0, icon: Users },
    { title: "Features", value: features?.length || 0, icon: Sparkles },
    { title: "Plans", value: plans?.length || 0, icon: Boxes },
    { title: "Secrets", value: status.has_secrets ? "Yes" : "No", icon: Key },
  ];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <div className="flex items-center gap-3">
            <Dog className="h-8 w-8 text-amber-500" />
            <h1 className="text-3xl font-bold">Dogfood Account</h1>
            <Badge variant="success">Configured</Badge>
            <Badge variant={environment === "live" ? "success" : "secondary"}>
              {environment === "live" ? "Live" : "Test"}
            </Badge>
          </div>
          <p className="text-muted-foreground mt-1">
            {status.account?.email} • {currentApp?.name}
          </p>
        </div>
        <div className="flex items-center gap-4">
          <EnvironmentToggle
            environment={environment}
            onToggle={setEnvironment}
          />
          <Button
            variant="outline"
            onClick={() => handleValidateUsers(true)}
            disabled={validateUsersMutation.isPending}
          >
            {validateUsersMutation.isPending ? (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            ) : (
              <ShieldCheck className="mr-2 h-4 w-4" />
            )}
            Validate Users
          </Button>
          <Button
            variant="outline"
            onClick={handleResetPassword}
            disabled={resetPasswordMutation.isPending}
          >
            {resetPasswordMutation.isPending ? (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            ) : (
              <KeyRound className="mr-2 h-4 w-4" />
            )}
            Reset Password
          </Button>
          <Button
            variant="outline"
            onClick={() => refreshK8sMutation.mutate()}
            disabled={refreshK8sMutation.isPending}
          >
            {refreshK8sMutation.isPending ? (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            ) : (
              <Cloud className="mr-2 h-4 w-4" />
            )}
            Sync K8s
          </Button>
        </div>
      </div>

      {/* Validate Users Dialog */}
      <Dialog open={showValidateDialog} onOpenChange={setShowValidateDialog}>
        <DialogContent className="sm:max-w-2xl max-h-[80vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <ShieldCheck className="h-5 w-5 text-blue-500" />
              User Validation Results
              {validateResult?.dry_run && (
                <Badge variant="secondary" className="ml-2">
                  Dry Run
                </Badge>
              )}
            </DialogTitle>
            <DialogDescription>
              Validates that dogfood users match their corresponding apps.
            </DialogDescription>
          </DialogHeader>

          {validateResult && (
            <div className="space-y-4">
              {/* Summary Stats */}
              <div className="grid grid-cols-4 gap-2 text-sm">
                <div className="p-2 bg-green-50 rounded text-center">
                  <div className="font-bold text-green-700">
                    {validateResult.valid}
                  </div>
                  <div className="text-green-600 text-xs">Valid</div>
                </div>
                <div className="p-2 bg-amber-50 rounded text-center">
                  <div className="font-bold text-amber-700">
                    {validateResult.wrong_environment}
                  </div>
                  <div className="text-amber-600 text-xs">Wrong Env</div>
                </div>
                <div className="p-2 bg-orange-50 rounded text-center">
                  <div className="font-bold text-orange-700">
                    {validateResult.email_mismatch}
                  </div>
                  <div className="text-orange-600 text-xs">Email Mismatch</div>
                </div>
                <div className="p-2 bg-red-50 rounded text-center">
                  <div className="font-bold text-red-700">
                    {validateResult.missing_app}
                  </div>
                  <div className="text-red-600 text-xs">Missing App</div>
                </div>
              </div>

              {/* Results Table */}
              {validateResult.results.length > 0 && (
                <div className="border rounded-lg max-h-60 overflow-y-auto">
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>User</TableHead>
                        <TableHead>App</TableHead>
                        <TableHead>Status</TableHead>
                        <TableHead>Issue</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {validateResult.results.map((result, i) => (
                        <TableRow key={i}>
                          <TableCell className="text-xs">
                            <div>
                              {result.user_name || result.user_public_id}
                            </div>
                            <div className="text-muted-foreground">
                              {result.user_email}
                            </div>
                          </TableCell>
                          <TableCell className="text-xs">
                            {result.app_name || result.app_public_id || "-"}
                          </TableCell>
                          <TableCell>
                            <Badge
                              variant={
                                result.status === "valid"
                                  ? "success"
                                  : result.status === "orphan"
                                    ? "destructive"
                                    : "secondary"
                              }
                            >
                              {result.status}
                            </Badge>
                            {result.fixed && (
                              <Badge variant="success" className="ml-1">
                                Fixed
                              </Badge>
                            )}
                          </TableCell>
                          <TableCell className="text-xs text-muted-foreground">
                            {result.issue || "-"}
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </div>
              )}

              {validateResult.results.length === 0 && (
                <div className="p-4 bg-green-50 border border-green-200 rounded-lg text-green-800 text-sm text-center">
                  ✅ All {validateResult.total_users} users are valid!
                </div>
              )}

              {validateResult.dry_run && validateResult.results.length > 0 && (
                <div className="p-3 bg-amber-50 border border-amber-200 rounded-lg text-amber-800 text-sm">
                  ⚠️ This was a dry run. Click "Apply Fixes" to fix the{" "}
                  {validateResult.wrong_environment +
                    validateResult.email_mismatch}{" "}
                  issues.
                </div>
              )}
            </div>
          )}

          <DialogFooter className="flex gap-2">
            <Button
              variant="outline"
              onClick={() => setShowValidateDialog(false)}
            >
              Close
            </Button>
            {validateResult?.dry_run && validateResult.results.length > 0 && (
              <Button
                onClick={() => handleValidateUsers(false)}
                disabled={validateUsersMutation.isPending}
              >
                {validateUsersMutation.isPending ? (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                ) : (
                  <Wrench className="mr-2 h-4 w-4" />
                )}
                Apply Fixes
              </Button>
            )}
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Password Reset Success Dialog */}
      <Dialog open={showPasswordDialog} onOpenChange={setShowPasswordDialog}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <CheckCircle2 className="h-5 w-5 text-green-500" />
              Password Reset Successful
            </DialogTitle>
            <DialogDescription>
              Use these credentials to log in to the dashboard as the dogfood
              account owner.
            </DialogDescription>
          </DialogHeader>

          {passwordResult && (
            <div className="space-y-4">
              <div className="p-4 bg-muted rounded-lg space-y-3">
                <div className="flex items-center justify-between">
                  <div>
                    <Label className="text-xs text-muted-foreground">
                      Username
                    </Label>
                    <p className="font-mono font-medium">
                      {passwordResult.username}
                    </p>
                  </div>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() =>
                      handleCopy(passwordResult.username, "username")
                    }
                  >
                    {copiedField === "username" ? (
                      <CheckCheck className="h-4 w-4 text-green-500" />
                    ) : (
                      <Copy className="h-4 w-4" />
                    )}
                  </Button>
                </div>

                <div className="flex items-center justify-between">
                  <div>
                    <Label className="text-xs text-muted-foreground">
                      Email
                    </Label>
                    <p className="font-mono font-medium">
                      {passwordResult.email}
                    </p>
                  </div>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleCopy(passwordResult.email, "email")}
                  >
                    {copiedField === "email" ? (
                      <CheckCheck className="h-4 w-4 text-green-500" />
                    ) : (
                      <Copy className="h-4 w-4" />
                    )}
                  </Button>
                </div>

                <div className="flex items-center justify-between">
                  <div>
                    <Label className="text-xs text-muted-foreground">
                      Password
                    </Label>
                    <p className="font-mono font-medium text-green-600">
                      {passwordResult.password}
                    </p>
                  </div>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() =>
                      handleCopy(passwordResult.password, "password")
                    }
                  >
                    {copiedField === "password" ? (
                      <CheckCheck className="h-4 w-4 text-green-500" />
                    ) : (
                      <Copy className="h-4 w-4" />
                    )}
                  </Button>
                </div>

                <div>
                  <Label className="text-xs text-muted-foreground">Realm</Label>
                  <p className="font-mono text-sm text-muted-foreground">
                    {passwordResult.realm}
                  </p>
                </div>
              </div>

              <div className="p-3 bg-amber-50 border border-amber-200 rounded-lg text-amber-800 text-sm">
                ⚠️ Make sure to copy these credentials now. The password won't
                be shown again.
              </div>
            </div>
          )}

          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setShowPasswordDialog(false)}
            >
              Close
            </Button>
            <Button
              onClick={() => {
                if (passwordResult) {
                  handleCopy(
                    `Username: ${passwordResult.username}\nEmail: ${passwordResult.email}\nPassword: ${passwordResult.password}`,
                    "all",
                  );
                }
              }}
            >
              {copiedField === "all" ? (
                <>
                  <CheckCheck className="mr-2 h-4 w-4" />
                  Copied!
                </>
              ) : (
                <>
                  <Copy className="mr-2 h-4 w-4" />
                  Copy All
                </>
              )}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {resetPasswordMutation.isError && (
        <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-red-800 text-sm">
          ❌ Failed to reset password:{" "}
          {(resetPasswordMutation.error as Error)?.message || "Unknown error"}
        </div>
      )}

      {refreshK8sMutation.isSuccess && (
        <div className="p-3 bg-green-50 border border-green-200 rounded-lg text-green-800 text-sm">
          ✅ {refreshK8sMutation.data.message}
        </div>
      )}

      <div className="grid gap-4 md:grid-cols-4">
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

      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="text-base">
              Current: {currentApp?.name}
            </CardTitle>
            <Badge variant={environment === "live" ? "success" : "secondary"}>
              {environment === "live" ? "Production" : "Development"}
            </Badge>
          </div>
        </CardHeader>
        <CardContent className="grid grid-cols-3 gap-4 text-sm">
          <div>
            <span className="text-muted-foreground">Name</span>
            <p className="font-medium">{currentApp?.name}</p>
          </div>
          <div>
            <span className="text-muted-foreground">Public ID</span>
            <p className="font-mono text-xs">{currentApp?.public_id}</p>
          </div>
          <div>
            <span className="text-muted-foreground">Internal ID</span>
            <p>{currentApp?.id}</p>
          </div>
        </CardContent>
      </Card>

      <div className="space-y-4">
        <div>
          <h2 className="text-xl font-semibold">Plans</h2>
          <p className="text-sm text-muted-foreground">Subscription tiers</p>
        </div>
        <div className="grid gap-4 md:grid-cols-2">
          {plans && plans.length > 0 ? (
            plans.map((plan) => <PlanCard key={plan.public_id} plan={plan} />)
          ) : (
            <Card>
              <CardContent className="pt-6 text-center text-muted-foreground">
                No plans found.
              </CardContent>
            </Card>
          )}
        </div>
      </div>

      <Tabs defaultValue="users" className="space-y-4">
        <TabsList>
          <TabsTrigger value="users">Users ({users?.length || 0})</TabsTrigger>
          <TabsTrigger value="features">
            Features ({features?.length || 0})
          </TabsTrigger>
        </TabsList>
        <TabsContent value="users">
          <UsersTable
            users={users || []}
            plans={plans || []}
            subscriptionFilter={subscriptionFilter}
            onFilterChange={setSubscriptionFilter}
            appId={currentApp?.id}
            environment={environment}
          />
        </TabsContent>
        <TabsContent value="features">
          <FeaturesTable
            features={features || []}
            appId={currentApp?.id || 0}
          />
        </TabsContent>
      </Tabs>
    </div>
  );
}

export default function DogfoodView() {
  return (
    <Suspense
      fallback={
        <div className="flex items-center justify-center h-64">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      }
    >
      <DogfoodDashboard />
    </Suspense>
  );
}
