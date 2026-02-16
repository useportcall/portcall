package tfflow

type Plan struct {
	StackDir     string
	Step         string
	Action       string
	Alias        string
	ClusterName  string
	Provider     string
	Region       string
	NodeSize     string
	NodeCount    int
	RedisSize    string
	RegistryName string
	VPCCIDR      string
	SpacesRegion string
	SpacesPrefix string
}

type Options struct {
	SkipPostgres bool
	SkipRedis    bool
	SkipSpaces   bool
	DryRun       bool
	Yes          bool
}

type Step struct {
	Name    string
	Targets []string
}

type Deps struct {
	RunCmd         func(name string, args ...string) error
	RunCmdWithEnv  func(env map[string]string, name string, args ...string) error
	ResolveDOToken func() (string, error)
	Section        func(title string)
	Plain          func(msg string, args ...any)
	Info           func(msg string, args ...any)
	Warn           func(msg string, args ...any)
	AskYesNo       func(prompt string, defaultYes bool) bool
	IsInteractive  func() bool
}

var CoreTargets = []string{"digitalocean_vpc.portcall", "digitalocean_kubernetes_cluster.portcall", "digitalocean_container_registry.portcall"}
var PostgresTargets = []string{"digitalocean_database_cluster.postgres", "digitalocean_database_db.main", "digitalocean_database_db.keycloak", "digitalocean_database_firewall.postgres"}
var RedisTargets = []string{"digitalocean_database_cluster.redis", "digitalocean_database_firewall.redis"}
var SpacesTargets = []string{"digitalocean_spaces_bucket.icon_logos", "digitalocean_spaces_bucket.quote_signatures", "digitalocean_spaces_key.portcall"}
