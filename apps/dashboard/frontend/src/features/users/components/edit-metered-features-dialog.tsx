import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { useCreateFeature, useListFeatures } from "@/hooks";
import {
  useCreateEntitlement,
  useListMeteredEntitlements,
} from "@/hooks/api/entitlements";
import { Plus } from "lucide-react";
import { useMemo, useState } from "react";
import { ActiveMeteredFeaturesPanel } from "./metered-features/active-metered-features-panel";
import { AvailableMeteredFeaturesPanel } from "./metered-features/available-metered-features-panel";
import { toSnakeCase } from "./metered-features/to-snake-case";

export function EditMeteredFeaturesDialog({ userId }: { userId: string }) {
  const [open, setOpen] = useState(false);
  const [searchValue, setSearchValue] = useState("");
  const { data: entitlements, isLoading: entitlementsLoading } =
    useListMeteredEntitlements({ userId });
  const { data: features, isLoading: featuresLoading } = useListFeatures({
    isMetered: true,
  });
  const { mutateAsync: createFeature, isPending: isCreatingFeature } =
    useCreateFeature({ isMetered: true });
  const { mutate: createEntitlement, isPending: isCreatingEntitlement } =
    useCreateEntitlement({ userId });

  const allFeatures = features?.data || [];
  const enabledFeatureIds = useMemo(
    () => new Set(entitlements?.data.map((e) => e.id) || []),
    [entitlements?.data],
  );
  const availableFeatures = useMemo(
    () =>
      allFeatures
        .filter((f) => !enabledFeatureIds.has(f.id))
        .filter(
          (f) =>
            !searchValue ||
            f.id.toLowerCase().includes(searchValue.toLowerCase()),
        ),
    [allFeatures, enabledFeatureIds, searchValue],
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
        maxWidth="sm:max-w-4xl"
        className="min-h-[700px] max-h-[90vh] overflow-hidden flex flex-col p-0"
      >
        <DialogHeader className="px-6 pt-6 pb-0">
          <DialogTitle className="text-xl">Manage Metered Features</DialogTitle>
          <DialogDescription>
            Configure usage quotas and track consumption for this user. Type a
            new feature name and press Enter to create it.
          </DialogDescription>
        </DialogHeader>
        <div className="flex-1 flex gap-6 overflow-hidden px-6 pb-6 pt-4">
          <AvailableMeteredFeaturesPanel
            searchValue={searchValue}
            setSearchValue={setSearchValue}
            canCreateNew={canCreateNew}
            handleCreateFeature={async () => {
              await createFeature({
                feature_id: toSnakeCase(searchValue),
                is_metered: true,
              });
              setSearchValue("");
            }}
            isCreatingFeature={isCreatingFeature}
            isLoading={isLoading}
            availableFeatures={availableFeatures}
            handleAddEntitlement={(featureId) =>
              createEntitlement({ feature_id: featureId, quota: 100 })
            }
            isCreatingEntitlement={isCreatingEntitlement}
          />
          <ActiveMeteredFeaturesPanel
            isLoading={isLoading}
            entitlements={entitlements?.data || []}
            userId={userId}
          />
        </div>
        <div className="flex items-center justify-end px-6 py-4 border-t bg-muted/30">
          <Button variant="outline" onClick={() => setOpen(false)}>
            Done
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
