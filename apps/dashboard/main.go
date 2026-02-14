package main

import (
	"log"
	"os"

	"github.com/useportcall/portcall/apps/dashboard/internal/router"
	"github.com/useportcall/portcall/libs/go/cronx"
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

	stopScheduler, err := cronx.StartSubscriptionResetScheduler(q)
	if err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}
	defer func() {
		if err := stopScheduler(); err != nil {
			log.Printf("Scheduler shutdown failed: %v", err)
		}
	}()

	router, err := router.Init(db, crypto, q)
	if err != nil {
		log.Fatalf("failed to init router: %v", err)
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
