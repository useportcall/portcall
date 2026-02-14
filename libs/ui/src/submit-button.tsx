import { getUserSubscription } from "./api/get-user-subscription";
import { SubscribeButton } from "./components/subscribe-button";
import { SwitchSubscriptionPlanButton } from "./components/switch-subscription-plan-button";
import { UnsubscribeButton } from "./components/unsubscribe-button";

export async function PlanSubmitButton({
  planId,
  className,
}: {
  planId: string;
  className?: string;
}) {
  const subscription = await getUserSubscription();

  if (!subscription.data) {
    return <SubscribeButton planId={planId} className={className} />;
  }

  if (subscription.data.plan?.id !== planId) {
    return (
      <SwitchSubscriptionPlanButton
        subscriptionId={subscription.data.id}
        planId={planId}
        className={className}
      />
    );
  }

  return <UnsubscribeButton id={subscription.data.id} className={className} />;
}
