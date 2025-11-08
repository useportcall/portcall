package qx

import (
	"encoding/json"
	"log"
	"os"

	"github.com/hibiken/asynq"
)

func New() IQueue {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR environment variable not set")
	}

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})

	return &queue_client{client: client}
}

type IQueue interface {
	Enqueue(name string, payload any, queue string) error
}

type queue_client struct {
	client *asynq.Client
}

func (c *queue_client) Enqueue(name string, payload any, queue string) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload: %v", err)
		return err
	}

	task := asynq.NewTask(name, payloadBytes)
	if _, err := c.client.Enqueue(task, asynq.Queue(queue)); err != nil {
		log.Printf("Failed to enqueue task %s: %v", name, err)
		return err
	}

	return nil
}
