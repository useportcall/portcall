package infra

import (
	"github.com/spf13/cobra"
	cleanup "github.com/useportcall/portcall/tools/dev-cli/internal/infra/cleanupcmd"
	doctor "github.com/useportcall/portcall/tools/dev-cli/internal/infra/doctorcmd"
	downscalecmd "github.com/useportcall/portcall/tools/dev-cli/internal/infra/downscalecmd"
	"github.com/useportcall/portcall/tools/dev-cli/internal/infra/initcmd"
	pull "github.com/useportcall/portcall/tools/dev-cli/internal/infra/pullcmd"
	statepkg "github.com/useportcall/portcall/tools/dev-cli/internal/infra/state"
	status "github.com/useportcall/portcall/tools/dev-cli/internal/infra/statuscmd"
	"github.com/useportcall/portcall/tools/dev-cli/internal/infra/tfflow"
)

func newInfraCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "infra",
		Short: "Provision and bootstrap Kubernetes infrastructure",
	}
	init, create, update := newInfraInitCommands()
	cmd.AddCommand(create)
	cmd.AddCommand(update)
	cmd.AddCommand(init)
	cmd.AddCommand(newInfraDownscaleCmd())
	cmd.AddCommand(newInfraStatusCmd())
	cmd.AddCommand(newInfraDoctorCmd())
	cmd.AddCommand(newInfraPullCmd())
	cmd.AddCommand(newInfraCleanupCmd())
	return cmd
}

func newInfraStatusCmd() *cobra.Command {
	return status.New(status.Deps{
		EnsureRootDir:        ensureRootDir,
		GetInfraClusterState: getInfraClusterState,
		ResolveClusterContext: func(cluster string) string {
			return ResolveClusterContext(rootDir, cluster)
		},
		Warn:  warn,
		Plain: plain,
	})
}

func newInfraDownscaleCmd() *cobra.Command {
	return downscalecmd.New(downscalecmd.Deps{
		EnsureRootDir:    ensureRootDir,
		ResolveDOToken:   resolveDigitalOceanToken,
		FetchNodePool:    fetchNodePool,
		DiscoverRedis:    discoverRedisCluster,
		DiscoverPostgres: discoverPostgresCluster,
		IsDowngradeSize:  isDowngradeSize,
		ApplyNodePool:    applyNodePoolDownscale,
		ApplyRedis:       applyRedisDownscale,
		ApplyPostgres:    applyPostgresResize,
		NodeMonthlyCost:  nodeMonthlyCost,
		GetInfraClusterState: func(cluster string) (statepkg.ClusterState, bool) {
			cfg, ok := getInfraClusterState(cluster)
			return statepkg.ClusterState(cfg), ok
		},
		AskText:       askText,
		AskYesNo:      askYesNo,
		IsInteractive: isInteractiveSession,
		Section:       section,
		Plain:         plain,
		Info:          info,
		Warn:          warn,
		OK:            ok,
		SwitchCluster: switchCluster,
		ResolveClusterContext: func(cluster string) string {
			return ResolveClusterContext(rootDir, cluster)
		},
		RunCmdOut: runCmdOut,
	})
}

func newInfraPullCmd() *cobra.Command {
	return pull.New(pull.Deps{
		EnsureRootDir: ensureRootDir,
		SwitchCluster: switchCluster,
		RunCmdOut:     runCmdOut,
		SetInfraClusterState: func(cluster string, cfg statepkg.ClusterState) error {
			return setInfraClusterState(cluster, ClusterState(cfg))
		},
		GetInfraClusterState: func(cluster string) (statepkg.ClusterState, bool) {
			cfg, ok := getInfraClusterState(cluster)
			return statepkg.ClusterState(cfg), ok
		},
		RootDir: func() string { return rootDir },
		OK:      ok,
		Plain:   plain,
	})
}

