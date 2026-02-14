# Follow-Up Policy Module

This package isolates payment follow-up decision logic from saga handlers.

## Invoice dunning policy

Use `DecideInvoiceFailure` to normalize:

- attempt/max-attempt values
- hard-decline behavior (`NoRetry`)
- follow-up stage (`invoice_first_reminder`, `invoice_retry_reminder`, `invoice_final_notice`)
- payment status (`past_due`, `uncollectible`)

`apps/billing/internal/sagas/on_payment` and `libs/go/services/payment` now use this policy.

## Quote unpaid policy

Use `DecideQuoteUnpaid` to determine quote follow-up stage based on accepted time:

- waiting
- reminder (default after 24h)
- final notice (default after 72h)

Thresholds are configurable per call via `QuoteUnpaidInput`.

## Why this exists

To keep follow-up rules provider-agnostic and editable in one place as billing
connections evolve (Stripe, Braintree, future providers).
