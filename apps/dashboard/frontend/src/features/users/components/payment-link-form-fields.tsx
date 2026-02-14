import { useTranslation } from "react-i18next";

export function PaymentLinkFormFields(props: {
  selectedPlanID: string;
  expiresInDays: string;
  plans: { id: string; name: string }[];
  onPlanChange: (value: string) => void;
  onDaysChange: (value: string) => void;
}) {
  const { t } = useTranslation();
  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <p className="text-sm text-muted-foreground">{t("views.user.payment_link.plan_label")}</p>
        <select
          data-testid="payment-link-plan-select"
          value={props.selectedPlanID}
          onChange={(event) => props.onPlanChange(event.target.value)}
          className="w-full h-10 px-3 rounded-lg border border-input bg-transparent text-sm"
        >
          <option value="">{t("views.user.payment_link.plan_placeholder")}</option>
          {props.plans.map((plan) => (
            <option key={plan.id} value={plan.id}>
              {plan.name}
            </option>
          ))}
        </select>
      </div>
      <div className="space-y-2">
        <p className="text-sm text-muted-foreground">{t("views.user.payment_link.expires_label")}</p>
        <input
          data-testid="payment-link-expiry-days"
          type="number"
          min={1}
          max={90}
          value={props.expiresInDays}
          onChange={(event) => props.onDaysChange(event.target.value)}
          className="w-full h-10 px-3 rounded-lg border border-input bg-transparent text-sm"
        />
      </div>
    </div>
  );
}