func newInfraCleanupCmd() *cobra.Command {
	return cleanup.New(cleanup.Deps{
		EnsureRootDir: ensureRootDir,
		SwitchCluster: switchCluster,
		ReadInput:     readInput,
		RunCmd:        runCmd,
		Warn:          warn,
		Plain:         plain,
		OK:            ok,
	})
}

func newInfraDoctorCmd() *cobra.Command {
	return doctor.New(doctor.Deps{
		EnsureRootDir: ensureRootDir,
		ResolveClusterContext: func(cluster string) string {
			return ResolveClusterContext(rootDir, cluster)
		},
		GetInfraClusterState: func(cluster string) (statepkg.ClusterState, bool) {
			cfg, ok := getInfraClusterState(cluster)
			return statepkg.ClusterState(cfg), ok
		},
		RunCmdOut: runCmdOut,
		OK:        ok,
		Warn:      warn,
		Fail:      fail,
	})
}

func newInfraInitCommands() (init, create, update *cobra.Command) {
	return initcmd.NewCommands(initcmd.Deps{
		EnsureRootDir:     ensureRootDir,
		RootDir:           func() string { return rootDir },
		RunCmd:            runCmd,
		RunCmdWithEnv:     runCmdWithEnv,
		RunCmdOut:         runCmdOut,
		RunCmdOutWithEnv:  runCmdOutWithEnv,
		RunShell:          runShell,
		ReadInput:         readInput,
		ResolveDOToken:    resolveDigitalOceanToken,
		VerifyDOAccess:    verifyDigitalOceanAccess,
		ResolveProviders:  resolveInfraProviders,
		ResolveAllowedIPs: resolveAllowedIPs,
		DecodeB64:         decodeB64,
		RandomSecret:      randomSecret,
		GetClusterState:   getInfraClusterState,
		SetClusterState:   setInfraClusterState,
		RunTerraformSteps: func(plan initcmd.Plan, opts initcmd.Options, steps []initcmd.Step) error {
			fp := toInitFlowPlan(plan)
			fo := toInitFlowOptions(opts)
			fd := initFlowDeps()
			return tfflow.RunTerraformSteps(fp, fo, steps, fd)
		},
		ResolveInitSteps: func(plan initcmd.Plan, opts initcmd.Options) ([]initcmd.Step, error) {
			return tfflow.ResolveInitSteps(toInitFlowPlan(plan), toInitFlowOptions(opts), initFlowDeps())
		},
		Plain:         plain,
		Info:          info,
		OK:            ok,
		Warn:          warn,
		Fail:          fail,
		Section:       section,
		AskYesNo:      askYesNo,
		AskText:       askText,
		IsInteractive: isInteractiveSession,
	})
}

func toInitFlowPlan(p initcmd.Plan) tfflow.Plan {
	return tfflow.Plan{
		StackDir: p.StackDir, Step: p.Step, Action: p.Action, Alias: p.Alias,
		ClusterName: p.ClusterName, Provider: p.Provider, Region: p.Region,
		NodeSize: p.NodeSize, NodeCount: p.NodeCount, RedisSize: p.RedisSize,
		RegistryName: p.RegistryName, VPCCIDR: p.VPCCIDR,
		SpacesRegion: p.SpacesRegion, SpacesPrefix: p.SpacesPrefix,
	}
}

func toInitFlowOptions(o initcmd.Options) tfflow.Options {
	return tfflow.Options{
		SkipPostgres: o.SkipPostgres, SkipRedis: o.SkipRedis,
		SkipSpaces: o.SkipSpaces, DryRun: o.DryRun, Yes: o.Yes,
	}
}

func initFlowDeps() tfflow.Deps {
	return tfflow.Deps{
		RunCmd: runCmd, RunCmdWithEnv: runCmdWithEnv,
		ResolveDOToken: resolveDigitalOceanToken,
		Section:        section, Plain: plain, Info: info, Warn: warn,
		AskYesNo: askYesNo, IsInteractive: isInteractiveSession,
	}
}
