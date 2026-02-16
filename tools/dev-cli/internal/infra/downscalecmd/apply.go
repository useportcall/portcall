package downscalecmd

import (
	"fmt"
	"strings"
	"time"
)

func PrintSummary(plan Plan, deps Deps) {
	deps.Section("Downscale Summary")
	deps.Plain("Cluster       : %s", plan.cluster)
	deps.Plain("Node pool     : %s (%s x%d -> %s x%d)", plan.nodePool.Name, plan.nodePool.Size, plan.nodePool.Count, plan.targetNodeSz, plan.targetNodeCnt)
	if plan.redis.ID != "" {
		deps.Plain("Redis cluster : %s (%s -> %s)", plan.redis.Name, plan.redis.Size, plan.targetRedisSz)
	}
	if plan.postgres.ID != "" {
		deps.Plain("Postgres      : %s (%s -> %s)", plan.postgres.Name, plan.postgres.Size, firstNonEmpty(strings.TrimSpace(plan.targetPgSz), "(no change)"))
	}
	if before, after, err := deps.NodeMonthlyCost(plan.env, plan.nodePool.Size, plan.nodePool.Count, plan.targetNodeSz, plan.targetNodeCnt); err == nil {
		deps.Plain("Node cost/mo  : $%.2f -> $%.2f (save $%.2f)", before, after, before-after)
	}
	if plan.targetPgSz != "" && deps.IsDowngradeSize(plan.postgres.Size, plan.targetPgSz) {
		deps.Warn("Postgres size downscale requires migration workflow to keep data safe")
	}
}

func RunApply(plan Plan, deps Deps) error {
	deps.Section("Step 1/4: Node Pool")
	if err := deps.ApplyNodePool(plan.env, plan.cluster, plan.nodePool, plan.targetNodeSz, plan.targetNodeCnt); err != nil {
		return err
	}
	deps.OK("Node pool step complete")
	deps.Section("Step 2/4: Redis")
	if err := deps.ApplyRedis(plan.env, plan.redis.ID, plan.redis.Size, plan.targetRedisSz); err != nil {
		return err
	}
	deps.OK("Redis step complete")
	deps.Section("Step 3/4: Postgres")
	if strings.TrimSpace(plan.targetPgSz) == "" {
		deps.Info("Skipping Postgres resize")
		return nil
	}
	if err := deps.ApplyPostgres(plan.env, plan.postgres, plan.targetPgSz); err != nil {
		return fmt.Errorf("postgres step failed: %w", err)
	}
	deps.OK("Postgres step complete")
	return nil
}

func RunSmokeCheck(alias string, timeoutSec int, deps Deps) error {
	if timeoutSec < 30 {
		timeoutSec = 30
	}
	deps.Section("Smoke Check")
	deps.Info("Switching kubectl context to %s", deps.ResolveClusterContext(alias))
	if err := deps.SwitchCluster(alias); err != nil {
		return err
	}
	deadline := time.Now().Add(time.Duration(timeoutSec) * time.Second)
	for {
		pending, err := unhealthyPods(deps.RunCmdOut, "portcall")
		if err != nil {
			return err
		}
		if len(pending) == 0 {
			deps.OK("All non-job pods are Running")
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("smoke check timed out with unhealthy pods: %v", pending)
		}
		deps.Info("Waiting for pods to stabilize (%d pending)...", len(pending))
		time.Sleep(10 * time.Second)
	}
}

func unhealthyPods(runCmdOut func(name string, args ...string) (string, error), namespace string) ([]string, error) {
	out, err := runCmdOut("kubectl", "-n", namespace, "get", "pods", "--no-headers")
	if err != nil {
		return nil, err
	}
	items := []string{}
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		status := fields[2]
		if status == "Running" || status == "Completed" || status == "Succeeded" {
			continue
		}
		items = append(items, fields[0]+":"+status)
	}
	return items, nil
}
