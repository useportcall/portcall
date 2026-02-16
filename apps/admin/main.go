package main

import (
	"log"
	"os"

	"github.com/useportcall/portcall/apps/admin/internal/router"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/envx"
	"github.com/useportcall/portcall/libs/go/qx"
)

func main() {
	envx.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	db, err := dbx.New()
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}
	crypto, err := cryptox.New()
	if err != nil {
		log.Fatalf("failed to init crypto: %v", err)
	}
	q, err := qx.New()
	if err != nil {
		log.Fatalf("failed to init queue: %v", err)
	}

	router := router.Init(db, crypto, q)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
