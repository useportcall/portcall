import { SaveIndicator } from "@/components/save-indicator";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { useRetrievePlan } from "@/hooks";
import { ArrowLeft } from "lucide-react";
import { useNavigate, useParams } from "react-router-dom";
import { BillInAdvanceSelect } from "./components/bill-in-advance-select";
import { FreeTrialComboBox } from "./components/free-trial-combo-box";
import MutableEntitlements from "./components/mutable-entitlements";
import { PlanCurrencySelect } from "./components/plan-currency-select";
import { PlanFixedUnitAmountInput } from "./components/plan-fixed-unit-amount-input";
import { PlanGroupSelect } from "./components/plan-group-select";
import { PlanIntervalSelect } from "./components/plan-interval-select";
import { PlanItems } from "./components/plan-items";
import { PlanNameInput } from "./components/plan-name-input";
import { ProrationSelect } from "./components/proration-select";
import { PublishPlanButton } from "./components/publish-plan-button";

export function EditPlan() {
  const { id } = useParams();
  const { data: plan } = useRetrievePlan(id!);

  const navigate = useNavigate();

  if (!plan) {
    return <div></div>;
  }

  return (
    <div className="w-full h-full flex flex-col gap-4 p-4 overflow-auto">
      <Button
        variant="ghost"
        className="mb-4 p-0 self-start text-sm flex flex-row justify-start items-center gap-2"
        onClick={() => navigate("/plans")}
      >
        <ArrowLeft className="w-4 h-4" /> Back to plans
      </Button>
      <div className="flex gap-4 w-full justify-between px-4">
        <PlanNameInput plan={plan.data} />
        <div className="flex flex-row justify-end gap-2">
          <SaveIndicator />
          <PublishPlanButton plan={plan.data} />
        </div>
      </div>
      <Separator />
      <div className="grid grid-cols-1 md:grid-cols-[20rem_auto_2fr] w-full h-full">
        <div className="flex flex-col gap-6 px-2 animate-fade-in">
          <Section title="Fixed price">
            <PlanFixedUnitAmountInput plan={plan.data} />
            <PlanCurrencySelect plan={plan.data} />
            <PlanIntervalSelect plan={plan.data} />
            <PlanGroupSelect plan={plan.data} />
          </Section>
          <MutableEntitlements />
          <Section title="Other properties">
            <FreeTrialComboBox plan={plan.data} />
            <BillInAdvanceSelect />
            <ProrationSelect />
          </Section>
        </div>
        <Separator orientation="vertical" className="hidden md:block" />
        <PlanItems />
      </div>
    </div>
  );
}

function Section(props: { title: string; children: React.ReactNode }) {
  return (
    <div className="space-y-2">
      <h2 className="text-sm text-muted-foreground px-2">{props.title}</h2>
      <div className="space-y-2">{props.children}</div>
    </div>
  );
}
