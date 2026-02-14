import { PlanFeature } from "@/models/plan-feature";
import { describe, expect, it } from "vitest";
import {
  getLinkedFeatureIds,
  getUnlinkedFeatures,
} from "./mutable-entitlements-utils";

function planFeature(id: string): PlanFeature {
  return {
    id: `pf_${id}`,
    interval: "no_reset",
    quota: -1,
    rollover: 0,
    feature: { id },
    plan_item: null,
    created_at: "2026-01-01T00:00:00Z",
    updated_at: "2026-01-01T00:00:00Z",
  };
}

describe("mutable-entitlements utils", () => {
  it("builds a set of linked feature ids", () => {
    const ids = getLinkedFeatureIds([planFeature("analytics"), planFeature("sso")]);
    expect(ids.has("analytics")).toBe(true);
    expect(ids.has("sso")).toBe(true);
    expect(ids.size).toBe(2);
  });

  it("returns only features not already linked to the plan", () => {
    const planFeatures = [planFeature("analytics"), planFeature("sso")];
    const options = [{ id: "analytics" }, { id: "priority_support" }, { id: "sso" }];
    expect(getUnlinkedFeatures(planFeatures, options)).toEqual([{ id: "priority_support" }]);
  });

  it("returns empty list when feature options are missing", () => {
    expect(getUnlinkedFeatures([planFeature("analytics")], undefined)).toEqual([]);
  });
});
