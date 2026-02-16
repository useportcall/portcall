package initcmd

import (
	"github.com/useportcall/portcall/tools/dev-cli/internal/infra/state"
	"github.com/useportcall/portcall/tools/dev-cli/internal/infra/tfflow"
)

type ClusterState = state.ClusterState
type Step = tfflow.Step

type Options struct {
	Cluster           string
	Provider          string
	Providers         []string
	Action            string
	Mode              string
	Step              string
	ClusterName       string
	Region            string
	NodeSize          string
	NodeCount         int
	RedisSize         string
	PostgresSize      string
	RegistryName      string
	VPCCIDR           string
	SpacesRegion      string
	SpacesPrefix      string
	Domain            string
	DNSProvider       string
	DNSAuto           bool
	CloudflareZoneID  string
	AllowedIPs        []string
	AllowedIPsFrom    string
	InstallIngress    bool
	SkipRegistryLogin bool
	SkipPostgres      bool
	SkipRedis         bool
	SkipSpaces        bool
	SkipSmokeCheck    bool
	SmokeTimeout      int
	Yes               bool
	DryRun            bool
}

type Plan struct {
	Alias        string
	Providers    []string
	Provider     string
	Action       string
	Step         string
	ClusterName  string
	Region       string
	NodeSize     string
	NodeCount    int
	RedisSize    string
	RegistryName string
	Registry     string
	VPCCIDR      string
	SpacesRegion string
	SpacesPrefix string
	Domain       string
	AllowedIPs   []string
	Mode         string
	Namespace    string
	StackDir     string
}

type Deps struct {
	EnsureRootDir     func() error
	RootDir           func() string
	RunCmd            func(string, ...string) error
	RunCmdWithEnv     func(map[string]string, string, ...string) error
	RunCmdOut         func(string, ...string) (string, error)
	RunCmdOutWithEnv  func(map[string]string, string, ...string) (string, error)
	RunShell          func(string) error
	ReadInput         func(string) string
	ResolveDOToken    func() (string, error)
	VerifyDOAccess    func(string) error
	ResolveProviders  func([]string, string) ([]string, error)
	ResolveAllowedIPs func([]string, string) ([]string, error)
	DecodeB64         func(string) string
	RandomSecret      func(int) (string, error)
	GetClusterState   func(string) (ClusterState, bool)
	SetClusterState   func(string, ClusterState) error
	RunTerraformSteps func(Plan, Options, []Step) error
	ResolveInitSteps  func(Plan, Options) ([]Step, error)
	Plain             func(string, ...any)
	Info              func(string, ...any)
	OK                func(string, ...any)
	Warn              func(string, ...any)
	Fail              func(string, ...any)
	Section           func(string)
	AskYesNo          func(string, bool) bool
	AskText           func(string, string) string
	IsInteractive     func() bool
}
