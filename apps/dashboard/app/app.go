// Package app exposes the dashboard router for cross-module integration tests.
package app

import (
	"github.com/useportcall/portcall/apps/dashboard/internal/router"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/qx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

// NewRouter creates the dashboard router with all routes and middleware.
func NewRouter(db dbx.IORM, crypto cryptox.ICrypto, q qx.IQueue) (routerx.IRouter, error) {
	return router.Init(db, crypto, q)
}
