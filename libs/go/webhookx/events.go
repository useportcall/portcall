package webhookx

var relevantStripeEvents = map[string]bool{
	"setup_intent.succeeded":                    true,
	"setup_intent.setup_failed":                 true,
	"payment_intent.succeeded":                  true,
	"payment_intent.payment_failed":             true,
	"payment_intent.canceled":                   true,
	"payment_intent.requires_action":            true,
	"charge.failed":                             true,
	"invoice.payment_failed":                    true,
	"invoice.payment_action_required":           true,
	"payment_method.attached":                   true,
	"payment_method.detached":                   true,
	"payment_method.updated":                    true,
	"payment_method.card_automatically_updated": true,
}

var relevantBraintreeEvents = map[string]bool{
	"transaction_settled":                 true,
	"transaction_settlement_declined":     true,
	"subscription_charged_successfully":   true,
	"subscription_charged_unsuccessfully": true,
}
