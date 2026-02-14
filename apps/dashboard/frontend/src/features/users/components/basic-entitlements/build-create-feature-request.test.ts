import { describe, expect, it } from "vitest";
import { buildCreateFeatureRequest } from "./build-create-feature-request";

describe("buildCreateFeatureRequest", () => {
  it("uses feature_id field (not id)", () => {
    const req = buildCreateFeatureRequest("My Feature");
    expect(req).toHaveProperty("feature_id");
    expect(req).not.toHaveProperty("id");
  });

  it("snake-cases the feature name", () => {
    const req = buildCreateFeatureRequest("My Feature Name");
    expect(req.feature_id).toBe("my_feature_name");
  });

  it("sets is_metered to false", () => {
    const req = buildCreateFeatureRequest("analytics");
    expect(req.is_metered).toBe(false);
  });

  it("handles separators in names", () => {
    const req = buildCreateFeatureRequest("feature/path-name");
    expect(req.feature_id).toBe("feature_path_name");
  });

  it("trims whitespace", () => {
    const req = buildCreateFeatureRequest("  spaced  ");
    expect(req.feature_id).toBe("spaced");
  });
});
