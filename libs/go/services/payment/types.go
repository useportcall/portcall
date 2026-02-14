package payment

import (
	"github.com/stripe/stripe-go"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// PayInput is the input for paying an invoice.
type PayInput struct {
	InvoiceID uint `json:"invoice_id"`
}

// PayResult is the result of paying an invoice.
type PayResult struct {
	Invoice *models.Invoice
}

// ResolveInput is the input for resolving an invoice.
type ResolveInput struct {
	InvoiceID uint `json:"invoice_id"`
}

// ResolveResult is the result of resolving an invoice.
type ResolveResult struct {
	Invoice      *models.Invoice
	EmailPayload *InvoicePaidEmailPayload
}

// InvoicePaidEmailPayload is the payload for the invoice paid email task.
type InvoicePaidEmailPayload struct {
	InvoiceNumber  string `json:"invoice_number"`
	AmountPaid     string `json:"amount_paid"`
	DatePaid       string `json:"date_paid"`
	CompanyName    string `json:"company_name"`
	RecipientEmail string `json:"recipient_email"`
	LogoURL        string `json:"logo_url"`
	PaymentStatus  string `json:"payment_status"`
}

// StripeWebhookInput is the input for processing a Stripe webhook event.
type StripeWebhookInput struct {
	Event stripe.Event `json:"event"`
}

// StripeResult is the result of processing a Stripe webhook event.
type StripeResult struct {
	Action          string
	SessionID       string
	PaymentMethodID string
	Failure         *StripeFailurePayload
	Handled         bool
}

type BraintreeWebhookInput struct {
	Kind               string `json:"kind"`
	OrderID            string `json:"order_id"`
	PaymentMethodToken string `json:"payment_method_token"`
	FailureCount       int    `json:"failure_count"`
	FailureReason      string `json:"failure_reason"`
}

type BraintreeResult struct {
	Action          string
	SessionID       string
	PaymentMethodID string
	Failure         *StripeFailurePayload
	Handled         bool
}

// StripeFailurePayload captures a payment failure routed from Stripe webhooks.
type StripeFailurePayload struct {
	InvoiceID     uint   `json:"invoice_id"`
	Attempt       int    `json:"attempt"`
	NoRetry       bool   `json:"no_retry"`
	EventType     string `json:"event_type"`
	FailureReason string `json:"failure_reason"`
}

// CreateMethodInput is the input for creating a payment method.
type CreateMethodInput struct {
	AppID                   uint   `json:"app_id"`
	UserID                  uint   `json:"user_id"`
	PlanID                  uint   `json:"plan_id"`
	ExternalPaymentMethodID string `json:"external_payment_method_id"`
}

// CreateMethodResult is the result of creating a payment method.
type CreateMethodResult struct {
	PaymentMethod *models.PaymentMethod
	AppID         uint
	UserID        uint
	PlanID        uint
}
