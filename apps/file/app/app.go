// Package app exposes the file-service router for cross-module integration tests.
package app

import (
	"time"

	"github.com/useportcall/portcall/apps/file/handlers"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/ratelimitx"
	"github.com/useportcall/portcall/libs/go/routerx"
	"github.com/useportcall/portcall/libs/go/storex"
)

// NewRouter creates the file-service router with all routes.
// templateDir is the path to the templates directory.
func NewRouter(db dbx.IORM, store storex.IStore, templateDir string) routerx.IRouter {
	handlers.TemplateDir = templateDir

	r := routerx.New(db, nil, nil)
	r.SetStore(store)

	r.GET("/ping", func(c *routerx.Context) {
		c.OK(map[string]any{"status": "ok"})
	})

	rl := ratelimitx.NewGinAdapter()
	r.Use(func(c *routerx.Context) {
		if c.Request.URL.Path == "/ping" {
			c.Next()
			return
		}
		rl.Middleware(ratelimitx.ByClientIP, 120, time.Minute)(c.Context)
	})

	r.GET("/invoice/:invoice_id", handlers.GetInvoice)
	r.GET("/invoices/:invoice_id/view", handlers.ViewInvoice)
	r.GET("/icon-logos/:filename", handlers.GetIconLogo)

	return r
}