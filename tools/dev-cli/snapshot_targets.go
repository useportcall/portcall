package main

// SnapshotMode categorises the type of snapshot artifact.
type SnapshotMode string

const (
	ModeFullscreen SnapshotMode = "fullscreen"
	ModeComponent  SnapshotMode = "component"
	ModeVideo      SnapshotMode = "video"
)

// SnapshotTarget describes a named, runnable snapshot.
// Add new entries here to make them available via `dev-cli snapshot`.
type SnapshotTarget struct {
	Name        string       // short CLI-friendly name
	Mode        SnapshotMode // category
	Grep        string       // Playwright --grep pattern
	Description string       // shown in --list output
}

// snapshotTargets is the central registry. Append to this slice to add flows.
var snapshotTargets = []SnapshotTarget{
	// ── fullscreen page snapshots ──────────────────────────────
	{Name: "dashboard-home", Mode: ModeFullscreen, Grep: "snapshot home page", Description: "Dashboard home page"},
	{Name: "dashboard-plans", Mode: ModeFullscreen, Grep: "snapshot plans list", Description: "Dashboard plans list"},
	{Name: "dashboard-plan-editor", Mode: ModeFullscreen, Grep: "snapshot plan editor", Description: "Dashboard plan editor"},
	{Name: "dashboard-users", Mode: ModeFullscreen, Grep: "snapshot users list", Description: "Dashboard users list"},
	{Name: "dashboard-user-detail", Mode: ModeFullscreen, Grep: "snapshot user detail", Description: "Dashboard user detail page"},
	{Name: "dashboard-subscriptions", Mode: ModeFullscreen, Grep: "snapshot subscriptions page", Description: "Dashboard subscriptions page"},
	{Name: "dashboard-invoices", Mode: ModeFullscreen, Grep: "snapshot invoices page", Description: "Dashboard invoices page"},
	{Name: "dashboard-developer", Mode: ModeFullscreen, Grep: "snapshot developer page", Description: "Dashboard developer page"},
	{Name: "dashboard-integrations", Mode: ModeFullscreen, Grep: "snapshot integrations page", Description: "Dashboard integrations page"},
	{Name: "dashboard-company", Mode: ModeFullscreen, Grep: "snapshot company page", Description: "Dashboard company settings"},
	{Name: "dashboard-quotes", Mode: ModeFullscreen, Grep: "snapshot quotes list page", Description: "Dashboard quotes list"},
	{Name: "dashboard-usage", Mode: ModeFullscreen, Grep: "snapshot usage page", Description: "Dashboard usage / billing page"},
	{Name: "dashboard-ja", Mode: ModeFullscreen, Grep: "snapshot dashboard home in japanese", Description: "Dashboard home (Japanese locale)"},
	{Name: "checkout-form", Mode: ModeFullscreen, Grep: "snapshot checkout form", Description: "Checkout form page"},
	{Name: "checkout-success", Mode: ModeFullscreen, Grep: "snapshot checkout success", Description: "Checkout success page"},
	{Name: "invoice-light", Mode: ModeFullscreen, Grep: "snapshot invoice light mode", Description: "Invoice view (light mode)"},
	{Name: "invoice-dark", Mode: ModeFullscreen, Grep: "snapshot invoice dark mode", Description: "Invoice view (dark mode)"},
	{Name: "connections-form", Mode: ModeFullscreen, Grep: "snapshot connections form", Description: "Add Provider dialog with Braintree"},
	{Name: "connections-braintree", Mode: ModeFullscreen, Grep: "snapshot connections braintree card", Description: "Connections page showing Braintree card"},

	// ── component snapshots (add entries as tests are created) ─
	// {Name: "plan-card", Mode: ModeComponent, Grep: "snapshot plan card component", Description: "Plan pricing card component"},

	// ── video walkthroughs (add entries as tests are created) ──
	// {Name: "checkout-flow-video", Mode: ModeVideo, Grep: "video checkout full flow", Description: "Video walkthrough of checkout"},
}

// snapshotGroups provides convenient multi-target aliases.
var snapshotGroups = map[string][]string{
	"dashboard-core": {
		"dashboard-home", "dashboard-plans", "dashboard-plan-editor",
		"dashboard-users", "dashboard-user-detail",
	},
	"dashboard-settings": {
		"dashboard-subscriptions", "dashboard-invoices", "dashboard-developer",
		"dashboard-integrations", "dashboard-company", "dashboard-quotes",
		"dashboard-usage",
	},
	"checkout":   {"checkout-form", "checkout-success"},
	"invoice":    {"invoice-light", "invoice-dark"},
	"braintree":  {"connections-form", "connections-braintree"},
}
