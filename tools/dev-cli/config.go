package main

// AppConfig defines an application and its configuration.
type AppConfig struct {
	Name          string
	ContainerName string
	ComposeFile   string // relative to docker-compose/
	BackendCmd    string // command to run backend locally
	FrontendCmd   string // command to run frontend locally (empty if none)
	IsWorker      bool   // true if this is a worker (runs in workers.yml)
}

var apps = []AppConfig{
	{Name: "api", ContainerName: "api", ComposeFile: "docker-compose.api.yml", BackendCmd: "apps/api"},
	{Name: "admin", ContainerName: "admin", ComposeFile: "docker-compose.admin.yml", BackendCmd: "apps/admin", FrontendCmd: "apps/admin/frontend"},
	{Name: "file", ContainerName: "file-api", ComposeFile: "docker-compose.file.yml", BackendCmd: "apps/file"},
	{Name: "quote", ContainerName: "quote", ComposeFile: "docker-compose.quote.yml", BackendCmd: "apps/quote"},
	{Name: "dashboard", ContainerName: "dashboard", ComposeFile: "docker-compose.dashboard.yml", BackendCmd: "apps/dashboard", FrontendCmd: "apps/dashboard/frontend"},
	{Name: "checkout", ContainerName: "checkout", ComposeFile: "docker-compose.checkout.yml", BackendCmd: "apps/checkout", FrontendCmd: "apps/checkout/frontend"},
	{Name: "billing", ContainerName: "billing_worker", BackendCmd: "apps/billing", IsWorker: true},
	{Name: "email", ContainerName: "email_worker", BackendCmd: "apps/email", IsWorker: true},
}

// Infrastructure compose files (always run).
var infraComposeFiles = []string{
	"docker-compose.db.yml",
	"docker-compose.auth.yml",
	"docker-compose.tools.yml",
	"docker-compose.workers.yml",
}

// Presets define common development configurations.
var presets = map[string]struct {
	Description string
	Docker      []string
	Terminal    []string
}{
	"dashboard": {
		Description: "Dashboard development (dashboard+checkout in terminal, others in Docker)",
		Docker:      []string{"api", "admin", "file", "quote"},
		Terminal:    []string{"dashboard", "checkout"},
	},
	"billing": {
		Description: "Billing worker development (billing in terminal, minimal Docker)",
		Docker:      []string{"api"},
		Terminal:    []string{"billing"},
	},
	"all-docker": {
		Description: "All apps in Docker (no terminal windows)",
		Docker:      []string{"api", "admin", "file", "quote", "dashboard", "checkout"},
		Terminal:    []string{},
	},
	"minimal": {
		Description: "Minimal setup - just infrastructure, no apps",
		Docker:      []string{},
		Terminal:    []string{},
	},
	"quick": {
		Description: "Quick dashboard dev (dashboard only, minimal Docker)",
		Docker:      []string{"api"},
		Terminal:    []string{"dashboard"},
	},
}

// E2E test targets (name → list of test package paths relative to repo root).
var e2eTargets = map[string][]string{
	"all": {"./e2etest/..."},
}
