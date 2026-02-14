import { Quote } from "@/models/quote";
import { UpdateAddressInput } from "./update-address-input";
import { UpdateQuoteInput } from "./update-quote-input";

export function RecipientDetailsFields({ quote }: { quote: Quote }) {
  return (
    <>
      <UpdateQuoteInput k="recipient_email" quote={quote} placeholder="Email" />
      <UpdateQuoteInput k="recipient_name" quote={quote} placeholder="Full name" />
      <UpdateQuoteInput
        k="recipient_title"
        quote={quote}
        placeholder="Title (e.g. CEO)"
      />
      <UpdateQuoteInput k="company_name" quote={quote} placeholder="Company name" />
      <UpdateAddressInput
        k="line1"
        address={quote.recipient_address}
        placeholder="Address line 1"
      />
      <UpdateAddressInput
        k="line2"
        address={quote.recipient_address}
        placeholder="Address line 2"
      />
      <div className="flex w-full gap-2">
        <UpdateAddressInput k="city" address={quote.recipient_address} placeholder="City" />
        <UpdateAddressInput k="state" address={quote.recipient_address} placeholder="State" />
      </div>
      <div className="flex w-full gap-2">
        <UpdateAddressInput
          k="postal_code"
          address={quote.recipient_address}
          placeholder="Postal code"
        />
        <UpdateAddressInput
          k="country"
          address={quote.recipient_address}
          placeholder="Country"
        />
      </div>
    </>
  );
}
