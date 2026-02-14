package main

import (
	"log"
	"net/http"
	"os"

	"github.com/useportcall/portcall/apps/billing/internal/sagas"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/envx"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

func main() {
	envx.Load()

	db, err := dbx.New()
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}
	crypto, err := cryptox.New()
	if err != nil {
		log.Fatalf("failed to init crypto: %v", err)
	}

	srv, err := server.New(db, crypto, map[string]int{
		"billing_queue": 10,
	})
	if err != nil {
		log.Fatalf("failed to init worker server: %v", err)
	}

	// lightweight health endpoint for orchestration/health checks
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "9090"
		}
		mux := http.NewServeMux()
		mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})
		mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})
		log.Printf("billing health server listening on :%s", port)
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			log.Printf("billing health server error: %v", err)
		}
	}()

	sagas.RegisterAll(srv)

	if err := srv.R(); err != nil {
		log.Fatalf("billing worker failed: %v", err)
	}
}
