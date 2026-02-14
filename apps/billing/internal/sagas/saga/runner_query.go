package saga

import (
	"encoding/json"
	"fmt"
)

// TaskPayload finds the first enqueued task with the given name and
// unmarshals its payload into dest.
func (r *Runner) TaskPayload(name string, dest any) error {
	for _, t := range r.AllTasks {
		if t.Name == name {
			return json.Unmarshal(t.Payload, dest)
		}
	}
	return fmt.Errorf("task %q not found in executed tasks", name)
}

// TaskPayloads finds all enqueued tasks with the given name and returns
// their raw payloads.
func (r *Runner) TaskPayloads(name string) []json.RawMessage {
	var payloads []json.RawMessage
	for _, t := range r.AllTasks {
		if t.Name == name {
			payloads = append(payloads, t.Payload)
		}
	}
	return payloads
}

// HasTask returns true if a task with the given name was executed.
func (r *Runner) HasTask(name string) bool {
	for _, n := range r.Executed {
		if n == name {
			return true
		}
	}
	return false
}
