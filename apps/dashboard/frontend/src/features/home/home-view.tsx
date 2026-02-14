import { useTranslation } from "react-i18next";
import { Separator } from "@/components/ui/separator";
import {
  useGetAccount,
  useListPlans,
  useListSubscriptions,
  useListUsers,
} from "@/hooks";
import { HomeOverview } from "./home-overview";
import { HomeQuickStart } from "./home-quick-start";
import { buildHomeSummary } from "./home-summary";

export default function HomeView() {
  const { t } = useTranslation();
  const { data: account } = useGetAccount();
  const { data: users } = useListUsers({ email: "" });
  const { data: plans } = useListPlans("*");
  const { data: subscriptions } = useListSubscriptions();
  const summary = buildHomeSummary(users.data, plans.data, subscriptions.data);
  const firstName = account.data.first_name?.trim();

  return (
    <div className="w-full p-4 lg:p-10 flex flex-col gap-6">
      <div className="flex flex-col gap-4 lg:flex-row justify-between items-start">
        <div className="flex flex-col space-y-2 justify-start">
          <h1 className="text-xl md:text-2xl font-bold">
            {firstName ? t("home.greeting", { name: firstName }) : t("home.greeting_anon")}
          </h1>
          <p className="text-muted-foreground text-sm md:text-base">
            {t("home.subtitle")}
          </p>
        </div>
      </div>
      <Separator />
      {!summary.hasBillingData && <HomeQuickStart />}
      {summary.hasBillingData && <HomeOverview summary={summary} />}
    </div>
  );
}
