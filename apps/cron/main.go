package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/useportcall/portcall/libs/go/envx"
	"github.com/useportcall/portcall/libs/go/qx"
)

func main() {
	envx.Load()

	s, err := gocron.NewScheduler(
		gocron.WithLocation(time.UTC),
		// optional: WithLogger, WithMonitor, WithDistributedJobLocker…
	)
	if err != nil {
		log.Fatal(err)
	}

	q := qx.New()

	_, err = s.NewJob(
		gocron.CronJob("*/10 * * * * *", true),
		gocron.NewTask(func() {
			log.Println("Finding subscriptions to reset...")
			q.Enqueue("find_subscriptions_to_reset", map[string]interface{}{}, "billing_queue")
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	s.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down…")
	if err := s.Shutdown(); err != nil {
		log.Printf("scheduler shutdown error: %v", err)
	}
}
