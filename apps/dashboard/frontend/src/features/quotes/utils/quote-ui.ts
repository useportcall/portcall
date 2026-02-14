import { Quote } from "@/models/quote";

const SEND_DISABLED_STATUSES: Quote["status"][] = [
  "accepted",
  "voided",
  "rejected",
];
const VOID_DISABLED_STATUSES: Quote["status"][] = [
  "accepted",
  "voided",
  "rejected",
  "draft",
];
const READONLY_STATUSES: Quote["status"][] = ["sent", "accepted", "voided", "rejected"];

export function getSendButtonTitle(status: Quote["status"]) {
  return status === "sent" ? "Resend" : "Send";
}

export function isSendDisabled(status: Quote["status"]) {
  return SEND_DISABLED_STATUSES.includes(status);
}

export function isVoidDisabled(status: Quote["status"]) {
  return VOID_DISABLED_STATUSES.includes(status);
}

export function isQuoteReadonly(status: Quote["status"]) {
  return READONLY_STATUSES.includes(status);
}

export function buildRecipientDetails(quote: Quote) {
  return [
    quote.recipient_email,
    quote.recipient_name,
    quote.recipient_title,
    quote.company_name,
    quote.recipient_address.line1,
    quote.recipient_address.city,
    quote.recipient_address.postal_code,
  ].filter((value): value is string => Boolean(value));
}

export function hasRecipientDetails(quote: Quote) {
  const { recipient_address } = quote;
  const hasAddress = [
    recipient_address.line1,
    recipient_address.line2,
    recipient_address.city,
    recipient_address.state,
    recipient_address.postal_code,
    recipient_address.country,
  ].some(Boolean);

  return (
    hasAddress ||
    Boolean(quote.recipient_email) ||
    Boolean(quote.recipient_name) ||
    Boolean(quote.company_name)
  );
}

export function summarizeTerms(text: string, maxLength = 50) {
  if (text.length <= maxLength) {
    return text;
  }
  return `${text.slice(0, maxLength)}...`;
}
