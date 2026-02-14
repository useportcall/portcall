import { useGetApp } from "@/hooks/api/apps";
import { UsageContent } from "./components/usage-content";

/**
 * UsageView is the main component for the usage feature.
 *
 * ENTITLEMENT GATE UI: This view should only be shown to non-dogfood apps.
 * Dogfood apps (billing_exempt === true) have unlimited usage and don't need
 * to see this view.
 */
export default function UsageView() {
  const { data: app } = useGetApp();

  // Don't show usage view for dogfood (billing-exempt) apps
  if (app?.data?.billing_exempt) {
    return null;
  }

  const isLiveApp = app?.data?.is_live ?? false;

  return <UsageContent isLiveApp={isLiveApp} />;
}
