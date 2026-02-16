package infra

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/useportcall/portcall/tools/dev-cli/internal/infra/doauth"
	"github.com/useportcall/portcall/tools/dev-cli/internal/infra/downscale"
	"github.com/useportcall/portcall/tools/dev-cli/internal/infra/inputs"
	"github.com/useportcall/portcall/tools/dev-cli/internal/infra/util"
)

type doNodePool = downscale.NodePool
type doDBRow = downscale.DBRow

func resolveInfraProviders(providers []string, fallback string) ([]string, error) {
	return inputs.ResolveProviders(providers, fallback)
}

func resolveAllowedIPs(explicit []string, source string) ([]string, error) {
	if err := ensureRootDir(); err != nil {
		return nil, err
	}
	return inputs.ResolveAllowedIPs(explicit, filepath.Join(rootDir, source))
}

func resolveDigitalOceanToken() (string, error) { return doauth.ResolveDigitalOceanToken(runCmdOut) }
func verifyDigitalOceanAccess(token string) error {
	return doauth.VerifyDigitalOceanAccess(token, runCmdOutWithEnv)
}
func findContextToken(cfg string, context string) string {
	return doauth.FindContextToken(cfg, context)
}

func isDowngradeSize(current, target string) bool { return downscale.IsDowngradeSize(current, target) }
func decodeB64(value string) string               { return util.DecodeB64(value) }
func randomSecret(byteLen int) (string, error)    { return util.RandomSecret(byteLen) }
func parseDBRows(out string) []doDBRow            { return downscale.ParseDBRows(out) }

func fetchNodePool(env map[string]string, cluster string) (doNodePool, error) {
	return downscale.FetchNodePool(runCmdOutWithEnv, env, cluster)
}

func discoverRedisCluster(env map[string]string) (doDBRow, error) {
	return downscale.DiscoverRedisCluster(runCmdOutWithEnv, env)
}

func discoverPostgresCluster(env map[string]string) (doDBRow, error) {
	return downscale.DiscoverPostgresCluster(runCmdOutWithEnv, env)
}

func nodeMonthlyCost(env map[string]string, beforeSize string, beforeCount int, afterSize string, afterCount int) (float64, float64, error) {
	return downscale.NodeMonthlyCost(runCmdOutWithEnv, env, beforeSize, beforeCount, afterSize, afterCount)
}

func applyNodePoolDownscale(env map[string]string, cluster string, current downscale.NodePool, targetSize string, targetCount int) error {
	if strings.TrimSpace(targetSize) == "" {
		targetSize = current.Size
	}
	if targetSize == current.Size {
		return runCmdWithEnv(env, "doctl", "kubernetes", "cluster", "node-pool", "update", cluster, current.ID, "--count", strconv.Itoa(targetCount))
	}
	newName := current.Name + "-downscale"
	if err := runCmdWithEnv(env, "doctl", "kubernetes", "cluster", "node-pool", "create", cluster, "--name", newName, "--size", targetSize, "--count", strconv.Itoa(targetCount)); err != nil {
		return err
	}
	return runCmdWithEnv(env, "doctl", "kubernetes", "cluster", "node-pool", "delete", cluster, current.ID, "--force")
}

func applyRedisDownscale(env map[string]string, redisID, currentSize, targetSize string) error {
	if redisID == "" || targetSize == "" || targetSize == currentSize {
		return nil
	}
	return runCmdWithEnv(env, "doctl", "databases", "resize", redisID, "--num-nodes", "1", "--size", targetSize, "--wait")
}

func applyPostgresResize(env map[string]string, pg downscale.DBRow, targetSize string) error {
	if pg.ID == "" || strings.TrimSpace(targetSize) == "" || pg.Size == targetSize {
		return nil
	}
	if isDowngradeSize(pg.Size, targetSize) {
		return fmt.Errorf("postgres in-place downsize is blocked for data safety. use migration workflow: fork/restore to a smaller cluster, switch secrets, and cut over")
	}
	nodes := pg.Nodes
	if nodes < 1 {
		nodes = 1
	}
	return runCmdWithEnv(env, "doctl", "databases", "resize", pg.ID, "--num-nodes", strconv.Itoa(nodes), "--size", targetSize, "--wait")
}
