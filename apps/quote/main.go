package main

import (
	"github.com/useportcall/portcall/apps/quote/internal/handlers"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/envx"
	"github.com/useportcall/portcall/libs/go/qx"
	"github.com/useportcall/portcall/libs/go/routerx"
	"github.com/useportcall/portcall/libs/go/storex"
)

func main() {
	envx.Load()

	r := routerx.New(dbx.New(), cryptox.New(), qx.New())
	r.SetStore(storex.New())

	// Load HTML templates from the templates directory
	r.LoadHTMLGlob("templates/*.html")

	r.GET("/quotes/:id", handlers.GetQuote)
	r.POST("/quotes/:id", handlers.SubmitQuote)
	r.GET("/quotes/:id/success", handlers.QuoteSuccess)

	r.Run(":8095")
}
