import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { getUserEntitlement, getUserSubscription } from "@/lib/api";
import { Check, X } from "lucide-react";
import { ReactNode } from "react";
import { MeteredFeatureForm } from "./metered-feature-form";
import { Separator } from "./ui/separator";

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
        <MeteredFeature featureId="credits" />
        <Separator className="my-2" />
        <MeteredFeature featureId="users" />
        <Separator className="my-2" />
        <MeteredFeature featureId="people_company_searches" />
      </div>
      <hr />
      <div className="flex flex-col gap-2 w-full">
        <p className="mb-4 font-medium text-sm">Unmetered features</p>
        <div className="flex flex-wrap gap-2">
          <BasicFeature
            featureId="sculptor"
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
            featureId="rollover_credits"
            disabled={
              <Badge className="bg-slate-100 text-slate-600">
                <X />
                Rollover Credits
              </Badge>
            }
            enabled={
              <Badge variant={"outline"}>
                <Check />
                Rollover Credits
              </Badge>
            }
          />
          <BasicFeature
            featureId="integration_providers"
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
            featureId="use_your_own_api_keys"
            disabled={
              <Badge className="bg-slate-100 text-slate-600">
                <X />
                Use your own API keys
              </Badge>
            }
            enabled={
              <Badge variant={"outline"}>
                <Check />
                Use your own API keys
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
          <BasicFeature
            featureId="integrate_with_any_http_api"
            disabled={
              <Badge className="bg-slate-100 text-slate-600">
                <X />
                Integrate with any HTTP API
              </Badge>
            }
            enabled={
              <Badge variant={"outline"}>
                <Check />
                Integrate with any HTTP API
              </Badge>
            }
          />
          <BasicFeature
            featureId="webhooks"
            disabled={
              <Badge className="bg-slate-100 text-slate-600">
                <X />
                Webhooks
              </Badge>
            }
            enabled={
              <Badge variant={"outline"}>
                <Check />
                Webhooks
              </Badge>
            }
          />
          <BasicFeature
            featureId="email_sequencing_integrations"
            disabled={
              <Badge className="bg-slate-100 text-slate-600">
                <X />
                Email sequencing integrations
              </Badge>
            }
            enabled={
              <Badge variant={"outline"}>
                <Check />
                Email sequencing integrations
              </Badge>
            }
          />
          <BasicFeature
            featureId="exclude_people_company_filters"
            disabled={
              <Badge className="bg-slate-100 text-slate-600">
                <X />
                Exclude people/company filters
              </Badge>
            }
            enabled={
              <Badge variant={"outline"}>
                <Check />
                Exclude people/company filters
              </Badge>
            }
          />
          <BasicFeature
            featureId="web_intent"
            disabled={
              <Badge className="bg-slate-100 text-slate-600">
                <X />
                Web Intent
              </Badge>
            }
            enabled={
              <Badge variant={"outline"}>
                <Check />
                Web Intent
              </Badge>
            }
          />
          <BasicFeature
            featureId="crm_integrations"
            disabled={
              <Badge className="bg-slate-100 text-slate-600">
                <X />
                CRM Integrations
              </Badge>
            }
            enabled={
              <Badge variant={"outline"}>
                <Check />
                CRM Integrations
              </Badge>
            }
          />
        </div>
      </div>
    </div>
  );
}

async function MeteredFeature({ featureId }: { featureId: string }) {
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
