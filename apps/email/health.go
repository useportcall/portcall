package main

import (
	"log"
	"net/http"
	"os"
)

func startHealthServer() {
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "9091"
		}
		mux := http.NewServeMux()
		mux.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})
		mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})
		log.Printf("email health server listening on :%s", port)
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			log.Printf("email health server error: %v", err)
		}
	}()
}
