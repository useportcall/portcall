package saga

import "fmt"

// ValidateSteps checks a single saga's steps for internal consistency:
// no empty routes, no nil handlers, no duplicate route names.
func ValidateSteps(steps []Step) error {
	seen := make(map[string]bool)
	for _, s := range steps {
		if s.Route.Name == "" {
			return fmt.Errorf("step has empty route name")
		}
		if s.Handler == nil {
			return fmt.Errorf("step %q has nil handler", s.Route.Name)
		}
		if seen[s.Route.Name] {
			return fmt.Errorf("duplicate route %q", s.Route.Name)
		}
		seen[s.Route.Name] = true
	}
	return nil
}

// ValidateAll checks cross-saga consistency:
// - no route registered by multiple sagas
// - every billing-queue route in Emits is registered by some saga
func ValidateAll(sagas ...[]Step) error {
	defined := make(map[string]bool)
	for _, steps := range sagas {
		if err := ValidateSteps(steps); err != nil {
			return err
		}
		for _, s := range steps {
			if defined[s.Route.Name] {
				return fmt.Errorf("route %q registered by multiple sagas", s.Route.Name)
			}
			defined[s.Route.Name] = true
		}
	}
	for _, steps := range sagas {
		for _, s := range steps {
			for _, e := range s.Emits {
				if e.Queue != BillingQueue {
					continue // external queues (email, etc.) are not validated
				}
				if !defined[e.Name] {
					return fmt.Errorf("step %q emits unregistered route %q", s.Route.Name, e.Name)
				}
			}
		}
	}
	return nil
}

// ExpectRoutes validates steps and asserts that routes appear in the given order.
// Useful in per-module tests to lock down the expected saga flow.
func ExpectRoutes(steps []Step, names ...string) error {
	if err := ValidateSteps(steps); err != nil {
		return err
	}
	if len(steps) != len(names) {
		return fmt.Errorf("expected %d steps, got %d", len(names), len(steps))
	}
	for i, name := range names {
		if steps[i].Route.Name != name {
			return fmt.Errorf("step %d: expected route %q, got %q", i, name, steps[i].Route.Name)
		}
	}
	return nil
}
