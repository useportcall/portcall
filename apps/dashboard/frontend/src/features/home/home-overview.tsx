import { useTranslation } from "react-i18next";
import { Card } from "@/components/ui/card";
import { HomeSummary } from "./home-summary";
import { HomeStatCard } from "./home-stat-card";
import { Link } from "react-router-dom";

function formatDate(value: string, locale: string) {
  return new Date(value).toLocaleDateString(locale, {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

export function HomeOverview({ summary }: { summary: HomeSummary }) {
  const { t, i18n } = useTranslation();
  const locale = i18n.resolvedLanguage || i18n.language || "en";

  return (
    <div data-testid="home-overview" className="w-full flex flex-col gap-4">
      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <HomeStatCard
          title={t("home.stats.users")}
          value={summary.stats.users}
          meta={t("home.stats.users_meta", { count: summary.stats.subscribedUsers })}
          testId="home-stat-users"
        />
        <HomeStatCard
          title={t("home.stats.plans")}
          value={summary.stats.plans}
          meta={t("home.stats.plans_meta")}
          testId="home-stat-plans"
        />
        <HomeStatCard
          title={t("home.stats.subscriptions")}
          value={summary.stats.subscriptions}
          meta={t("home.stats.subscriptions_meta", { count: summary.stats.activeSubscriptions })}
          testId="home-stat-subscriptions"
        />
        <HomeStatCard
          title={t("home.stats.health")}
          value={summary.stats.subscriptions === 0 ? 0 : Math.round((summary.stats.activeSubscriptions / summary.stats.subscriptions) * 100)}
          meta={t("home.stats.health_meta")}
          testId="home-stat-health"
        />
      </div>
      <div className="grid gap-4 lg:grid-cols-2">
        <Card data-testid="home-recent-users" className="p-4 space-y-3">
          <p className="font-semibold">{t("home.recent_users")}</p>
          {summary.recentUsers.map((user) => (
            <Link key={user.id} to={`/users/${user.id}`} className="flex items-center justify-between text-sm hover:underline">
              <span className="truncate max-w-[16rem]">{user.email}</span>
              <span className="text-muted-foreground">{formatDate(user.created_at, locale)}</span>
            </Link>
          ))}
          {!summary.recentUsers.length && <p className="text-sm text-muted-foreground">{t("home.no_users")}</p>}
        </Card>
        <Card data-testid="home-recent-subscriptions" className="p-4 space-y-3">
          <p className="font-semibold">{t("home.recent_subscriptions")}</p>
          {summary.recentSubscriptions.map((subscription) => (
            <Link key={subscription.id} to={`/subscriptions/${subscription.id}`} className="flex items-center justify-between text-sm hover:underline">
              <span className="truncate max-w-[16rem]">{subscription.user?.email || subscription.user_id}</span>
              <span className="text-muted-foreground capitalize">{subscription.status}</span>
            </Link>
          ))}
          {!summary.recentSubscriptions.length && (
            <p className="text-sm text-muted-foreground">{t("home.no_subscriptions")}</p>
          )}
        </Card>
      </div>
    </div>
  );
}
