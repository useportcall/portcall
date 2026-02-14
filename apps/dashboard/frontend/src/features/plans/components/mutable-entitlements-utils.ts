import { PlanFeature } from "@/models/plan-feature";

type FeatureRecord = { id: string };

export function getLinkedFeatureIds(planFeatures: PlanFeature[]): Set<string> {
  return new Set(planFeatures.map((feature) => feature.feature.id));
}

export function getUnlinkedFeatures(
  planFeatures: PlanFeature[],
  features?: FeatureRecord[],
): FeatureRecord[] {
  if (!features) return [];
  const linked = getLinkedFeatureIds(planFeatures);
  return features.filter((feature) => !linked.has(feature.id));
}
