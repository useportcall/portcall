import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { useGetUserSubscription, useListInvoices } from "@/hooks";
import { Link, useParams } from "react-router-dom";
import { SubscribeUserDialog } from "../subscribe-user-dialog";
import { CreatePaymentLinkDialog } from "../create-payment-link-dialog";
import { CancelUserSubscriptionButton } from "./cancel-user-subscription-button";
import { formatShortDate } from "@/features/billing-meters/utils/format";
import { UserPaymentStatusCard } from "./user-payment-status-card";
import { useTranslation } from "react-i18next";

export function UserSubscriptionSection() {
  const { t } = useTranslation();
  const { id } = useParams();
  const { data: subscription } = useGetUserSubscription(id!);
  const { data: invoices } = useListInvoices({ userId: id! });
  const latestInvoice = invoices?.data?.[0];

  if (!subscription?.data) {
    return (
      <div className="space-y-4">
        <h2 className="text-xl font-semibold">
          {t("views.user.subscription.title")}
        </h2>
        <p className="text-sm text-muted-foreground">
          {t("views.user.subscription.no_active")}
        </p>
        <div className="flex flex-col gap-2">
          <SubscribeUserDialog userId={id!} />
          <CreatePaymentLinkDialog userId={id!} />
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <h2 className="text-xl font-semibold">
        {t("views.user.subscription.title")}
      </h2>
      <SubscriptionField
        label={t("views.user.subscription.id")}
        value={subscription.data.id}
        href={`/subscriptions/${subscription.data.id}`}
      />
      <div className="space-y-1">
        <p className="text-sm text-muted-foreground">
          {t("views.user.subscription.status")}
        </p>
        <Badge variant={subscription.data.status === "active" ? "success" : "outline"}>
          {subscription.data.status}
        </Badge>
      </div>
      <SubscriptionField
        label={t("views.user.subscription.reset_date")}
        value={formatShortDate(subscription.data.next_reset_at)}
      />
      <UserPaymentStatusCard
        subscriptionStatus={subscription.data.status}
        latestInvoice={latestInvoice}
      />
      <div className="space-y-1">
        <p className="text-sm text-muted-foreground">
          {t("views.user.subscription.manage")}
        </p>
        <div className="flex flex-col gap-2">
          <CreatePaymentLinkDialog
            userId={id!}
            userEmail={subscription.data.user?.email}
          />
          <Button disabled variant="outline">
            {t("views.user.subscription.apply_discount_soon")}
          </Button>
          <Button disabled variant="outline">
            {t("views.user.subscription.change_plan_soon")}
          </Button>
          <CancelUserSubscriptionButton subscription={subscription.data} />
        </div>
      </div>
    </div>
  );
}

function SubscriptionField({
  label,
  value,
  href,
}: {
  label: string;
  value: string;
  href?: string;
}) {
  return (
    <div className="space-y-1">
      <p className="text-sm text-muted-foreground">{label}</p>
      {href ? (
        <Link
          to={href}
          target="_blank"
          className="text-xs truncate whitespace-nowrap overflow-ellipsis hover:underline"
        >
          {value}
        </Link>
      ) : (
        <p className="text-xs">{value}</p>
      )}
    </div>
  );
}
