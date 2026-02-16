package downscale

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func FetchNodePool(run RunCmdOutWithEnv, env map[string]string, cluster string) (NodePool, error) {
	out, err := run(env, "doctl", "kubernetes", "cluster", "node-pool", "list", cluster, "--output", "json")
	if err != nil {
		return NodePool{}, fmt.Errorf("failed to read node pools for %s: %w", cluster, err)
	}
	var pools []NodePool
	if err := json.Unmarshal([]byte(out), &pools); err != nil || len(pools) == 0 {
		return NodePool{}, fmt.Errorf("failed to parse node pools for %s", cluster)
	}
	return pools[0], nil
}

func DiscoverRedisCluster(run RunCmdOutWithEnv, env map[string]string) (DBRow, error) {
	return DiscoverDatabaseCluster(run, env, "redis", "valkey")
}

func DiscoverPostgresCluster(run RunCmdOutWithEnv, env map[string]string) (DBRow, error) {
	return DiscoverDatabaseCluster(run, env, "pg", "postgres")
}

func DiscoverDatabaseCluster(run RunCmdOutWithEnv, env map[string]string, engines ...string) (DBRow, error) {
	out, err := run(env, "doctl", "databases", "list", "--format", "ID,Name,Engine,Status,Size,NumNodes", "--no-header")
	if err != nil {
		return DBRow{}, err
	}
	rows := ParseDBRows(out)
	if len(rows) == 1 {
		return rows[0], nil
	}
	for _, row := range rows {
		for _, engine := range engines {
			if strings.EqualFold(row.Engine, engine) {
				return row, nil
			}
		}
	}
	return DBRow{}, nil
}

func ParseDBRows(out string) []DBRow {
	rows := []DBRow{}
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}
		nodes, _ := strconv.Atoi(fields[5])
		rows = append(rows, DBRow{ID: fields[0], Name: fields[1], Engine: fields[2], Status: fields[3], Size: fields[4], Nodes: nodes})
	}
	return rows
}
