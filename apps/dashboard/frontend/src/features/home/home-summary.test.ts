import { Plan } from "@/models/plan";
import { Subscription } from "@/models/subscription";
import { User } from "@/models/user";
import { describe, expect, it } from "vitest";
import { buildHomeSummary } from "./home-summary";

function userFixture(id: string, subscribed: boolean, createdAt: string): User {
  return {
    id,
    email: `${id}@test.dev`,
    name: id,
    company_title: null,
    created_at: createdAt,
    updated_at: createdAt,
    subscribed,
    payment_method_added: subscribed,
    billing_address: null,
  };
}

function planFixture(id: string): Plan {
  return {
    id,
    name: id,
    currency: "usd",
    status: "published",
    trial_period_days: 0,
    created_at: "2026-01-01T00:00:00Z",
    updated_at: "2026-01-01T00:00:00Z",
    items: [],
    interval: "month",
    interval_count: 1,
    plan_group: null,
    features: [],
    metered_features: [],
    plan_group_id: null,
    discount_pct: 0,
    discount_qty: 0,
  };
}

function subscriptionFixture(id: string, status: string, createdAt: string): Subscription {
  return {
    id,
    app_id: 1,
    user_id: `usr_${id}`,
    status,
    next_reset_at: createdAt,
    stripe_payment_method_id: null,
    items: [],
    created_at: createdAt,
    user: { email: `${id}@test.dev` },
    plan: { name: "Starter" },
  } as Subscription;
}

describe("buildHomeSummary", () => {
  it("returns an empty summary for new accounts", () => {
    const summary = buildHomeSummary([], [], []);
    expect(summary.hasBillingData).toBe(false);
    expect(summary.stats).toEqual({
      users: 0,
      subscribedUsers: 0,
      plans: 0,
      subscriptions: 0,
      activeSubscriptions: 0,
    });
  });

  it("computes counts and orders recent records", () => {
    const users = [
      userFixture("u1", false, "2026-01-01T00:00:00Z"),
      userFixture("u2", true, "2026-01-03T00:00:00Z"),
      userFixture("u3", true, "2026-01-02T00:00:00Z"),
    ];
    const plans = [planFixture("p1"), planFixture("p2")];
    const subscriptions = [
      subscriptionFixture("s1", "active", "2026-01-02T00:00:00Z"),
      subscriptionFixture("s2", "canceled", "2026-01-01T00:00:00Z"),
      subscriptionFixture("s3", "trialing", "2026-01-03T00:00:00Z"),
    ];

    const summary = buildHomeSummary(users, plans, subscriptions);
    expect(summary.hasBillingData).toBe(true);
    expect(summary.stats.users).toBe(3);
    expect(summary.stats.subscribedUsers).toBe(2);
    expect(summary.stats.plans).toBe(2);
    expect(summary.stats.subscriptions).toBe(3);
    expect(summary.stats.activeSubscriptions).toBe(2);
    expect(summary.recentUsers.map((user) => user.id)).toEqual(["u2", "u3", "u1"]);
    expect(summary.recentSubscriptions.map((subscription) => subscription.id)).toEqual([
      "s3",
      "s1",
      "s2",
    ]);
  });
});
