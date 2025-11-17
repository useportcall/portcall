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
          <BasicFeature featureId="sculptor" title="Sculptor" />
          <BasicFeature featureId="sequencer" title="Sequencer" />
          <BasicFeature featureId="exporting" title="Exporting" />
          <BasicFeature featureId="ai_claygent" title="AI/Claygent" />
          <BasicFeature featureId="rollover_credits" title="Rollover Credits" />
          <BasicFeature
            featureId="integration_providers"
            title="100+ integration providers"
          />
          <BasicFeature featureId="scheduling" title="Scheduling" />
          <BasicFeature
            featureId="phone_number_enrichments"
            title="Phone number enrichments"
          />
          <BasicFeature
            featureId="use_your_own_api_keys"
            title="Use your own API keys"
          />
          <BasicFeature featureId="signals" title="Signals" />
          <BasicFeature
            featureId="integrate_with_any_http_api"
            title="Integrate with any HTTP API"
          />
          <BasicFeature featureId="webhooks" title="Webhooks" />
          <BasicFeature
            featureId="email_sequencing_integrations"
            title="Email sequencing integrations"
          />
          <BasicFeature
            featureId="exclude_people_company_filters"
            title="Exclude people/company filters"
          />
          <BasicFeature featureId="web_intent" title="Web Intent" />
          <BasicFeature featureId="crm_integrations" title="CRM Integrations" />
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
  title,
}: {
  featureId: string;
  title: string;
}) {
  const entitlement = await getUserEntitlement(featureId);

  if (!entitlement || !entitlement.enabled) {
    return (
      <Badge className="bg-slate-100 text-slate-600">
        <X />
        {title}
      </Badge>
    );
  }

  return (
    <Badge variant={"outline"}>
      <Check />
      {title}
    </Badge>
  );
}
