import { Section } from "@/components/sections";
import { Quote } from "@/models/quote";
import { BillInAdvanceSelect } from "../../plans/components/bill-in-advance-select";
import { FreeTrialComboBox } from "../../plans/components/free-trial-combo-box";
import MutableEntitlements from "../../plans/components/mutable-entitlements";
import { PlanCurrencySelect } from "../../plans/components/plan-currency-select";
import { PlanFixedUnitAmountInput } from "../../plans/components/plan-fixed-unit-amount-input";
import { PlanGroupSelect } from "../../plans/components/plan-group-select";
import { PlanIntervalSelect } from "../../plans/components/plan-interval-select";
import { ProrationSelect } from "../../plans/components/proration-select";
import { CheckoutRequiredToggle, QuoteExpiryPicker } from "./quote-selects";
import { PreparedByEmailInput } from "./quote-inputs";
import { QuoteAdditionalDetailsDialog, RecipientDetailsDialog } from "./quote-dialogs";
import { RecipientDetailsSummary } from "./recipient-details-summary";
import { TermsSummary } from "./terms-summary";
import { AnimatedSuspense } from "./animated-suspense";

export function QuoteViewSidebar({ quote }: { quote: Quote }) {
  return (
    <div className="flex flex-col gap-6 px-2 animate-fade-in">
      <AnimatedSuspense>
        <Section title="Recipient details">
          <RecipientDetailsSummary quote={quote} />
          <RecipientDetailsDialog quote={quote} />
        </Section>
      </AnimatedSuspense>
      <AnimatedSuspense>
        <Section title="Prepared by">
          <PreparedByEmailInput quote={quote} />
        </Section>
      </AnimatedSuspense>
      <Section title="Additional settings">
        <QuoteExpiryPicker quote={quote} />
        <CheckoutRequiredToggle quote={quote} />
        <TermsSummary quote={quote} />
        <QuoteAdditionalDetailsDialog quote={quote} />
      </Section>
      <Section title="Fixed price">
        <PlanFixedUnitAmountInput id={quote.plan.id} />
        <PlanCurrencySelect id={quote.plan.id} />
        <PlanIntervalSelect id={quote.plan.id} />
        <PlanGroupSelect id={quote.plan.id} />
      </Section>
      <MutableEntitlements id={quote.plan.id} />
      <Section title="Other properties">
        <FreeTrialComboBox id={quote.plan.id} />
        <BillInAdvanceSelect />
        <ProrationSelect />
      </Section>
    </div>
  );
}
