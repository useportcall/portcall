package paymentx

var stripeWebhookEvents = []string{
	"setup_intent.succeeded",
	"setup_intent.setup_failed",
	"payment_intent.succeeded",
	"payment_intent.payment_failed",
	"payment_intent.canceled",
	"payment_intent.requires_action",
	"charge.failed",
	"invoice.payment_failed",
	"invoice.payment_action_required",
	"payment_method.attached",
	"payment_method.detached",
	"payment_method.updated",
	"payment_method.card_automatically_updated",
}
