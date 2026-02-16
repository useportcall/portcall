package deploy

import (
	"fmt"
	"time"
)

func waitForObservability() error {
	info("Watching rollout for observability...")
	checks := [][]string{
		{"rollout", "status", "statefulset/loki", "-n", "portcall", "--request-timeout=15s", "--timeout=300s"},
		{"rollout", "status", "deployment/grafana", "-n", "portcall", "--request-timeout=15s", "--timeout=300s"},
		{"rollout", "status", "daemonset/promtail", "-n", "portcall", "--request-timeout=15s", "--timeout=300s"},
	}
	for _, c := range checks {
		if err := runCmdWithTimeout(6*time.Minute, "kubectl", c...); err != nil {
			if pods, podsErr := runCmdOut("kubectl", "get", "pods", "-n", "portcall", "-l", "app=grafana", "-o", "wide"); podsErr == nil && pods != "" {
				warn("Grafana pod status:\n%s", pods)
			}
			return fmt.Errorf("observability rollout check failed (%v): %w", c, err)
		}
	}
	ok("Successfully deployed observability")
	return nil
}
