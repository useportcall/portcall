package queues

import (
	"encoding/json"
	"strconv"

	"github.com/useportcall/portcall/libs/go/qx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type ListTasksResponse struct {
	Tasks []*qx.TaskInfo `json:"tasks"`
	Queue string         `json:"queue"`
	State string         `json:"state"`
	Limit int            `json:"limit"`
}

type QueueStatsResponse struct {
	Stats []*qx.QueueStats `json:"stats"`
}

// GetQueueStats returns statistics for queues
func GetQueueStats(c *routerx.Context) {
	queue := c.Query("queue")

	if queue != "" {
		stats, err := c.Queue().GetQueueStats(queue)
		if err != nil {
			c.ServerError("Failed to get queue stats", err)
			return
		}
		c.OK(&QueueStatsResponse{Stats: []*qx.QueueStats{stats}})
		return
	}

	// Get stats for known queues
	knownQueues := []string{"billing_queue", "email_queue"}
	allStats := make([]*qx.QueueStats, 0, len(knownQueues))

	for _, q := range knownQueues {
		stats, err := c.Queue().GetQueueStats(q)
		if err != nil {
			// Queue might not exist yet, skip it
			continue
		}
		allStats = append(allStats, stats)
	}

	c.OK(&QueueStatsResponse{Stats: allStats})
}

// ListTasks returns tasks from a queue
func ListTasks(c *routerx.Context) {
	queue := c.Query("queue")
	if queue == "" {
		queue = "billing_queue"
	}

	state := c.Query("state")        // active, pending, scheduled, retry, archived, completed, or empty for all
	taskType := c.Query("task_type") // filter by task type name

	limitStr := c.Query("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	tasks, err := c.Queue().ListTasks(queue, state, limit, taskType)
	if err != nil {
		c.ServerError("Failed to list tasks", err)
		return
	}

	c.OK(&ListTasksResponse{
		Tasks: tasks,
		Queue: queue,
		State: state,
		Limit: limit,
	})
}

// GetTask returns info about a specific task
func GetTask(c *routerx.Context) {
	queue := c.Query("queue")
	if queue == "" {
		c.BadRequest("queue is required")
		return
	}

	taskID := c.Param("task_id")
	if taskID == "" {
		c.BadRequest("task_id is required")
		return
	}

	task, err := c.Queue().GetTaskInfo(queue, taskID)
	if err != nil {
		c.ServerError("Failed to get task info", err)
		return
	}

	c.OK(task)
}

type ArchiveTaskRequest struct {
	Queue  string `json:"queue"`
	TaskID string `json:"task_id"`
}

// ArchiveTask archives a task
func ArchiveTask(c *routerx.Context) {
	var req ArchiveTaskRequest
	if err := c.BindJSON(&req); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	if req.Queue == "" || req.TaskID == "" {
		c.BadRequest("queue and task_id are required")
		return
	}

	err := c.Queue().ArchiveTask(req.Queue, req.TaskID)
	if err != nil {
		c.ServerError("Failed to archive task", err)
		return
	}

	c.OK(map[string]interface{}{
		"success": true,
		"message": "Task archived successfully",
	})
}

type DeleteTaskRequest struct {
	Queue  string `json:"queue"`
	TaskID string `json:"task_id"`
}

// DeleteTask deletes a task
func DeleteTask(c *routerx.Context) {
	var req DeleteTaskRequest
	if err := c.BindJSON(&req); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	if req.Queue == "" || req.TaskID == "" {
		c.BadRequest("queue and task_id are required")
		return
	}

	err := c.Queue().DeleteTask(req.Queue, req.TaskID)
	if err != nil {
		c.ServerError("Failed to delete task", err)
		return
	}

	c.OK(map[string]interface{}{
		"success": true,
		"message": "Task deleted successfully",
	})
}

type RunTaskRequest struct {
	Queue  string `json:"queue"`
	TaskID string `json:"task_id"`
}

// RunTask runs a task immediately
func RunTask(c *routerx.Context) {
	var req RunTaskRequest
	if err := c.BindJSON(&req); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	if req.Queue == "" || req.TaskID == "" {
		c.BadRequest("queue and task_id are required")
		return
	}

	err := c.Queue().RunTask(req.Queue, req.TaskID)
	if err != nil {
		c.ServerError("Failed to run task", err)
		return
	}

	c.OK(map[string]interface{}{
		"success": true,
		"message": "Task queued to run immediately",
	})
}

type RetryWithModifiedPayloadRequest struct {
	Queue    string          `json:"queue"`
	TaskType string          `json:"task_type"`
	Payload  json.RawMessage `json:"payload"`
}

// RetryWithModifiedPayload creates a new task with modified payload (for dead letter debugging)
func RetryWithModifiedPayload(c *routerx.Context) {
	var req RetryWithModifiedPayloadRequest
	if err := c.BindJSON(&req); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	if req.Queue == "" || req.TaskType == "" {
		c.BadRequest("queue and task_type are required")
		return
	}

	// Parse the payload
	var payload map[string]interface{}
	if len(req.Payload) > 0 {
		if err := json.Unmarshal(req.Payload, &payload); err != nil {
			c.BadRequest("Invalid payload JSON")
			return
		}
	}

	// Enqueue the new task
	err := c.Queue().Enqueue(req.TaskType, payload, req.Queue)
	if err != nil {
		c.ServerError("Failed to enqueue task", err)
		return
	}

	c.OK(map[string]interface{}{
		"success":   true,
		"message":   "New task enqueued with modified payload",
		"queue":     req.Queue,
		"task_type": req.TaskType,
	})
}
