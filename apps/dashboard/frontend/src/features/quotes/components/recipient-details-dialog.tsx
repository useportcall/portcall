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
import { hasRecipientDetails } from "../utils/quote-ui";
import { RecipientDetailsFields } from "./recipient-details-fields";

export function RecipientDetailsDialog({ quote }: { quote: Quote }) {
  const buttonLabel = hasRecipientDetails(quote)
    ? "Edit recipient details"
    : "Add recipient details";

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline" size="sm" className="w-fit justify-start gap-2 text-left">
          {buttonLabel}
        </Button>
      </DialogTrigger>
      <DialogContent>
        <div className="max-w-2xl w-full space-y-6">
          <DialogHeader>
            <DialogTitle className="text-start">Recipient details</DialogTitle>
          </DialogHeader>
          <Section title="Basic details">
            <div className="px-2">
              <RecipientDetailsFields quote={quote} />
            </div>
          </Section>
        </div>
      </DialogContent>
    </Dialog>
  );
}
