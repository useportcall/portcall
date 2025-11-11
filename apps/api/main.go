package main

import (
	"log"
	"os"

	"github.com/useportcall/portcall/apps/api/internal/modules/v1/checkout_session"
	"github.com/useportcall/portcall/apps/api/internal/modules/v1/user"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/envx"
	"github.com/useportcall/portcall/libs/go/qx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func main() {
	envx.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	db := dbx.New()
	crypto := cryptox.New()
	q := qx.New()

	r := routerx.New(db, crypto, q)

	r.GET("/v1/users/", user.ListUsers)

	r.POST("/v1/checkout-sessions", checkout_session.CreateCheckoutSession)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
