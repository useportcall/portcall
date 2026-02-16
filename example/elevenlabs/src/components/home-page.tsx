import { listFeatures } from "@repo/ui/api/list-features.ts";
import { OnOffFeature } from "@repo/ui/on-off-feature";
import { MeteredFeature } from "@repo/ui/metered-feature";

export default async function HomePage() {
  const FEATURES = await listFeatures(false);

  return (
    <div className="pt-20 h-screen">
      <div className="grid grid-cols-[1fr_1px_1fr] px-10 h-full">
        <div className="space-x-2 space-y-4 h-full p-4">
          <h2 className="text-4xl">Basic Features</h2>
          {FEATURES.data.map((feature) => (
            <OnOffFeature key={feature.id} featureId={feature.id}>
              {feature.id}
            </OnOffFeature>
          ))}
        </div>
        <div className="h-[calc(90%-2rem)] bg-slate-200"></div>
        <MeteredFeatureList />
      </div>

      <pre className="text-xs font-mono p-2 text-slate-100 bg-slate-900">
        {JSON.stringify(FEATURES.data, null, 4)}
      </pre>
    </div>
  );
}

async function MeteredFeatureList() {
  const FEATURES = await listFeatures(true);

  return (
    <div className="space-x-2 space-y-4 h-full p-4">
      <h2 className="text-4xl">Metered Features</h2>
      {FEATURES.data.map((feature) => (
        <MeteredFeature key={feature.id} featureId={feature.id} />
      ))}
    </div>
  );
}
