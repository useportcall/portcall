package main

import (
	"log"
	"os"

	quoteapp "github.com/useportcall/portcall/apps/quote/app"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/envx"
	"github.com/useportcall/portcall/libs/go/qx"
	"github.com/useportcall/portcall/libs/go/storex"
)

func main() {
	envx.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8010"
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
	store, err := storex.New()
	if err != nil {
		log.Fatalf("failed to init store: %v", err)
	}

	r := quoteapp.NewRouter(db, crypto, q, store, "templates/*.html")
	r.Run(":" + port)
}
