import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import {
  useCreateFeature,
  useCreatePlanFeature,
  useDeletePlanFeature,
  useListFeatures,
  useListPlanFeatures,
} from "@/hooks";
import { cn } from "@/lib/utils";
import { PlanFeature } from "@/models/plan-feature";
import { Plus, X } from "lucide-react";
import { useMemo, useState } from "react";
import { useParams } from "react-router-dom";
import { toSnakeCase } from "../utils";
import { PopoverPortal } from "@radix-ui/react-popover";

export default function MutableEntitlements() {
  const { id } = useParams();
  const { data: planFeatures } = useListPlanFeatures({ planId: id! });

  if (!planFeatures?.data) return null;

  return (
    <div className="space-y-2 px-2">
      <div className="flex items-center justify-between">
        <h2 className="text-sm text-muted-foreground">Included features</h2>
        <AddNewFeature planFeatures={planFeatures.data} />
      </div>
      <div className="flex flex-wrap gap-2">
        {planFeatures.data.map((pf) => {
          return <LinkedFeature key={pf.id} planFeature={pf} />;
        })}
      </div>
    </div>
  );
}

function LinkedFeature({ planFeature }: { planFeature: PlanFeature }) {
  const { mutate: removeFeature } = useDeletePlanFeature(planFeature);

  return (
    <Badge variant={"outline"} className={cn("gap-1 cursor-default")}>
      <h3>{planFeature.feature.id}</h3>{" "}
      <X
        className="w-3 h-3 cursor-pointer hover:text-red-800"
        onClick={() => {
          removeFeature({});
        }}
      />
    </Badge>
  );
}

function AddNewFeature({ planFeatures }: { planFeatures: PlanFeature[] }) {
  const [open, setOpen] = useState(false);
  const [value, setValue] = useState("");
  const { id } = useParams();

  const { mutateAsync: addFeature } = useCreateFeature({ isMetered: true });

  const { data: features } = useListFeatures({ isMetered: false });

  const filteredFeatures = useMemo(() => {
    if (!features) return [];

    return (
      features.data.filter((f) => {
        return !planFeatures.find((pf) => pf.feature.id === f.id);
      }) || []
    );
  }, [planFeatures.length, features]);

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          size="icon"
          variant="ghost"
          className="size-[22px] rounded-md hover:bg-slate-100"
        >
          <Plus className="" />
        </Button>
      </PopoverTrigger>
      <PopoverPortal>
        <PopoverContent align="end" side="right">
          <Command>
            <CommandInput
              value={value}
              onValueChange={(v) => setValue(toSnakeCase(v))}
              placeholder="Find entitlements or create a new one"
              className="text-xs w-full outline-none px-2 py-1"
              onKeyDown={async (e) => {
                if (e.key === "Enter") {
                  await addFeature({
                    feature_id: toSnakeCase(value),
                    plan_id: id,
                  });
                  setValue("");
                  setOpen(false);
                }
              }}
            />
            <CommandList>
              <CommandEmpty>No features found.</CommandEmpty>
              <CommandGroup>
                {planFeatures.map((pf) => (
                  <DeletePlanFeatureCommandItem
                    key={pf.id}
                    planFeature={pf}
                    featureId={pf.feature.id}
                  />
                ))}
              </CommandGroup>
              <CommandGroup>
                {filteredFeatures.map((f) => (
                  <PlanFeatureToggle key={f.id} featureId={f.id} />
                ))}
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </PopoverPortal>
    </Popover>
  );
}

function PlanFeatureToggle(props: {
  planFeature?: PlanFeature;
  featureId: string;
}) {
  if (props.planFeature) {
    return (
      <DeletePlanFeatureCommandItem
        planFeature={props.planFeature}
        featureId={props.featureId}
      />
    );
  } else {
    return <AddPlanFeatureCommandItem featureId={props.featureId} />;
  }
}

function DeletePlanFeatureCommandItem(props: {
  planFeature: PlanFeature;
  featureId: string;
}) {
  const { mutate } = useDeletePlanFeature(props.planFeature);

  return (
    <CommandItem
      key={props.featureId}
      onSelect={() => {
        mutate({});
      }}
    >
      <Checkbox checked={true} />
      {props.featureId}
    </CommandItem>
  );
}

function AddPlanFeatureCommandItem(props: { featureId: string }) {
  const { id } = useParams();
  const { mutate } = useCreatePlanFeature(id!);

  return (
    <CommandItem
      key={props.featureId}
      onSelect={() => {
        mutate({ feature_id: props.featureId, plan_id: id! });
      }}
    >
      <Checkbox checked={false} />
      {props.featureId}
    </CommandItem>
  );
}
