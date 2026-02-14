package router

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/useportcall/portcall/apps/dashboard/internal/middleware"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/account"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/address"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/app"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/app_config"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/billing"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/billing_meter"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/checkout_session"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/company"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/config"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/connection"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/entitlement"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/feature"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/invoice"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/payment_link"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/plan"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/plan_feature"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/plan_group"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/plan_item"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/quote"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/secret"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/subscription"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/user"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/qx"
	"github.com/useportcall/portcall/libs/go/ratelimitx"
	"github.com/useportcall/portcall/libs/go/routerx"
	"github.com/useportcall/portcall/libs/go/storex"
	"github.com/useportcall/portcall/libs/go/webhookx"
)

func Init(db dbx.IORM, crypto cryptox.ICrypto, q qx.IQueue) (routerx.IRouter, error) {
	r := routerx.New(db, crypto, q)
	store, err := storex.New()
	if err != nil {
		return nil, fmt.Errorf("init store: %w", err)
	}
	r.SetStore(store)

	r.GET("/ping", func(c *routerx.Context) { c.OK(map[string]any{"message": "pong"}) })
	r.GET("/health", func(c *routerx.Context) { c.OK(map[string]any{"status": "healthy"}) })
	r.GET("/healthz", func(c *routerx.Context) { c.OK(map[string]any{"status": "healthy"}) })

	// Config endpoint for runtime configuration (Keycloak URL, etc.)
	r.GET("/api/config", config.GetConfig)

	// --- Serve static files ---
	staticDir := "./frontend/dist"
	if v := os.Getenv("DASHBOARD_STATIC_DIR"); v != "" {
		staticDir = v
	}
	r.Use(routerx.StaticFileMiddleware(staticDir))

	// --- Fallback to index.html for SPA routing ---
	r.NoRoute(func(c *routerx.Context) {
		c.Status(http.StatusOK)
		c.File(staticDir + "/index.html")
	})

	// Public webhook endpoints (no dashboard auth).
	webhookx.RegisterRoutes(r)

	r.Use(middleware.Auth(db))

	// Rate limiting middleware - 100 requests per minute per user
	rateLimiter := ratelimitx.NewGinAdapter()
	e2eMode := os.Getenv("E2E_MODE") == "true"
	r.Use(func(c *routerx.Context) {
		if e2eMode {
			c.Next()
			return
		}

		// Only rate limit API endpoints (not static files or healthz)
		if !strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.Next()
			return
		}

		// Apply rate limiting by user email
		middleware := rateLimiter.Middleware(
			ratelimitx.ByEmail,
			100,
			1*time.Minute,
		)
		middleware(c.Context)
	})

	// Account routes
	r.GET("/api/account", account.GetAccount)

	// App routes
	r.GET("/api/apps", app.ListApps)
	r.POST("/api/apps", app.CreateApp)
	r.GET("/api/apps/:app_id", app.GetApp)

	// Company routes
	r.GET("/api/apps/:app_id/company", company.GetCompany)
	r.POST("/api/apps/:app_id/company", company.UpsertCompany)

	// AppConfig routes
	r.GET("/api/apps/:app_id/config", app_config.GetAppConfig)
	r.POST("/api/apps/:app_id/config", app_config.UpdateAppConfig)

	// Address routes
	r.GET("/api/api/apps/:app_id/addresses", address.ListAddresses)
	r.POST("/api/apps/:app_id/addresses", address.CreateAddress)
	r.GET("/api/apps/:app_id/addresses/:id", address.GetAddress)
	r.POST("/api/apps/:app_id/addresses/:id", address.UpdateAddress)
	r.DELETE("/api/apps/:app_id/addresses/:id", address.DeleteAddress)

	// Feature routes
	r.GET("/api/apps/:app_id/features", feature.ListFeatures)
	r.POST("/api/apps/:app_id/features", feature.CreateFeature)

	// Secret routes
	r.GET("/api/apps/:app_id/secrets", secret.ListSecrets)
	r.POST("/api/apps/:app_id/secrets", secret.CreateSecret)
	r.POST("/api/apps/:app_id/secrets/:id/disable", secret.DisableSecret)

	// User routes
	r.GET("/api/apps/:app_id/users", user.ListUsers)
	r.POST("/api/apps/:app_id/users", user.CreateUser)
	r.GET("/api/apps/:app_id/users/:id", user.GetUser)
	r.POST("/api/apps/:app_id/users/:id", user.UpdateUser)
	r.DELETE("/api/apps/:app_id/users/:id", user.DeleteUser)

	// Subscription routes
	r.GET("/api/apps/:app_id/subscriptions", subscription.ListSubscriptions)
	r.POST("/api/apps/:app_id/subscriptions", subscription.CreateSubscription)
	r.GET("/api/apps/:app_id/subscriptions/:id", subscription.GetSubscription)
	r.GET("/api/apps/:app_id/users/:id/subscription", subscription.GetUserSubscription)
	r.POST("/api/apps/:app_id/subscriptions/:subscription_id/cancel", subscription.CancelSubscription)

	// Invoice routes
	r.GET("/api/apps/:app_id/invoices", invoice.ListInvoices)

	// Connection routes
	r.GET("/api/apps/:app_id/connections", connection.ListConnections)
	r.POST("/api/apps/:app_id/connections", connection.CreateConnection)
	r.DELETE("/api/apps/:app_id/connections/:id", connection.DeleteConnection)

	// Plan routes
	r.GET("/api/apps/:app_id/plans", plan.ListPlans)
	r.POST("/api/apps/:app_id/plans", plan.CreatePlan)
	r.GET("/api/apps/:app_id/plans/:id", plan.GetPlan)
	r.POST("/api/apps/:app_id/plans/:id", plan.UpdatePlan)
	r.DELETE("/api/apps/:app_id/plans/:id", plan.DeletePlan)
	r.POST("/api/apps/:app_id/plans/:id/publish", plan.PublishPlan)
	r.POST("/api/apps/:app_id/plans/:id/duplicate", plan.DuplicatePlan)
	r.POST("/api/apps/:app_id/plans/:id/copy", plan.CopyPlan)

	// PlanItem routes
	r.GET("/api/apps/:app_id/plan-items", plan_item.ListPlanItems)
	r.POST("/api/apps/:app_id/plan-items", plan_item.CreatePlanItem)
	r.POST("/api/apps/:app_id/plan-items/:id", plan_item.UpdatePlanItem)
	r.DELETE("/api/apps/:app_id/plan-items/:id", plan_item.DeletePlanItem)

	// PlanFeature routes
	r.GET("/api/apps/:app_id/plan-features", plan_feature.ListPlanFeatures)
	r.POST("/api/apps/:app_id/plan-features", plan_feature.CreatePlanFeature)
	r.POST("/api/apps/:app_id/plan-features/:id", plan_feature.UpdatePlanFeature)
	r.DELETE("/api/apps/:app_id/plan-features/:id", plan_feature.DeletePlanFeature)

	// PlanGroup routes
	r.GET("/api/apps/:app_id/groups", plan_group.ListPlanGroups)
	r.POST("/api/apps/:app_id/groups", plan_group.CreatePlanGroup)

	// Quote routes
	r.GET("/api/apps/:app_id/quotes", quote.ListQuotes)
	r.POST("/api/apps/:app_id/quotes", quote.CreateQuote)
	r.GET("/api/apps/:app_id/quotes/:id", quote.GetQuote)
	r.GET("/api/apps/:app_id/quotes/:id/signature", quote.GetQuoteSignature)
	r.POST("/api/apps/:app_id/quotes/:id", quote.UpdateQuote)
	r.POST("/api/apps/:app_id/quotes/:id/send", quote.SendQuote)
	r.POST("/api/apps/:app_id/quotes/:id/void", quote.VoidQuote)

	// CheckoutSession routes
	r.POST("/api/apps/:app_id/checkout-sessions", checkout_session.CreateCheckoutSession)
	r.POST("/api/apps/:app_id/payment-links", payment_link.CreatePaymentLink)

	// Entitlement routes
	r.GET("/api/apps/:app_id/entitlements", entitlement.ListEntitlements)
	r.POST("/api/apps/:app_id/entitlements/:user_id", entitlement.CreateEntitlement)
	r.POST("/api/apps/:app_id/entitlements/:user_id/:feature_id", entitlement.UpdateEntitlement)
	r.DELETE("/api/apps/:app_id/entitlements/:user_id/:feature_id", entitlement.DeleteEntitlement)
	r.POST("/api/apps/:app_id/entitlements/:user_id/toggle", entitlement.ToggleEntitlement)

	// Billing Meter routes (for tracking metered usage and billing projections)
	r.GET("/api/apps/:app_id/billing-meters", billing_meter.ListBillingMeters)
	r.GET("/api/apps/:app_id/billing-meters/:subscription_id/:feature_id", billing_meter.GetBillingMeter)
	r.POST("/api/apps/:app_id/billing-meters/:subscription_id/:feature_id", billing_meter.UpdateBillingMeter)
	r.POST("/api/apps/:app_id/billing-meters/:subscription_id/:feature_id/reset", billing_meter.ResetBillingMeter)

	// Billing routes (for dogfood billing)
	// ENTITLEMENT GATES: These endpoints are used to check usage limits
	r.GET("/api/apps/:app_id/billing/quota", billing.CheckSubscriptionQuota)
	r.GET("/api/apps/:app_id/billing/quota/users", billing.CheckUserQuota)
	r.GET("/api/apps/:app_id/billing/quota/subscriptions", billing.CheckSubscriptionQuota)
	r.GET("/api/apps/:app_id/billing/subscription", billing.GetSubscriptionInfo)
	r.POST("/api/apps/:app_id/billing/upgrade-to-pro", billing.UpgradeToPro)
	r.POST("/api/apps/:app_id/billing/downgrade-to-free", billing.DowngradeToFree)
	r.GET("/api/apps/:app_id/billing/invoices", billing.ListBillingInvoices)
	r.GET("/api/apps/:app_id/billing/address", billing.GetBillingAddress)
	r.POST("/api/apps/:app_id/billing/address", billing.UpsertBillingAddressHandler)

	return r, nil
}
