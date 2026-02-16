import { getUserSubscription } from "@/lib/api";
import { MeteredFeature } from "./metered-feature";
import { OnOffFeature } from "./on-off-feature";
import { UnsubscribeButton } from "./unsubscribe-button";
import { Separator } from "./ui/separator";

const ON_OFF_FEATURES = [
  { id: "sculptor", name: "Sculptor" },
  { id: "sequencer", name: "Sequencer" },
  { id: "exporting", name: "Exporting" },
  { id: "ai_claygent", name: "AI/Claygent" },
  { id: "rollover_credits", name: "Rollover Credits" },
  { id: "integration_providers", name: "100+ integration providers" },
  { id: "scheduling", name: "Scheduling" },
  { id: "phone_number_enrichments", name: "Phone number enrichments" },
  { id: "use_your_own_api_keys", name: "Use your own API keys" },
  { id: "signals", name: "Signals" },
  {
    id: "integrate_with_any_http_api",
    name: "Integrate with any HTTP API",
  },
  { id: "webhooks", name: "Webhooks" },
  {
    id: "email_sequencing_integrations",
    name: "Email sequencing integrations",
  },
  {
    id: "exclude_people_company_filters",
    name: "Exclude people/company filters",
  },
  { id: "web_intent", name: "Web Intent" },
  { id: "crm_integrations", name: "CRM Integrations" },
];

export default async function Features() {
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
      <div className="flex flex-row justify-between items-center w-full">
        <h2 className="text-xl font-medium">Enabled Features</h2>
        <UnsubscribeButton subscriptionId={subscription.id} />
      </div>
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
          {ON_OFF_FEATURES.map((feature) => (
            <OnOffFeature key={feature.id} featureId={feature.id}>
              {feature.name}
            </OnOffFeature>
          ))}
        </div>
      </div>
    </div>
  );
}
