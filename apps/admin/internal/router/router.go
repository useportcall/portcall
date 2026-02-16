package router

import (
	"net/http"

	"github.com/useportcall/portcall/apps/admin/internal/middleware"
	"github.com/useportcall/portcall/apps/admin/internal/modules/apikeys"
	"github.com/useportcall/portcall/apps/admin/internal/modules/apps"
	"github.com/useportcall/portcall/apps/admin/internal/modules/connections"
	"github.com/useportcall/portcall/apps/admin/internal/modules/dogfood"
	"github.com/useportcall/portcall/apps/admin/internal/modules/plans"
	"github.com/useportcall/portcall/apps/admin/internal/modules/queues"
	"github.com/useportcall/portcall/apps/admin/internal/modules/quotes"
	"github.com/useportcall/portcall/apps/admin/internal/modules/stats"
	"github.com/useportcall/portcall/apps/admin/internal/modules/subscriptions"
	"github.com/useportcall/portcall/apps/admin/internal/modules/users"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/qx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func Init(db dbx.IORM, crypto cryptox.ICrypto, q qx.IQueue) routerx.IRouter {
	r := routerx.New(db, crypto, q)

	r.GET("/ping", func(c *routerx.Context) { c.OK(map[string]any{"message": "pong"}) })
	r.GET("/healthz", func(c *routerx.Context) { c.OK(map[string]any{"status": "healthy"}) })

	// --- Serve static files ---
	r.Use(routerx.StaticFileMiddleware("./frontend/dist"))

	// Dogfood management routes (internal billing for Portcall dashboard)
	// These are placed before auth middleware for initial setup accessibility
	// Protected by admin API key header
	r.GET("/api/dogfood/status", dogfood.GetDogfoodStatus)
	r.POST("/api/dogfood/setup", dogfood.SetupDogfood)
	r.POST("/api/dogfood/refresh-k8s", dogfood.RefreshK8sSecrets)
	r.POST("/api/dogfood/regenerate-secrets", dogfood.RegenerateSecrets)
	r.POST("/api/dogfood/reset-password", dogfood.ResetDogfoodPassword)
	r.POST("/api/dogfood/sync", dogfood.SyncApps)
	r.POST("/api/dogfood/sync/:app_id", dogfood.SyncSingleApp)
	r.GET("/api/dogfood/apps/:app_id/users", dogfood.ListUsers)
	r.GET("/api/dogfood/apps/:app_id/users/:user_id", dogfood.GetUser)
	r.GET("/api/dogfood/apps/:app_id/features", dogfood.ListFeatures)
	r.POST("/api/dogfood/apps/:app_id/features", dogfood.CreateFeature)
	r.GET("/api/dogfood/apps/:app_id/plans", dogfood.ListPlans)
	r.POST("/api/dogfood/apps/:app_id/fix-users", dogfood.FixUsers)
	r.POST("/api/dogfood/setup-payment", dogfood.SetupPaymentConnection)
	r.POST("/api/dogfood/validate-users", dogfood.ValidateUsers)

	r.Use(middleware.Auth())

	// Stats routes
	r.GET("/api/stats", stats.GetStats)

	// Apps routes
	r.GET("/api/apps", apps.ListApps)
	r.GET("/api/apps/:app_id", apps.GetApp)
	r.DELETE("/api/apps/:app_id", apps.DeleteApp)

	// Connection management routes
	r.GET("/api/connections/check-mock", connections.CheckMockConnections)
	r.POST("/api/connections/clear-mock", connections.ClearMockConnections)
	r.GET("/api/apps/:app_id/connections", connections.ListConnectionsForApp)

	// API Key management routes
	r.POST("/api/apikeys/generate", apikeys.GenerateAdminAPIKey)
	r.GET("/api/apikeys/validate", apikeys.ValidateAPIKey)
	r.GET("/api/apikeys/info", apikeys.GetCurrentKeyInfo)

	// Users routes (per app)
	r.GET("/api/apps/:app_id/users", users.ListUsers)
	r.GET("/api/apps/:app_id/users/:user_id", users.GetUser)

	// Subscriptions routes (per app)
	r.GET("/api/apps/:app_id/subscriptions", subscriptions.ListSubscriptions)
	r.GET("/api/apps/:app_id/subscriptions/:subscription_id", subscriptions.GetSubscription)

	// Plans routes (per app)
	r.GET("/api/apps/:app_id/plans", plans.ListPlans)
	r.GET("/api/apps/:app_id/plans/:plan_id", plans.GetPlan)

	// Quotes routes (per app)
	r.GET("/api/apps/:app_id/quotes", quotes.ListQuotes)
	r.GET("/api/apps/:app_id/quotes/:quote_id", quotes.GetQuote)

	// Queue management routes
	r.GET("/api/queues", queues.GetQueues)
	r.POST("/api/queues/enqueue", queues.EnqueueTask)

	// Queue inspection routes (new)
	r.GET("/api/queues/stats", queues.GetQueueStats)
	r.GET("/api/queues/tasks", queues.ListTasks)
	r.GET("/api/queues/tasks/:task_id", queues.GetTask)
	r.POST("/api/queues/tasks/archive", queues.ArchiveTask)
	r.POST("/api/queues/tasks/delete", queues.DeleteTask)
	r.POST("/api/queues/tasks/run", queues.RunTask)
	r.POST("/api/queues/tasks/retry", queues.RetryWithModifiedPayload)

	// --- Fallback to index.html for SPA routing ---
	// Must be at the end after all API routes are registered
	r.NoRoute(func(c *routerx.Context) {
		c.Status(http.StatusOK)
		c.File("./frontend/dist/index.html")
	})

	return r
}
