import { Badge } from "./ui/badge";
import { Progress } from "./ui/progress";
import { MeteredFeatureForm } from "./metered-feature-form";
import { getUserEntitlement } from "@/lib/get-user-entitlement";

export async function MeteredFeature({ featureId }: { featureId: string }) {
  const entitlement = await getUserEntitlement(featureId);

  if (!entitlement) {
    return <></>;
  }

  return (
    <div className="w-full space-y-2 flex flex-col">
      <div className="flex justify-between items-center gap-2">
        <Badge variant={"outline"}>{featureId}</Badge>
        {entitlement.quota > 0 && (
          <span className="text-xs whitespace-nowrap">
            {entitlement.usage} / {entitlement.quota}
          </span>
        )}
        {entitlement.quota < 0 && (
          <span className="text-xs">{entitlement.usage}</span>
        )}
      </div>
      {entitlement.quota > 0 && (
        <Progress
          value={(entitlement.usage / entitlement.quota) * 100}
          className="w-full"
        />
      )}
      <MeteredFeatureForm featureId={featureId} />
    </div>
  );
}
