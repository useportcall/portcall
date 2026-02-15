package main

import (
	"log"
	"os"
	"strings"

	"github.com/useportcall/portcall/apps/email/internal/smtprelay"
	"github.com/useportcall/portcall/libs/go/envx"
)

func main() {
	envx.Load()
	startHealthServer()

	workerErr := make(chan error, 1)
	go func() {
		workerErr <- runEmailWorker()
	}()

	var relayErr <-chan error
	if relayEnabled() {
		ch := make(chan error, 1)
		relayErr = ch
		go func() {
			ch <- smtprelay.RunFromEnv()
		}()
	} else {
		log.Printf("SMTP relay disabled (SMTP_RELAY_ENABLED=false)")
	}

	select {
	case err := <-workerErr:
		log.Fatalf("email worker failed: %v", err)
	case err := <-relayErr:
		log.Fatalf("smtp relay failed: %v", err)
	}
}

func relayEnabled() bool {
	return !strings.EqualFold(os.Getenv("SMTP_RELAY_ENABLED"), "false")
}
