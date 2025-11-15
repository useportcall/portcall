package main

import (
	"log"
	"os"

	"github.com/useportcall/portcall/apps/file/handlers"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/envx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func main() {
	envx.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	db := dbx.New()

	r := routerx.New(db, nil, nil)

	r.GET("/invoice/:invoice_id", handlers.GetInvoice)
	r.GET("/invoices/:invoice_id/view", handlers.ViewInvoice)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
