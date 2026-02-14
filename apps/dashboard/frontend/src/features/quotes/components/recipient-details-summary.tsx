import { Quote } from "@/models/quote";
import { buildRecipientDetails } from "../utils/quote-ui";

export function RecipientDetailsSummary({ quote }: { quote: Quote }) {
  const details = buildRecipientDetails(quote);

  if (details.length === 0) {
    return (
      <div className="px-2 text-sm text-muted-foreground">
        <p>No recipient details added yet.</p>
      </div>
    );
  }

  return (
    <div className="px-2 text-sm text-muted-foreground">
      {details.map((detail) => (
        <p key={detail}>{detail}</p>
      ))}
    </div>
  );
}
