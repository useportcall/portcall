import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { useCreateBasicFeature, useListFeatures } from "@/hooks";
import {
  useListBasicEntitlements,
  useToggleUserFeature,
} from "@/hooks/api/entitlements";
import { Loader2, Plus, Search } from "lucide-react";
import { useMemo, useState } from "react";
import { buildCreateFeatureRequest } from "./basic-entitlements/build-create-feature-request";
import { FeatureToggleRow } from "./basic-entitlements/feature-toggle-row";
import { toSnakeCase } from "./metered-features/to-snake-case";

export function EditBasicEntitlementsDialog({ userId }: { userId: string }) {
  const [open, setOpen] = useState(false);
  const [searchValue, setSearchValue] = useState("");
  const { data: entitlements, isLoading: entitlementsLoading } =
    useListBasicEntitlements({ userId });
  const { data: features, isLoading: featuresLoading } = useListFeatures({
    isMetered: false,
  });
  const { mutateAsync: createFeature, isPending: isCreatingFeature } =
    useCreateBasicFeature();
  const { mutate: toggleUserFeature, isPending } = useToggleUserFeature({
    userId,
  });

  const enabledFeatureIds = useMemo(
    () => new Set(entitlements?.data.map((e) => e.id) || []),
    [entitlements?.data],
  );
  const allFeatures = features?.data || [];
  const filteredFeatures = useMemo(
    () =>
      !searchValue
        ? allFeatures
        : allFeatures.filter((f) =>
            f.id.toLowerCase().includes(searchValue.toLowerCase()),
          ),
    [allFeatures, searchValue],
  );
  const canCreateNew = useMemo(
    () =>
      !!searchValue &&
      !allFeatures.find((f) => f.id === toSnakeCase(searchValue)),
    [searchValue, allFeatures],
  );
  const isLoading = entitlementsLoading || featuresLoading;

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button
          size="icon"
          variant="ghost"
          className="size-6 rounded-md hover:bg-accent"
        >
          <Plus className="size-4" />
        </Button>
      </DialogTrigger>
      <DialogContent
        maxWidth="sm:max-w-3xl"
        className="min-h-[600px] max-h-[85vh] overflow-hidden flex flex-col p-0"
      >
        <DialogHeader className="px-6 pt-6 pb-0">
          <DialogTitle className="text-xl">Manage Feature Access</DialogTitle>
          <DialogDescription>
            Toggle features for this user. Type a new feature and press Enter to
            create it.
          </DialogDescription>
        </DialogHeader>
        <div className="flex-1 flex flex-col overflow-hidden px-6 pb-6">
          <div className="relative mt-4 mb-2">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
            <input
              type="text"
              value={searchValue}
              onChange={(e) => setSearchValue(e.target.value)}
              onKeyDown={async (e) =>
                e.key === "Enter" &&
                canCreateNew &&
                (e.preventDefault(),
                await createFeature(buildCreateFeatureRequest(searchValue)),
                setSearchValue(""))
              }
              placeholder="Search features or type to create new..."
              className="w-full h-10 pl-10 pr-4 rounded-lg border border-input bg-transparent text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 dark:focus:ring-offset-background"
            />
          </div>
          {canCreateNew && (
            <div className="flex items-center justify-between px-3 py-2.5 mb-2 rounded-lg bg-primary/5 dark:bg-primary/10 border border-primary/20">
              <div className="flex items-center gap-2">
                <Plus className="size-4 text-primary" />
                <span className="text-sm">
                  Create new feature:{" "}
                  <code className="px-1.5 py-0.5 rounded bg-muted text-xs font-mono">
                    {toSnakeCase(searchValue)}
                  </code>
                </span>
              </div>
              <Button
                size="sm"
                variant="default"
                onClick={async () => {
                  await createFeature(buildCreateFeatureRequest(searchValue));
                  setSearchValue("");
                }}
                disabled={isCreatingFeature}
                className="h-7 text-xs"
              >
                {isCreatingFeature ? (
                  <Loader2 className="size-3 animate-spin" />
                ) : (
                  "Create"
                )}
              </Button>
            </div>
          )}
          <div className="flex-1 overflow-auto rounded-lg border bg-card">
            {isLoading ? (
              <div className="flex items-center justify-center h-32">
                <Loader2 className="size-6 animate-spin text-muted-foreground" />
              </div>
            ) : filteredFeatures.length === 0 && !canCreateNew ? (
              <div className="flex flex-col items-center justify-center h-32 text-muted-foreground">
                <p className="text-sm">No features found</p>
                <p className="text-xs mt-1">
                  Type a name above to create a new feature
                </p>
              </div>
            ) : (
              <div className="divide-y divide-border">
                {filteredFeatures.map((feature) => (
                  <FeatureToggleRow
                    key={feature.id}
                    feature={feature}
                    isEnabled={enabledFeatureIds.has(feature.id)}
                    isPending={isPending}
                    onToggle={() =>
                      !isPending &&
                      toggleUserFeature({
                        feature_id: feature.id,
                        enabled: !enabledFeatureIds.has(feature.id),
                      })
                    }
                  />
                ))}
              </div>
            )}
          </div>
          <div className="flex items-center justify-between pt-4 text-sm text-muted-foreground">
            <span>
              {enabledFeatureIds.size} of {allFeatures.length} features enabled
            </span>
            <Button
              variant="outline"
              size="sm"
              onClick={() => setOpen(false)}
              className="h-8"
            >
              Done
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
