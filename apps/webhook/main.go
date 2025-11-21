package main

import (
	"log"
	"os"

	"github.com/useportcall/portcall/apps/webhook/internal/handlers"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/envx"
	"github.com/useportcall/portcall/libs/go/qx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func main() {
	envx.Load()

	port, exists := os.LookupEnv("PORT")
	if !exists {
		log.Fatal("PORT environment variable not set")
	}

	db := dbx.New()
	queue := qx.New()
	crypto := cryptox.New()

	r := routerx.New(db, crypto, queue)

	r.POST("/stripe/:connection_id", handlers.HandleStripeWebhook)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
