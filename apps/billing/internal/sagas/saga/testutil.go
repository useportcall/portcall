package saga

import "fmt"

// TraceResult holds the result of tracing a saga's step connectivity.
type TraceResult struct {
	// Reachable is the ordered list of route names reachable from the entry step.
	Reachable []string
	// Unreachable lists routes defined but not reachable from the entry step.
	Unreachable []string
}

// Trace walks the Emits graph starting from the first step and returns
// which routes are reachable. This is useful for verifying that all steps
// in a saga are connected from the entry point.
func Trace(steps []Step) (*TraceResult, error) {
	if len(steps) == 0 {
		return nil, fmt.Errorf("no steps to trace")
	}

	stepMap := make(map[string]*Step, len(steps))
	for i := range steps {
		stepMap[steps[i].Route.Name] = &steps[i]
	}

	visited := make(map[string]bool)
	var reachable []string

	var walk func(name string)
	walk = func(name string) {
		if visited[name] {
			return
		}
		visited[name] = true
		reachable = append(reachable, name)
		s, ok := stepMap[name]
		if !ok {
			return
		}
		for _, e := range s.Emits {
			walk(e.Name)
		}
	}

	walk(steps[0].Route.Name)

	var unreachable []string
	for _, s := range steps {
		if !visited[s.Route.Name] {
			unreachable = append(unreachable, s.Route.Name)
		}
	}

	return &TraceResult{
		Reachable:   reachable,
		Unreachable: unreachable,
	}, nil
}

// ExpectFullyConnected verifies that every step in the saga is reachable
// from the entry step via the Emits graph.
func ExpectFullyConnected(steps []Step) error {
	result, err := Trace(steps)
	if err != nil {
		return err
	}
	if len(result.Unreachable) > 0 {
		return fmt.Errorf("unreachable steps: %v", result.Unreachable)
	}
	return nil
}
