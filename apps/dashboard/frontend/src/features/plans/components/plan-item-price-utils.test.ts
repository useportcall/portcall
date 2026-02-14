import { describe, expect, it } from "vitest";
import { formatPriceInput, parsePriceToCents } from "./plan-item-price-utils";

describe("plan-item price utils", () => {
  it("parses decimal dollar values to integer cents", () => {
    expect(parsePriceToCents("12.34")).toBe(1234);
    expect(parsePriceToCents("0.75")).toBe(75);
  });

  it("handles invalid input safely", () => {
    expect(parsePriceToCents("")).toBe(0);
    expect(parsePriceToCents("abc")).toBe(0);
    expect(formatPriceInput("abc")).toBe("0.00");
  });

  it("formats typed values to two decimals", () => {
    expect(formatPriceInput("10")).toBe("10.00");
    expect(formatPriceInput("10.5")).toBe("10.50");
  });
});
