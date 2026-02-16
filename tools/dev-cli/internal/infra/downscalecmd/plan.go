package downscalecmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func BuildPlan(c *cobra.Command, opts Options, deps Deps) (Plan, error) {
	if err := deps.EnsureRootDir(); err != nil {
		return Plan{}, err
	}
	token, err := deps.ResolveDOToken()
	if err != nil {
		return Plan{}, err
	}
	plan := Plan{env: map[string]string{"DIGITALOCEAN_TOKEN": token}, targetNodeCnt: 1, targetRedisSz: "db-s-1vcpu-1gb"}
	plan.cluster = strings.TrimSpace(opts.ClusterName)
	if plan.cluster == "" || !c.Flags().Changed("name") {
		plan.cluster = firstNonEmpty(clusterFromState(opts.Cluster, deps), "portcall-prod")
	}
	plan.nodePool, err = deps.FetchNodePool(plan.env, plan.cluster)
	if err != nil {
		return Plan{}, err
	}
	plan.targetNodeSz = firstNonEmpty(strings.TrimSpace(opts.NodeSize), plan.nodePool.Size)
	if c.Flags().Changed("node-count") {
		plan.targetNodeCnt = opts.NodeCount
	}
	if plan.targetNodeCnt < 1 {
		return Plan{}, fmt.Errorf("--node-count must be >= 1")
	}
	plan.redis, _ = deps.DiscoverRedis(plan.env)
	plan.postgres, _ = deps.DiscoverPostgres(plan.env)
	if strings.TrimSpace(opts.RedisSize) != "" {
		plan.targetRedisSz = strings.TrimSpace(opts.RedisSize)
	} else if plan.redis.Size != "" {
		plan.targetRedisSz = plan.redis.Size
	}
	plan.targetPgSz = strings.TrimSpace(opts.PostgresSize)
	if !opts.Yes && deps.IsInteractive() {
		return runWizard(plan, deps)
	}
	if plan.targetPgSz != "" && deps.IsDowngradeSize(plan.postgres.Size, plan.targetPgSz) {
		return Plan{}, fmt.Errorf("postgres downsize (%s -> %s) requires migration flow; in-place resize is blocked\n%s", plan.postgres.Size, plan.targetPgSz, postgresMigrationHint(plan.postgres.Name, plan.targetPgSz))
	}
	return plan, nil
}

func runWizard(plan Plan, deps Deps) (Plan, error) {
	deps.Section("Infra Downscale Wizard")
	deps.Warn("Possible brief downtime during node and database resizing")
	plan.cluster = deps.AskText("Cluster name ["+plan.cluster+"]: ", plan.cluster)
	plan.targetNodeSz = deps.AskText("Node size slug ["+plan.nodePool.Size+"]: ", plan.targetNodeSz)
	nodeCountText := deps.AskText("Node count [1]: ", strconv.Itoa(plan.targetNodeCnt))
	nodeCount, err := strconv.Atoi(strings.TrimSpace(nodeCountText))
	if err != nil || nodeCount < 1 {
		return Plan{}, fmt.Errorf("invalid node count %q (must be >= 1)", nodeCountText)
	}
	plan.targetNodeCnt = nodeCount
	if plan.redis.ID != "" {
		plan.targetRedisSz = deps.AskText("Redis size ["+plan.redis.Size+"]: ", firstNonEmpty(plan.targetRedisSz, plan.redis.Size))
	}
	if plan.postgres.ID != "" && deps.AskYesNo("Include Postgres resize step? [y/N]: ", false) {
		plan.targetPgSz = deps.AskText("Postgres size ["+plan.postgres.Size+"]: ", plan.postgres.Size)
	}
	if plan.targetPgSz != "" && deps.IsDowngradeSize(plan.postgres.Size, plan.targetPgSz) {
		deps.Warn("Postgres in-place downsize is unsafe and blocked to protect data")
		deps.Plain(postgresMigrationHint(plan.postgres.Name, plan.targetPgSz))
		if deps.AskYesNo("Skip Postgres downsize and continue with other steps? [Y/n]: ", true) {
			plan.targetPgSz = ""
		} else {
			return Plan{}, fmt.Errorf("downscale canceled")
		}
	}
	return plan, nil
}

func postgresMigrationHint(currentName, targetSize string) string {
	name := firstNonEmpty(strings.TrimSpace(currentName), "portcall-db")
	return "Safe migration steps:\n1) fork/create new Postgres cluster at target size\n2) restore snapshot/logical dump into new cluster\n3) verify app queries + migrations in staging mode\n4) rotate Kubernetes DB secrets/connection URL\n5) deploy + smoke-check, then retire old cluster\nSuggested target: " + name + "-downscale (" + targetSize + ")"
}

func clusterFromState(alias string, deps Deps) string {
	if cfg, ok := deps.GetInfraClusterState(strings.TrimSpace(alias)); ok {
		return strings.TrimSpace(cfg.Cluster)
	}
	return ""
}

func firstNonEmpty(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
