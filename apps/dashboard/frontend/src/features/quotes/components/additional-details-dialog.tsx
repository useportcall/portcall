import { Section } from "@/components/sections";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Quote } from "@/models/quote";
import { TermsAndConditionsInput } from "./terms-and-conditions-input";

export function QuoteAdditionalDetailsDialog({ quote }: { quote: Quote }) {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline" size="sm" className="w-fit justify-start gap-2 text-left">
          Edit terms
        </Button>
      </DialogTrigger>
      <DialogContent>
        <div className="max-w-2xl w-full space-y-6">
          <DialogHeader>
            <DialogTitle>Add terms</DialogTitle>
          </DialogHeader>
          <Section title="Terms and conditions">
            <TermsAndConditionsInput quote={quote} />
          </Section>
        </div>
      </DialogContent>
    </Dialog>
  );
}
