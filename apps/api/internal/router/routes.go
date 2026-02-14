package router

import (
	"strings"
	"time"

	"github.com/useportcall/portcall/apps/api/internal/modules/v1/checkout_session"
	"github.com/useportcall/portcall/apps/api/internal/modules/v1/entitlement"
	"github.com/useportcall/portcall/apps/api/internal/modules/v1/feature"
	"github.com/useportcall/portcall/apps/api/internal/modules/v1/invoice"
	"github.com/useportcall/portcall/apps/api/internal/modules/v1/meter_event"
	"github.com/useportcall/portcall/apps/api/internal/modules/v1/payment_link"
	"github.com/useportcall/portcall/apps/api/internal/modules/v1/plan"
	"github.com/useportcall/portcall/apps/api/internal/modules/v1/plan_feature"
	"github.com/useportcall/portcall/apps/api/internal/modules/v1/subscription"
	"github.com/useportcall/portcall/apps/api/internal/modules/v1/user"
	"github.com/useportcall/portcall/libs/go/ratelimitx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func rateLimitMiddleware(c *routerx.Context) {
	if c.Request.URL.Path == "/ping" || c.Request.URL.Path == "/healthz" ||
		strings.HasPrefix(c.Request.URL.Path, "/docs") {
		c.Next()
		return
	}
	limiter := ratelimitx.NewGinAdapter()
	middleware := limiter.Middleware(ratelimitx.ByAPIKeyHeader, 1000, 1*time.Hour)
	middleware(c.Context)
}

func registerRoutes(r routerx.IRouter) {
	r.GET("/v1/users", user.ListUsers)
	r.POST("/v1/users", user.CreateUser)
	r.GET("/v1/users/:id", user.GetUser)
	r.POST("/v1/users/:id", user.UpdateUser)
	r.GET("/v1/users/:id/subscription", user.GetUserSubscription)
	r.POST("/v1/users/:id/billing-address", user.UpsertBillingAddress)

	r.GET("/v1/entitlements", entitlement.ListEntitlements)
	r.GET("/v1/entitlements/:user_id/:id", entitlement.GetEntitlement)

	r.GET("/v1/subscriptions", subscription.ListSubscriptions)
	r.POST("/v1/subscriptions", subscription.CreateSubscription)
	r.GET("/v1/subscriptions/:id", subscription.GetSubscription)
	r.POST("/v1/subscriptions/:subscription_id", subscription.UpdateSubscription)
	r.POST("/v1/subscriptions/:subscription_id/cancel", subscription.CancelSubscription)

	r.GET("/v1/invoices", invoice.ListInvoices)

	r.POST("/v1/meter-events", meter_event.CreateMeterEvent)

	r.POST("/v1/checkout-sessions", checkout_session.CreateCheckoutSession)
	r.POST("/v1/payment-links", payment_link.CreatePaymentLink)

	r.GET("/v1/plans", plan.ListPlans)
	r.POST("/v1/plans", plan.CreatePlan)
	r.GET("/v1/plans/:id", plan.GetPlan)
	r.POST("/v1/plans/:id", plan.UpdatePlan)
	r.POST("/v1/plans/:id/publish", plan.PublishPlan)

	r.GET("/v1/features", feature.ListFeatures)
	r.POST("/v1/features", feature.CreateFeature)

	r.GET("/v1/plan-features", plan_feature.ListPlanFeatures)
	r.POST("/v1/plan-features", plan_feature.CreatePlanFeature)
}
