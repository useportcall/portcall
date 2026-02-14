import { describe, expect, it } from "vitest";
import { formatShortDate, formatUsd, getPricingLabel } from "./format";

describe("billing meter format utils", () => {
  it("formats cents as USD", () => {
    expect(formatUsd(12345)).toBe("$123.45");
  });

  it("maps pricing model labels", () => {
    expect(getPricingLabel("unit")).toBe("per unit");
    expect(getPricingLabel("tiered")).toBe("tiered pricing");
  });

  it("formats short dates", () => {
    expect(formatShortDate("2024-01-02T00:00:00.000Z")).toContain("2024");
  });
});
