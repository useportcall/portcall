package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/useportcall/portcall/libs/go/envx"
	"github.com/useportcall/portcall/libs/go/qx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func main() {
	envx.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	q := qx.New()

	r := routerx.New(nil, nil, q)

	r.GET("/ping", func(c *routerx.Context) { c.JSON(http.StatusOK, gin.H{"message": "pong"}) })
	r.POST("/queues/:queue_id/tasks/:task_id", HandleQueueTask)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func HandleQueueTask(c *routerx.Context) {
	queueID := c.Param("queue_id")
	taskID := c.Param("task_id")

	var payload map[string]any
	if err := c.BindJSON(&payload); err != nil {
		log.Printf("JSON_BIND_ERROR: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := c.Queue().Enqueue(taskID, payload, queueID); err != nil {
		log.Printf("ENQUEUE_TASK_ERROR: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Task enqueued successfully"})
}
