package harness

import "github.com/useportcall/portcall/libs/go/qx"

// NopQueue is a no-op queue implementation for tests. All operations
// succeed silently without connecting to Redis.
type NopQueue struct{}

var _ qx.IQueue = (*NopQueue)(nil)

func (q *NopQueue) Enqueue(string, any, string) error               { return nil }
func (q *NopQueue) Close() error                                     { return nil }
func (q *NopQueue) ListQueues() ([]string, error)                    { return nil, nil }
func (q *NopQueue) GetQueueStats(string) (*qx.QueueStats, error)    { return nil, nil }
func (q *NopQueue) ArchiveTask(string, string) error                 { return nil }
func (q *NopQueue) DeleteTask(string, string) error                  { return nil }
func (q *NopQueue) RunTask(string, string) error                     { return nil }
func (q *NopQueue) GetTaskInfo(string, string) (*qx.TaskInfo, error) { return nil, nil }

func (q *NopQueue) ListTasks(string, string, int, string) ([]*qx.TaskInfo, error) {
	return nil, nil
}
