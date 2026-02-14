import { Button } from "@/components/ui/button";
import { Loader2, Plus, Search } from "lucide-react";
import { toSnakeCase } from "./to-snake-case";

export function AvailableMeteredFeaturesPanel(props: {
  searchValue: string;
  setSearchValue: (value: string) => void;
  canCreateNew: boolean;
  handleCreateFeature: () => Promise<void>;
  isCreatingFeature: boolean;
  isLoading: boolean;
  availableFeatures: { id: string }[];
  handleAddEntitlement: (featureId: string) => void;
  isCreatingEntitlement: boolean;
}) {
  return (
    <div className="w-72 shrink-0 flex flex-col gap-3">
      <h3 className="text-sm font-medium text-foreground">Add Metered Feature</h3>
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
        <input
          type="text"
          value={props.searchValue}
          onChange={(e) => props.setSearchValue(e.target.value)}
          onKeyDown={async (e) => e.key === "Enter" && props.canCreateNew && (e.preventDefault(), await props.handleCreateFeature())}
          placeholder="Search or create..."
          className="w-full h-9 pl-9 pr-4 rounded-md border border-input bg-transparent text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-1 dark:focus:ring-offset-background"
        />
      </div>
      {props.canCreateNew && (
        <div className="flex items-center justify-between px-3 py-2 rounded-md bg-primary/5 dark:bg-primary/10 border border-primary/20">
          <div className="flex items-center gap-2 min-w-0"><Plus className="size-3.5 text-primary shrink-0" /><code className="text-xs font-mono truncate">{toSnakeCase(props.searchValue)}</code></div>
          <Button size="sm" variant="default" onClick={props.handleCreateFeature} disabled={props.isCreatingFeature} className="h-6 text-xs px-2">
            {props.isCreatingFeature ? <Loader2 className="size-3 animate-spin" /> : "Create"}
          </Button>
        </div>
      )}
      <div className="flex-1 overflow-auto rounded-md border bg-card">
        {props.isLoading ? <LoadingState /> : <AvailableFeaturesList {...props} />}
      </div>
    </div>
  );
}

function LoadingState() {
  return <div className="flex items-center justify-center h-24"><Loader2 className="size-5 animate-spin text-muted-foreground" /></div>;
}

function AvailableFeaturesList(props: {
  searchValue: string;
  availableFeatures: { id: string }[];
  handleAddEntitlement: (featureId: string) => void;
  isCreatingEntitlement: boolean;
}) {
  if (!props.availableFeatures.length) {
    return <div className="flex flex-col items-center justify-center h-24 text-muted-foreground text-center p-4"><p className="text-xs">{props.searchValue ? "No matching features" : "All features assigned"}</p></div>;
  }
  return (
    <div className="divide-y divide-border">
      {props.availableFeatures.map((feature) => (
        <button
          key={feature.id}
          onClick={() => props.handleAddEntitlement(feature.id)}
          disabled={props.isCreatingEntitlement}
          className="w-full flex items-center gap-2 px-3 py-2.5 text-left hover:bg-accent/50 transition-colors disabled:opacity-50"
        >
          <Plus className="size-3.5 text-muted-foreground" />
          <span className="font-mono text-sm truncate">{feature.id}</span>
        </button>
      ))}
    </div>
  );
}
