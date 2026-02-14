import { Progress } from "@radix-ui/react-progress";
import { getUserEntitlement } from "./api/get-user-entitlement";
import { Badge } from "./components/ui/badge";
import { MeteredFeatureForm } from "./metered-feature-form";

export async function MeteredFeature({ featureId }: { featureId: string }) {
  const entitlement = await getUserEntitlement(featureId);

  if (!entitlement || entitlement.quota === 0) {
    return <></>;
  }

  const quota = entitlement.quota ?? 0;
  const usage = entitlement.usage ?? 0;

  return (
    <div className="w-full space-y-2 flex flex-col">
      <div className="flex justify-between items-center gap-2">
        <Badge variant={"outline"}>{featureId}</Badge>
        {quota > 0 && (
          <span className="text-xs whitespace-nowrap">
            {usage} / {quota}
          </span>
        )}
        {quota < 0 && (
          <span className="text-xs">{usage}</span>
        )}
      </div>
      {quota > 0 && (
        <Progress
          value={(usage / quota) * 100}
          className="w-full"
        />
      )}
      <MeteredFeatureForm featureId={featureId} />
    </div>
  );
}
