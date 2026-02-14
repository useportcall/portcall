import { CreateFeatureRequest } from "@/models/feature";
import { toSnakeCase } from "../metered-features/to-snake-case";

export function buildCreateFeatureRequest(
  rawName: string,
): CreateFeatureRequest {
  return {
    feature_id: toSnakeCase(rawName),
    is_metered: false,
  };
}
