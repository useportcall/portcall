package downscalecmd

import (
	infraDownscale "github.com/useportcall/portcall/tools/dev-cli/internal/infra/downscale"
	infraState "github.com/useportcall/portcall/tools/dev-cli/internal/infra/state"
)

type Options struct {
	Cluster        string
	ClusterName    string
	NodeSize       string
	NodeCount      int
	RedisSize      string
	PostgresSize   string
	SkipSmokeCheck bool
	SmokeTimeout   int
	Yes            bool
	DryRun         bool
}

type Plan struct {
	env           map[string]string
	cluster       string
	targetNodeSz  string
	targetNodeCnt int
	targetRedisSz string
	targetPgSz    string
	nodePool      infraDownscale.NodePool
	redis         infraDownscale.DBRow
	postgres      infraDownscale.DBRow
}

type Deps struct {
	EnsureRootDir         func() error
	ResolveDOToken        func() (string, error)
	FetchNodePool         func(env map[string]string, cluster string) (infraDownscale.NodePool, error)
	DiscoverRedis         func(env map[string]string) (infraDownscale.DBRow, error)
	DiscoverPostgres      func(env map[string]string) (infraDownscale.DBRow, error)
	IsDowngradeSize       func(current string, target string) bool
	ApplyNodePool         func(env map[string]string, cluster string, current infraDownscale.NodePool, targetSize string, targetCount int) error
	ApplyRedis            func(env map[string]string, redisID string, currentSize string, targetSize string) error
	ApplyPostgres         func(env map[string]string, pg infraDownscale.DBRow, targetSize string) error
	NodeMonthlyCost       func(env map[string]string, beforeSize string, beforeCount int, afterSize string, afterCount int) (float64, float64, error)
	GetInfraClusterState  func(cluster string) (infraState.ClusterState, bool)
	AskText               func(prompt string, fallback string) string
	AskYesNo              func(prompt string, defaultYes bool) bool
	IsInteractive         func() bool
	Section               func(title string)
	Plain                 func(msg string, args ...any)
	Info                  func(msg string, args ...any)
	Warn                  func(msg string, args ...any)
	OK                    func(msg string, args ...any)
	SwitchCluster         func(cluster string) error
	ResolveClusterContext func(cluster string) string
	RunCmdOut             func(name string, args ...string) (string, error)
}
