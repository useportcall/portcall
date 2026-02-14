import { Plan } from "@/models/plan";
import { Subscription } from "@/models/subscription";
import { User } from "@/models/user";

const ACTIVE_STATUSES = new Set(["active", "trialing"]);
const RECENT_LIMIT = 5;

type WithDate = { created_at: string };

function sortByCreatedDesc<T extends WithDate>(items: T[]) {
  return [...items].sort((a, b) => {
    return new Date(b.created_at).getTime() - new Date(a.created_at).getTime();
  });
}

export interface HomeSummary {
  hasBillingData: boolean;
  stats: {
    users: number;
    subscribedUsers: number;
    plans: number;
    subscriptions: number;
    activeSubscriptions: number;
  };
  recentUsers: User[];
  recentSubscriptions: Subscription[];
}

export function buildHomeSummary(
  users: User[],
  plans: Plan[],
  subscriptions: Subscription[],
): HomeSummary {
  const recentUsers = sortByCreatedDesc(users).slice(0, RECENT_LIMIT);
  const recentSubscriptions = sortByCreatedDesc(subscriptions).slice(
    0,
    RECENT_LIMIT,
  );

  return {
    hasBillingData: users.length > 0 || plans.length > 0 || subscriptions.length > 0,
    stats: {
      users: users.length,
      subscribedUsers: users.filter((user) => user.subscribed).length,
      plans: plans.length,
      subscriptions: subscriptions.length,
      activeSubscriptions: subscriptions.filter((subscription) =>
        ACTIVE_STATUSES.has(subscription.status),
      ).length,
    },
    recentUsers,
    recentSubscriptions,
  };
}
