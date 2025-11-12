import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import { getUserEntitlement, getUserSubscription } from "@/lib/api";
import { Check, X } from "lucide-react";
import { ReactNode } from "react";

export default async function EnabledFeatures() {
  const subscription = await getUserSubscription();

  if (!subscription) {
    return (
      <div className=" p-8 flex justify-center items-center h-full text-muted-foreground border-r">
        No subscription found
      </div>
    );
  }

  return (
    <div className="text-xs gap-4 space-y-4 p-8">
      <h2 className="text-xl font-medium">Enabled Features</h2>
      <hr />
      <div className="flex flex-col w-full gap-2">
        <p className="mb-4 font-medium text-sm">Metered features</p>
        <FeatureProgress featureId="credits" />
        <FeatureProgress featureId="people_company_searches" />
        <FeatureProgress featureId="users" />
      </div>
      <hr />
      <div className="flex flex-col gap-2 w-full">
        <p className="mb-4 font-medium text-sm">Unmetered features</p>
        <div className="flex flex-wrap gap-2">
          <BasicFeature
            featureId="scultor"
            disabled={
              <Badge className="bg-slate-100 text-slate-600">
                <X />
                Sculptor
              </Badge>
            }
            enabled={
              <Badge variant={"outline"}>
                <Check />
                Sculptor
              </Badge>
            }
          />
          <BasicFeature
            featureId="sequencer"
            disabled={
              <Badge className="bg-slate-100 text-slate-600">
                <X />
                Sequencer
              </Badge>
            }
            enabled={
              <Badge variant={"outline"}>
                <Check />
                Sequencer
              </Badge>
            }
          />
          <BasicFeature
            featureId="exporting"
            disabled={
              <Badge className="bg-slate-100 text-slate-600">
                <X />
                Exporting
              </Badge>
            }
            enabled={
              <Badge variant={"outline"}>
                <Check />
                Exporting
              </Badge>
            }
          />
          <BasicFeature
            featureId="ai_claygent"
            disabled={
              <Badge className="bg-slate-100 text-slate-600">
                <X />
                AI/Claygent
              </Badge>
            }
            enabled={
              <Badge variant={"outline"}>
                <Check />
                AI/Claygent
              </Badge>
            }
          />
          <BasicFeature
            featureId="100_integration_providers"
            disabled={
              <Badge className="bg-slate-100 text-slate-600">
                <X /> 100+ integration providers
              </Badge>
            }
            enabled={
              <Badge variant={"outline"}>
                <Check />
                100+ integration providers
              </Badge>
            }
          />
          <BasicFeature
            featureId="scheduling"
            disabled={
              <Badge className="bg-slate-100 text-slate-600">
                <X />
                Scheduling
              </Badge>
            }
            enabled={
              <Badge variant={"outline"}>
                <Check />
                Scheduling
              </Badge>
            }
          />
          <BasicFeature
            featureId="phone_number_enrichments"
            disabled={
              <Badge className="bg-slate-100 text-slate-600">
                <X /> Phone number enrichments
              </Badge>
            }
            enabled={
              <Badge variant={"outline"}>
                <Check />
                Phone number enrichments
              </Badge>
            }
          />
          <BasicFeature
            featureId="use_own_api_keys"
            disabled={
              <Badge className="bg-slate-100 text-slate-600">
                <X />
                Use own API keys
              </Badge>
            }
            enabled={
              <Badge variant={"outline"}>
                <Check />
                Use own API keys
              </Badge>
            }
          />
          <BasicFeature
            featureId="signals"
            disabled={
              <Badge className="bg-slate-100 text-slate-600">
                <X />
                Signals
              </Badge>
            }
            enabled={
              <Badge variant={"outline"}>
                <Check />
                Signals
              </Badge>
            }
          />
        </div>
      </div>
    </div>
  );
}

async function FeatureProgress({ featureId }: { featureId: string }) {
  const entitlement = await getUserEntitlement(featureId);
  // const { mutateAsync } = useCreateMeterEvent(featureId);

  if (!entitlement) {
    return <></>;
  }

  return (
    <div className="w-full space-y-2 flex flex-col">
      <div className="flex justify-between">
        <Badge variant={"outline"}>{featureId}</Badge>
        {entitlement.quota > 0 && (
          <span className="text-xs">
            {entitlement.usage} / {entitlement.quota}
          </span>
        )}
        {entitlement.quota < 0 && <span className="text-xs">Unlimited</span>}
      </div>
      <Progress
        value={(entitlement.usage / entitlement.quota) * 100}
        className="w-full"
      />
      <Button size="sm" className="w-fit mt-1 text-xs">
        Increment usage 🔼
      </Button>
    </div>
  );
}

async function BasicFeature({
  featureId,
  enabled,
  disabled,
}: {
  featureId: string;
  enabled: ReactNode;
  disabled: ReactNode;
}) {
  const entitlement = await getUserEntitlement(featureId);

  if (!entitlement) {
    return disabled;
  }

  if (!entitlement.enabled) {
    return disabled;
  }

  return enabled;
}
