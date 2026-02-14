package webhookx

import (
	"net"
	"os"

	"github.com/useportcall/portcall/libs/go/routerx"
	"golang.org/x/time/rate"
)

type Router struct {
	defaultStore   *limiterStore
	stripeStore    *limiterStore
	braintreeStore *limiterStore
	stripeCIDRs    []*net.IPNet
}

func NewRouter() *Router {
	defaultRPS := loadPositiveInt("WEBHOOK_RPS", 10)
	defaultBurst := loadPositiveInt("WEBHOOK_BURST", 20)
	stripeRPS := loadPositiveInt("WEBHOOK_STRIPE_RPS", 100)
	stripeBurst := loadPositiveInt("WEBHOOK_STRIPE_BURST", 200)
	braintreeRPS := loadPositiveInt("WEBHOOK_BRAINTREE_RPS", 40)
	braintreeBurst := loadPositiveInt("WEBHOOK_BRAINTREE_BURST", 80)

	return &Router{
		defaultStore:   newLimiterStore(rate.Limit(defaultRPS), defaultBurst),
		stripeStore:    newLimiterStore(rate.Limit(stripeRPS), stripeBurst),
		braintreeStore: newLimiterStore(rate.Limit(braintreeRPS), braintreeBurst),
		stripeCIDRs:    parseCIDRList(os.Getenv("WEBHOOK_ALLOWED_IPS")),
	}
}

func RegisterRoutes(r routerx.IRouter) {
	NewRouter().Register(r)
}

func (w *Router) Register(r routerx.IRouter) {
	r.POST("/stripe/:connection_id", w.HandleStripeWebhook)
	r.GET("/braintree/:connection_id", w.HandleBraintreeChallenge)
	r.POST("/braintree/:connection_id", w.HandleBraintreeWebhook)
	r.POST("/postmark/:connection_id", w.HandlePostmarkWebhook)
}
