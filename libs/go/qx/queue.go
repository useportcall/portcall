package qx

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/hibiken/asynq"
)

func New() (IQueue, error) { return NewFromEnv() }

func NewFromEnv() (IQueue, error) {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		return nil, fmt.Errorf("REDIS_ADDR environment variable not set")
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisTLS := os.Getenv("REDIS_TLS")

	var redisOpt asynq.RedisClientOpt
	if redisTLS == "true" {
		redisOpt = asynq.RedisClientOpt{
			Addr:     redisAddr,
			Password: redisPassword,
			TLSConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		}
	} else if redisPassword != "" {
		redisOpt = asynq.RedisClientOpt{
			Addr:     redisAddr,
			Password: redisPassword,
		}
	} else {
		redisOpt = asynq.RedisClientOpt{
			Addr: redisAddr,
		}
	}

	client := asynq.NewClient(redisOpt)
	inspector := asynq.NewInspector(redisOpt)

	return &queue_client{client: client, inspector: inspector}, nil
}

// TaskInfo represents a queue task with all its metadata
type TaskInfo struct {
	ID            string `json:"id"`
	Queue         string `json:"queue"`
	Type          string `json:"type"`
	Payload       string `json:"payload"`
	State         string `json:"state"`
	MaxRetry      int    `json:"max_retry"`
	Retried       int    `json:"retried"`
	LastErr       string `json:"last_err"`
	LastFailedAt  string `json:"last_failed_at,omitempty"`
	NextProcessAt string `json:"next_process_at,omitempty"`
	CompletedAt   string `json:"completed_at,omitempty"`
}

// QueueStats represents queue statistics
type QueueStats struct {
	Queue     string `json:"queue"`
	Active    int    `json:"active"`
	Pending   int    `json:"pending"`
	Scheduled int    `json:"scheduled"`
	Retry     int    `json:"retry"`
	Archived  int    `json:"archived"`
	Completed int    `json:"completed"`
}

type IQueue interface {
	Enqueue(name string, payload any, queue string) error
	Close() error
	// Inspector methods
	ListQueues() ([]string, error)
	GetQueueStats(queue string) (*QueueStats, error)
	ListTasks(queue string, state string, limit int, taskType string) ([]*TaskInfo, error)
	ArchiveTask(queue string, taskID string) error
	DeleteTask(queue string, taskID string) error
	RunTask(queue string, taskID string) error
	GetTaskInfo(queue string, taskID string) (*TaskInfo, error)
}

type queue_client struct {
	client    *asynq.Client
	inspector *asynq.Inspector
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

func (c *queue_client) Close() error {
	c.inspector.Close()
	return c.client.Close()
}

// ListQueues returns all available queues
func (c *queue_client) ListQueues() ([]string, error) {
	return c.inspector.Queues()
}

// GetQueueStats returns statistics for a queue
func (c *queue_client) GetQueueStats(queue string) (*QueueStats, error) {
	info, err := c.inspector.GetQueueInfo(queue)
	if err != nil {
		return nil, err
	}
	return &QueueStats{
		Queue:     queue,
		Active:    info.Active,
		Pending:   info.Pending,
		Scheduled: info.Scheduled,
		Retry:     info.Retry,
		Archived:  info.Archived,
		Completed: info.Completed,
	}, nil
}

// convertTaskInfo converts asynq.TaskInfo to our TaskInfo
func convertTaskInfo(t *asynq.TaskInfo) *TaskInfo {
	info := &TaskInfo{
		ID:       t.ID,
		Queue:    t.Queue,
		Type:     t.Type,
		Payload:  string(t.Payload),
		State:    t.State.String(),
		MaxRetry: t.MaxRetry,
		Retried:  t.Retried,
		LastErr:  t.LastErr,
	}
	if !t.LastFailedAt.IsZero() {
		info.LastFailedAt = t.LastFailedAt.Format("2006-01-02T15:04:05Z")
	}
	if !t.NextProcessAt.IsZero() {
		info.NextProcessAt = t.NextProcessAt.Format("2006-01-02T15:04:05Z")
	}
	if !t.CompletedAt.IsZero() {
		info.CompletedAt = t.CompletedAt.Format("2006-01-02T15:04:05Z")
	}
	return info
}

// ListTasks returns tasks from a queue filtered by state
func (c *queue_client) ListTasks(queue string, state string, limit int, taskType string) ([]*TaskInfo, error) {
	opts := []asynq.ListOption{asynq.PageSize(limit)}
	var tasks []*asynq.TaskInfo
	var err error

	switch state {
	case "active":
		tasks, err = c.inspector.ListActiveTasks(queue, opts...)
	case "pending":
		tasks, err = c.inspector.ListPendingTasks(queue, opts...)
	case "scheduled":
		tasks, err = c.inspector.ListScheduledTasks(queue, opts...)
	case "retry":
		tasks, err = c.inspector.ListRetryTasks(queue, opts...)
	case "archived":
		tasks, err = c.inspector.ListArchivedTasks(queue, opts...)
	case "completed":
		tasks, err = c.inspector.ListCompletedTasks(queue, opts...)
	default:
		// Get all states and merge
		allTasks := make([]*asynq.TaskInfo, 0)

		active, _ := c.inspector.ListActiveTasks(queue, opts...)
		allTasks = append(allTasks, active...)

		pending, _ := c.inspector.ListPendingTasks(queue, opts...)
		allTasks = append(allTasks, pending...)

		scheduled, _ := c.inspector.ListScheduledTasks(queue, opts...)
		allTasks = append(allTasks, scheduled...)

		retry, _ := c.inspector.ListRetryTasks(queue, opts...)
		allTasks = append(allTasks, retry...)

		archived, _ := c.inspector.ListArchivedTasks(queue, opts...)
		allTasks = append(allTasks, archived...)

		completed, _ := c.inspector.ListCompletedTasks(queue, opts...)
		allTasks = append(allTasks, completed...)

		tasks = allTasks
	}

	if err != nil {
		return nil, err
	}

	result := make([]*TaskInfo, 0, len(tasks))
	for _, t := range tasks {
		// Filter by task type if specified
		if taskType != "" && t.Type != taskType {
			continue
		}
		result = append(result, convertTaskInfo(t))
	}

	// Limit the results
	if len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}

// ArchiveTask archives a task
func (c *queue_client) ArchiveTask(queue string, taskID string) error {
	return c.inspector.ArchiveTask(queue, taskID)
}

// DeleteTask deletes a task
func (c *queue_client) DeleteTask(queue string, taskID string) error {
	return c.inspector.DeleteTask(queue, taskID)
}

// RunTask runs an archived or retry task immediately
func (c *queue_client) RunTask(queue string, taskID string) error {
	return c.inspector.RunTask(queue, taskID)
}

// GetTaskInfo returns info about a specific task
func (c *queue_client) GetTaskInfo(queue string, taskID string) (*TaskInfo, error) {
	t, err := c.inspector.GetTaskInfo(queue, taskID)
	if err != nil {
		return nil, err
	}
	return convertTaskInfo(t), nil
}
