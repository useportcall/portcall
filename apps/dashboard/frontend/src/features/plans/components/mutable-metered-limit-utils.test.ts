import { describe, expect, it } from "vitest";
import {
  getQuotaTitle,
  parseQuotaInput,
  sanitizeQuotaInput,
} from "./mutable-metered-limit-utils";

describe("mutable-metered-limit utils", () => {
  it("formats quota values for display", () => {
    expect(getQuotaTitle(-1)).toBe("no limit");
    expect(getQuotaTitle(0)).toBe("no limit");
    expect(getQuotaTitle(2500)).toBe("2,500");
  });

  it("sanitizes in-progress quota input", () => {
    expect(sanitizeQuotaInput("-")).toBe("-");
    expect(sanitizeQuotaInput("")).toBe("");
    expect(sanitizeQuotaInput("42.1")).toBe("42");
    expect(sanitizeQuotaInput("not-number")).toBeNull();
  });

  it("parses quota input on submit", () => {
    expect(parseQuotaInput("250")).toBe(250);
    expect(parseQuotaInput("")).toBe(-1);
    expect(parseQuotaInput("x")).toBe(-1);
  });
});
