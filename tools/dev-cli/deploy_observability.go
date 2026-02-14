package main

import "fmt"

func waitForObservability() error {
	info("Watching rollout for observability...")
	checks := [][]string{
		{"rollout", "status", "statefulset/loki", "-n", "portcall", "--timeout=120s"},
		{"rollout", "status", "deployment/grafana", "-n", "portcall", "--timeout=120s"},
		{"rollout", "status", "daemonset/promtail", "-n", "portcall", "--timeout=120s"},
	}
	for _, c := range checks {
		if err := runCmd("kubectl", c...); err != nil {
			return fmt.Errorf("observability rollout check failed (%v): %w", c, err)
		}
	}
	ok("Successfully deployed observability")
	return nil
}
