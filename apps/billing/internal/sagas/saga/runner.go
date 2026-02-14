package saga

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

// Runner executes saga handler chains end-to-end, processing enqueued tasks
// through their registered handlers until the queue drains.
type Runner struct {
	handlers map[string]server.HandlerFunc
	queue    *RecordingQueue
	db       dbx.IORM
	crypto   cryptox.ICrypto
	// Executed records every task name processed, in order.
	Executed []string
	// AllTasks records every enqueued task (name + payload) for assertions.
	AllTasks []EnqueuedTask
	// MaxSteps prevents runaway loops; default 100.
	MaxSteps int
}

// NewRunner creates a Runner that can replay handler chains.
// Pass all saga Steps that participate in the chain being tested.
func NewRunner(db dbx.IORM, crypto cryptox.ICrypto, allSteps ...[]Step) *Runner {
	r := &Runner{
		handlers: make(map[string]server.HandlerFunc),
		queue:    &RecordingQueue{},
		db:       db,
		crypto:   crypto,
		MaxSteps: 100,
	}
	for _, steps := range allSteps {
		for _, s := range steps {
			r.handlers[s.Route.Name] = s.Handler
		}
	}
	return r
}

// Run seeds the runner with an initial task and processes all tasks
// until the queue drains or MaxSteps is reached.
func (r *Runner) Run(taskName string, payload any) error {
	raw, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal initial payload: %w", err)
	}

	r.queue.Reset()
	r.Executed = nil
	r.AllTasks = nil

	// Seed the queue.
	r.queue.Tasks = append(r.queue.Tasks, EnqueuedTask{
		Name:    taskName,
		Payload: raw,
		Queue:   BillingQueue,
	})

	steps := 0
	for r.queue.Len() > 0 {
		if steps >= r.MaxSteps {
			return fmt.Errorf("max steps (%d) exceeded, possible infinite loop", r.MaxSteps)
		}
		task := r.queue.Pop()

		handler, ok := r.handlers[task.Name]
		if !ok {
			// Task targets a handler we don't have (e.g. email queue).
			// Record it but don't execute.
			r.Executed = append(r.Executed, task.Name)
			r.AllTasks = append(r.AllTasks, *task)
			steps++
			continue
		}

		ctx := &mockContext{
			ctx:     context.Background(),
			payload: task.Payload,
			db:      r.db,
			queue:   r.queue,
			crypto:  r.crypto,
		}

		r.Executed = append(r.Executed, task.Name)
		r.AllTasks = append(r.AllTasks, *task)

		if err := handler(ctx); err != nil {
			return fmt.Errorf("handler %q failed: %w", task.Name, err)
		}
		steps++
	}
	return nil
}

