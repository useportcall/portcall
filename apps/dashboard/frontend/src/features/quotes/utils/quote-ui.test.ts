import { Quote } from "@/models/quote";
import {
  buildRecipientDetails,
  getSendButtonTitle,
  hasRecipientDetails,
  isQuoteReadonly,
  isSendDisabled,
  isVoidDisabled,
  summarizeTerms,
} from "./quote-ui";
import { describe, expect, it } from "vitest";

function createQuote(overrides: Partial<Quote> = {}): Quote {
  return {
    id: "quote_1",
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z",
    expires_at: null,
    plan: {
      id: "plan_1",
      name: "Starter",
      currency: "USD",
      status: "active",
      trial_period_days: 0,
      created_at: "2024-01-01T00:00:00Z",
      updated_at: "2024-01-01T00:00:00Z",
      items: [],
      interval: "month",
      interval_count: 1,
      plan_group: null,
      features: [],
      metered_features: [],
      plan_group_id: null,
      discount_pct: 0,
      discount_qty: 0,
    },
    user_id: null,
    status: "draft",
    company_name: null,
    direct_checkout_enabled: false,
    recipient_address: {
      id: "addr_1",
      line1: "",
      city: "",
      postal_code: "",
      country: "",
      created_at: "2024-01-01T00:00:00Z",
      updated_at: "2024-01-01T00:00:00Z",
    },
    recipient_email: null,
    recipient_name: null,
    recipient_title: null,
    toc: "",
    prepared_by_email: null,
    ...overrides,
  };
}

describe("quote-ui", () => {
  it("returns correct send button title", () => {
    expect(getSendButtonTitle("sent")).toBe("Resend");
    expect(getSendButtonTitle("draft")).toBe("Send");
  });

  it("disables send and void actions for terminal statuses", () => {
    expect(isSendDisabled("rejected")).toBe(true);
    expect(isVoidDisabled("rejected")).toBe(true);
    expect(isSendDisabled("sent")).toBe(false);
    expect(isVoidDisabled("sent")).toBe(false);
    expect(isQuoteReadonly("sent")).toBe(true);
    expect(isQuoteReadonly("draft")).toBe(false);
  });

  it("builds a compact recipient details summary", () => {
    const quote = createQuote({
      recipient_email: "a@example.com",
      recipient_address: {
        ...createQuote().recipient_address,
        line1: "123 Main",
        city: "NYC",
      },
    });

    expect(buildRecipientDetails(quote)).toEqual([
      "a@example.com",
      "123 Main",
      "NYC",
    ]);
  });

  it("detects when recipient details exist", () => {
    expect(hasRecipientDetails(createQuote())).toBe(false);
    expect(hasRecipientDetails(createQuote({ recipient_name: "Jane" }))).toBe(
      true,
    );
  });

  it("truncates long terms text", () => {
    expect(summarizeTerms("short", 10)).toBe("short");
    expect(summarizeTerms("12345678901", 10)).toBe("1234567890...");
  });
});
