package router

import (
	"net/http"

	"github.com/useportcall/portcall/apps/dashboard/internal/middleware"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/account"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/address"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/app"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/app_config"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/checkout_session"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/company"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/connection"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/entitlement"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/feature"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/invoice"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/plan"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/plan_feature"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/plan_group"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/plan_item"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/secret"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/subscription"
	"github.com/useportcall/portcall/apps/dashboard/internal/modules/user"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/qx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func Init(db dbx.IORM, crypto cryptox.ICrypto, q qx.IQueue) routerx.IRouter {
	r := routerx.New(db, crypto, q)

	// --- Serve static files ---
	r.Use(middleware.StaticFile("./frontend/dist"))

	// --- Fallback to index.html for SPA routing ---
	r.NoRoute(func(c *routerx.Context) {
		c.Status(http.StatusOK)
		c.File("./frontend/dist/index.html")
	})

	r.GET("/ping", func(c *routerx.Context) { c.OK(map[string]any{"message": "pong"}) })

	r.Use(middleware.Auth(db))

	// App routes
	r.GET("/api/apps", app.ListApps)
	r.POST("/api/apps", app.CreateApp)
	r.GET("/api/apps/:app_id", app.GetApp)

	// Account routes
	r.GET("/api/apps/:app_id/account", account.GetAccount)

	// // Company routes
	r.GET("/api/apps/:app_id/company", company.GetCompany)
	r.POST("/api/apps/:app_id/company", company.UpsertCompany)

	// // AppConfig routes
	r.GET("/api/apps/:app_id/config", app_config.GetAppConfig)

	// // Address routes
	r.GET("/api/api/apps/:app_id/addresses", address.ListAddresses)
	r.POST("/api/apps/:app_id/addresses", address.CreateAddress)
	r.GET("/api/apps/:app_id/addresses/:id", address.GetAddress)
	r.POST("/api/apps/:app_id/addresses/:id", address.UpdateAddress)
	r.DELETE("/api/apps/:app_id/addresses/:id", address.DeleteAddress)

	// // Feature routes
	r.GET("/api/apps/:app_id/features", feature.ListFeatures)
	r.POST("/api/apps/:app_id/features", feature.CreateFeature)

	// // Secret routes
	r.GET("/api/apps/:app_id/secrets", secret.ListSecrets)
	r.POST("/api/apps/:app_id/secrets", secret.CreateSecret)
	r.POST("/api/apps/:app_id/secrets/:id/disable", secret.DisableSecret)

	// // User routes
	r.GET("/api/apps/:app_id/users", user.ListUsers)
	r.POST("/api/apps/:app_id/users", user.CreateUser)
	r.GET("/api/apps/:app_id/users/:id", user.GetUser)
	r.POST("/api/apps/:app_id/users/:id", user.UpdateUser)
	r.DELETE("/api/apps/:app_id/users/:id", user.DeleteUser)

	// // Subscription routes
	r.GET("/api/apps/:app_id/subscriptions", subscription.ListSubscriptions)
	r.GET("/api/apps/:app_id/users/:id/subscription", subscription.GetUserSubscription)

	// // Invoice routes
	r.GET("/api/apps/:app_id/invoices", invoice.ListInvoices)

	// // Connection routes
	r.GET("/api/apps/:app_id/connections", connection.ListConnections)
	r.POST("/api/apps/:app_id/connections", connection.CreateConnection)
	r.DELETE("/api/apps/:app_id/connections/:id", connection.DeleteConnection)

	// // Plan routes
	r.GET("/api/apps/:app_id/plans", plan.ListPlans)
	r.POST("/api/apps/:app_id/plans", plan.CreatePlan)
	r.GET("/api/apps/:app_id/plans/:id", plan.GetPlan)
	r.POST("/api/apps/:app_id/plans/:id", plan.UpdatePlan)
	r.DELETE("/api/apps/:app_id/plans/:id", plan.DeletePlan)
	r.POST("/api/apps/:app_id/plans/:id/publish", plan.PublishPlan)
	r.POST("/api/apps/:app_id/plans/:id/duplicate", plan.DuplicatePlan)

	// // PlanItem routes
	r.GET("/api/apps/:app_id/plan-items", plan_item.ListPlanItems)
	r.POST("/api/apps/:app_id/plan-items", plan_item.CreatePlanItem)
	r.POST("/api/apps/:app_id/plan-items/:id", plan_item.UpdatePlanItem)
	r.DELETE("/api/apps/:app_id/plan-items/:id", plan_item.DeletePlanItem)

	// // PlanFeature routes
	r.GET("/api/apps/:app_id/plan-features", plan_feature.ListPlanFeatures)
	r.POST("/api/apps/:app_id/plan-features", plan_feature.CreatePlanFeature)
	r.POST("/api/apps/:app_id/plan-features/:id", plan_feature.UpdatePlanFeature)
	r.DELETE("/api/apps/:app_id/plan-features/:id", plan_feature.DeletePlanFeature)

	// // PlanGroup routes
	r.GET("/api/apps/:app_id/groups", plan_group.ListPlanGroups)
	r.POST("/api/apps/:app_id/groups", plan_group.CreatePlanGroup)

	// // CheckoutSession routes
	r.POST("/api/apps/:app_id/checkout-sessions", checkout_session.CreateCheckoutSession)

	// // Entitlement routes
	r.GET("/api/apps/:app_id/entitlements", entitlement.ListEntitlements)

	return r
}
