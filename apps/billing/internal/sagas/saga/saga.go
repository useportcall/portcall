package saga

import "github.com/useportcall/portcall/libs/go/qx/server"

const (
	BillingQueue = "billing_queue"
	EmailQueue   = "email_queue"
)

// Route is an immutable, typed reference to a queue task.
type Route struct {
	Name  string
	Queue string
}

// Enqueue sends a payload to this route's queue.
func (r Route) Enqueue(q Enqueuer, payload any) error {
	return q.Enqueue(r.Name, payload, r.Queue)
}

// Enqueuer abstracts the queue client (satisfied by server.IContext.Queue()).
type Enqueuer interface {
	Enqueue(taskType string, payload any, queue string) error
}

// Step binds a route to its handler and declares which routes it may emit.
// The Emits field serves as both documentation and a testable contract.
type Step struct {
	Route   Route
	Handler server.HandlerFunc
	Emits   []Route
}

// Register registers all steps in a saga with the queue server.
func Register(srv server.IServer, steps []Step) {
	for _, s := range steps {
		srv.H(s.Route.Name, s.Handler)
	}
}
