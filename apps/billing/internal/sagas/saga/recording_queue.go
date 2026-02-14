package saga

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/useportcall/portcall/libs/go/qx"
)

// EnqueuedTask captures a task that was enqueued during a handler run.
type EnqueuedTask struct {
	Name    string
	Payload json.RawMessage
	Queue   string
}

// RecordingQueue implements qx.IQueue by recording all enqueued tasks.
type RecordingQueue struct {
	mu    sync.Mutex
	Tasks []EnqueuedTask
}

func (q *RecordingQueue) Enqueue(name string, payload any, queue string) error {
	raw, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal enqueue payload: %w", err)
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	q.Tasks = append(q.Tasks, EnqueuedTask{Name: name, Payload: raw, Queue: queue})
	return nil
}

func (q *RecordingQueue) Close() error                                 { return nil }
func (q *RecordingQueue) ListQueues() ([]string, error)                { return nil, nil }
func (q *RecordingQueue) GetQueueStats(string) (*qx.QueueStats, error) { return nil, nil }
func (q *RecordingQueue) ListTasks(string, string, int, string) ([]*qx.TaskInfo, error) {
	return nil, nil
}
func (q *RecordingQueue) ArchiveTask(string, string) error { return nil }
func (q *RecordingQueue) DeleteTask(string, string) error  { return nil }
func (q *RecordingQueue) RunTask(string, string) error     { return nil }
func (q *RecordingQueue) GetTaskInfo(string, string) (*qx.TaskInfo, error) {
	return nil, nil
}

// Pop removes and returns the first recorded task, or nil if empty.
func (q *RecordingQueue) Pop() *EnqueuedTask {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.Tasks) == 0 {
		return nil
	}
	t := q.Tasks[0]
	q.Tasks = q.Tasks[1:]
	return &t
}

// Len returns the number of pending tasks.
func (q *RecordingQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.Tasks)
}

// Reset clears all recorded tasks.
func (q *RecordingQueue) Reset() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.Tasks = nil
}
