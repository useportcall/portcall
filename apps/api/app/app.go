// Package app exposes the public API router for cross-module integration tests.
package app

import (
	"github.com/useportcall/portcall/apps/api/internal/router"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/qx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

// NewRouter creates the public API router with API-key auth and all routes.
func NewRouter(db dbx.IORM, crypto cryptox.ICrypto, q qx.IQueue) routerx.IRouter {
	return router.Init(db, crypto, q)
}
