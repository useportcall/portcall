package cronx

import (
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/useportcall/portcall/libs/go/qx"
)

// StartSubscriptionResetScheduler enqueues periodic subscription reset discovery jobs.
func StartSubscriptionResetScheduler(queue qx.IQueue) (func() error, error) {
	scheduler, err := gocron.NewScheduler(gocron.WithLocation(time.UTC))
	if err != nil {
		return nil, err
	}

	_, err = scheduler.NewJob(
		gocron.CronJob("*/10 * * * * *", true),
		gocron.NewTask(func() {
			log.Println("[cronx] enqueue find_subscriptions_to_reset")
			if err := queue.Enqueue("find_subscriptions_to_reset", map[string]any{}, "billing_queue"); err != nil {
				log.Printf("[cronx] enqueue failed: %v", err)
			}
		}),
	)
	if err != nil {
		return nil, err
	}

	scheduler.Start()
	return scheduler.Shutdown, nil
}
