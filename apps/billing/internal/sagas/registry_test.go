package sagas_test

import (
	"fmt"
	"testing"

	"github.com/useportcall/portcall/apps/billing/internal/sagas"
	"github.com/useportcall/portcall/apps/billing/internal/sagas/saga"
)

func TestAllSagasValid(t *testing.T) {
	all := sagas.All()
	for i, steps := range all {
		if err := saga.ValidateSteps(steps); err != nil {
			t.Fatalf("saga %d invalid: %v", i, err)
		}
	}
}

func TestNoDuplicateRoutes(t *testing.T) {
	if err := saga.ValidateAll(sagas.All()...); err != nil {
		t.Fatal(err)
	}
}

func TestTotalRouteCount(t *testing.T) {
	total := 0
	for _, steps := range sagas.All() {
		total += len(steps)
	}
	// 25 billing routes across 11 sagas — update this if you add/remove routes.
	if total != 25 {
		t.Fatalf("expected 25 total routes, got %d", total)
	}
}

func TestAllSagasFullyConnected(t *testing.T) {
	for i, steps := range sagas.All() {
		if i == 0 {
			if err := expectReachableFrom(steps,
				"process_stripe_webhook_event",
				"process_braintree_webhook_event",
			); err != nil {
				t.Fatalf("saga %d not fully connected: %v", i, err)
			}
			continue
		}
		if err := saga.ExpectFullyConnected(steps); err != nil {
			t.Fatalf("saga %d not fully connected: %v", i, err)
		}
	}
}

func expectReachableFrom(steps []saga.Step, roots ...string) error {
	graph := make(map[string][]string, len(steps))
	for _, step := range steps {
		for _, emit := range step.Emits {
			graph[step.Route.Name] = append(graph[step.Route.Name], emit.Name)
		}
		if _, ok := graph[step.Route.Name]; !ok {
			graph[step.Route.Name] = nil
		}
	}
	visited := map[string]bool{}
	var walk func(string)
	walk = func(route string) {
		if visited[route] {
			return
		}
		visited[route] = true
		for _, next := range graph[route] {
			walk(next)
		}
	}
	for _, root := range roots {
		walk(root)
	}
	for _, step := range steps {
		if !visited[step.Route.Name] {
			return fmt.Errorf("unreachable steps: [%s]", step.Route.Name)
		}
	}
	return nil
}
