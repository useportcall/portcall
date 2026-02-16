// Package main provides the Portcall CLI for generating typed Go code
// from your Portcall app configuration.
//
// Usage:
//
//	go run github.com/useportcall/portcall/sdks/go/cmd/portcall-gen@latest generate
//	go run github.com/useportcall/portcall/sdks/go/cmd/portcall-gen@latest generate --key sk_xxx --output ./portcall
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"
)

const (
	defaultBaseURL     = "https://api.useportcall.com"
	defaultLocalURL    = "http://localhost:9080"
	defaultOutputDir   = "./portcall"
	defaultPackageName = "portcall"
)

// Plan represents a billing plan from the API
type Plan struct {
	ID              string        `json:"id"`
	Name            string        `json:"name"`
	Status          string        `json:"status"`
	Currency        string        `json:"currency"`
	Interval        string        `json:"interval"`
	IntervalCount   int           `json:"interval_count"`
	TrialPeriodDays int           `json:"trial_period_days"`
	IsFree          bool          `json:"is_free"`
	Items           []PlanItem    `json:"items,omitempty"`
	Features        []PlanFeature `json:"features,omitempty"`
}

// PlanItem represents an item within a plan
type PlanItem struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Amount        int64  `json:"amount"`
	Interval      string `json:"interval"`
	IntervalCount int    `json:"interval_count"`
}

// PlanFeature represents a feature attached to a plan
type PlanFeature struct {
	FeatureID string `json:"feature_id"`
	Quota     int64  `json:"quota"`
	Interval  string `json:"interval"`
}

// Feature represents a feature definition from the API
type Feature struct {
	ID        string `json:"id"`
	IsMetered bool   `json:"is_metered"`
}

// AppData holds all fetched app configuration
type AppData struct {
	Plans           []Plan
	Features        []Feature
	MeteredFeatures []Feature
}

// APIResponse wraps the API response format
type APIResponse[T any] struct {
	Data T `json:"data"`
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "generate", "gen":
		if err := runGenerate(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "help", "-h", "--help":
		printUsage()
	case "version", "-v", "--version":
		fmt.Println("portcall-gen v0.2.0")
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Portcall Go Code Generator

Usage:
  portcall-gen <command> [options]

Commands:
  generate, gen    Generate typed Go code from your Portcall app
  help             Show this help message
  version          Show version information

Options for generate:
  --key, -k        Portcall API secret key (or set PC_API_SECRET env var)
  --url, -u        Portcall API base URL (default: prompts for environment)
  --output, -o     Output directory (default: prompts, ./portcall)
  --package, -p    Package name for generated code (default: derived from output dir)
  --yaml           Also generate YAML config file (default: true)
  --no-yaml        Skip generating YAML config file

Examples:
  # Interactive mode (prompts for everything)
  portcall-gen generate

  # With API key
  portcall-gen generate --key sk_xxx

  # Specify local environment and custom output
  portcall-gen generate --url http://localhost:9080 --output ./internal/pc

  # Non-interactive
  portcall-gen generate --key sk_xxx --url http://localhost:9080 --output ./portcall`)
}

func runGenerate() error {
	fmt.Print("\n🚢 Portcall Go Code Generator\n\n")

	// Parse command line arguments
	args := os.Args[2:]
	opts := parseArgs(args)

	// Determine base URL
	baseURL := opts["url"]
	if baseURL == "" {
		env := promptEnvironment()
		if env == "local" {
			port := promptPort()
			baseURL = fmt.Sprintf("http://localhost:%s", port)
		} else {
			baseURL = defaultBaseURL
		}
	}

	// Get API key
	apiKey := opts["key"]
	if apiKey == "" {
		apiKey = os.Getenv("PC_API_SECRET")
	}
	if apiKey == "" {
		apiKey = promptAPIKey()
	}

	// Get output directory (prompt if not provided)
	outputDir := opts["output"]
	if outputDir == "" {
		outputDir = promptOutputDir()
	}

	// Get package name
	packageName := opts["package"]
	if packageName == "" {
		packageName = filepath.Base(outputDir)
		if packageName == "." || packageName == "/" {
			packageName = defaultPackageName
		}
	}

	// Generate YAML unless disabled
	generateYAML := opts["no-yaml"] == ""

	fmt.Println("\n⏳ Fetching app configuration...")

	// Fetch app data from API
	appData, err := fetchAppData(apiKey, baseURL)
	if err != nil {
		return fmt.Errorf("failed to fetch app data: %w", err)
	}

	fmt.Println("✅ Fetched app configuration")
	fmt.Printf("   • %d plan(s)\n", len(appData.Plans))
	fmt.Printf("   • %d feature(s)\n", len(appData.Features))
	fmt.Printf("   • %d metered feature(s)\n", len(appData.MeteredFeatures))

	fmt.Println("\n⏳ Generating Go code...")

	// Generate Go code
	goContent := generateGoCode(appData, packageName, baseURL)

	fmt.Println("✅ Generated Go code")

	// Generate YAML if requested
	var yamlContent string
	if generateYAML {
		fmt.Println("⏳ Generating YAML configuration...")
		yamlContent = generateYAMLConfig(appData)
		fmt.Println("✅ Generated YAML configuration")
	}

	// Write files
	fmt.Println("\n⏳ Writing files...")

	files, err := writeFiles(outputDir, goContent, yamlContent)
	if err != nil {
		return fmt.Errorf("failed to write files: %w", err)
	}

	fmt.Println("✅ Files written successfully")

	// Summary
	fmt.Print("\n✨ Generation complete!\n\n")
	fmt.Println("Generated files:")
	for _, f := range files {
		fmt.Printf("  • %s\n", f)
	}

	fmt.Print("\n📖 Usage example:\n\n")
	fmt.Printf("  import \"yourproject/%s\"\n\n", packageName)
	fmt.Printf("  // Create client (reads PC_API_SECRET from env, uses configured base URL)\n")
	fmt.Printf("  pc := %s.New()\n\n", packageName)
	fmt.Printf("  // Check entitlement for a feature\n")
	fmt.Printf("  ent, err := pc.Check.YourFeature.For(ctx, \"user_id\")\n\n")
	fmt.Printf("  // Create checkout session for a plan\n")
	fmt.Printf("  session, err := pc.Checkout.YourPlan.For(ctx, \"user_id\", successURL, cancelURL)\n\n")
	fmt.Printf("  // Record usage for a metered feature\n")
	fmt.Printf("  err := pc.Record.YourFeature.For(ctx, \"user_id\", 1)\n")

	return nil
}

func parseArgs(args []string) map[string]string {
	opts := make(map[string]string)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--key", "-k":
			if i+1 < len(args) {
				opts["key"] = args[i+1]
				i++
			}
		case "--url", "-u":
			if i+1 < len(args) {
				opts["url"] = args[i+1]
				i++
			}
		case "--output", "-o":
			if i+1 < len(args) {
				opts["output"] = args[i+1]
				i++
			}
		case "--package", "-p":
			if i+1 < len(args) {
				opts["package"] = args[i+1]
				i++
			}
		case "--no-yaml":
			opts["no-yaml"] = "true"
		case "--yaml":
			// Default behavior, just ignore
		}
	}
	return opts
}

func promptEnvironment() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Select your Portcall environment:")
	fmt.Println("  1. Local (localhost)")
	fmt.Println("  2. Production (api.useportcall.com)")
	fmt.Print("\nChoice [1]: ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "2" {
		return "production"
	}
	return "local"
}

func promptPort() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the API port [9080]: ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return "9080"
	}
	return input
}

func promptAPIKey() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your Portcall API secret key: ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if !strings.HasPrefix(input, "sk_") {
		fmt.Println("Warning: API key should start with 'sk_'")
	}

	return input
}

func promptOutputDir() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter output directory [%s]: ", defaultOutputDir)

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultOutputDir
	}
	return input
}

func fetchAppData(apiKey, baseURL string) (*AppData, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	appData := &AppData{}

	// Fetch plans
	plans, err := fetchPlans(client, apiKey, baseURL)
	if err != nil {
		fmt.Printf("Warning: Could not fetch plans: %v\n", err)
	} else {
		appData.Plans = plans
	}

	// Fetch non-metered features
	features, err := fetchFeatures(client, apiKey, baseURL, false)
	if err != nil {
		fmt.Printf("Warning: Could not fetch features: %v\n", err)
	} else {
		appData.Features = features
	}

	// Fetch metered features
	meteredFeatures, err := fetchFeatures(client, apiKey, baseURL, true)
	if err != nil {
		fmt.Printf("Warning: Could not fetch metered features: %v\n", err)
	} else {
		appData.MeteredFeatures = meteredFeatures
	}

	return appData, nil
}

func fetchPlans(client *http.Client, apiKey, baseURL string) ([]Plan, error) {
	req, err := http.NewRequest("GET", baseURL+"/v1/plans", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result APIResponse[[]Plan]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func fetchFeatures(client *http.Client, apiKey, baseURL string, isMetered bool) ([]Feature, error) {
	url := fmt.Sprintf("%s/v1/features?is_metered=%v", baseURL, isMetered)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result APIResponse[[]Feature]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// toGoIdentifier converts a string to a valid Go exported identifier (PascalCase)
func toGoIdentifier(s string) string {
	// Replace non-alphanumeric with underscore
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	s = re.ReplaceAllString(s, "_")

	// Split by underscore and capitalize each part
	parts := strings.Split(s, "_")
	var result strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		// Capitalize first letter
		runes := []rune(part)
		runes[0] = unicode.ToUpper(runes[0])
		result.WriteString(string(runes))
	}

	// If starts with number, prefix with underscore
	out := result.String()
	if len(out) > 0 && unicode.IsDigit(rune(out[0])) {
		out = "_" + out
	}

	return out
}

func generateGoCode(appData *AppData, packageName, baseURL string) string {
	var sb strings.Builder

	allFeatures := append(appData.Features, appData.MeteredFeatures...)

	// Header
	sb.WriteString("// Code generated by portcall-gen. DO NOT EDIT.\n")
	sb.WriteString(fmt.Sprintf("// Generated at: %s\n", time.Now().Format(time.RFC3339)))
	sb.WriteString("//\n")
	sb.WriteString("// Run 'go run github.com/useportcall/portcall/sdks/go/cmd/portcall-gen generate' to regenerate.\n\n")
	sb.WriteString(fmt.Sprintf("package %s\n\n", packageName))

	// Imports
	sb.WriteString("import (\n")
	sb.WriteString("\t\"context\"\n")
	sb.WriteString("\t\"os\"\n")
	sb.WriteString("\n")
	sb.WriteString("\tportcallsdk \"github.com/useportcall/portcall/sdks/go\"\n")
	sb.WriteString(")\n\n")

	// Configuration constants
	sb.WriteString("// =============================================================================\n")
	sb.WriteString("// Configuration\n")
	sb.WriteString("// =============================================================================\n\n")
	sb.WriteString("// BaseURL is the Portcall API URL configured at generation time\n")
	sb.WriteString(fmt.Sprintf("const BaseURL = %q\n\n", baseURL))

	// Plan IDs section
	sb.WriteString("// =============================================================================\n")
	sb.WriteString("// Plan IDs\n")
	sb.WriteString("// =============================================================================\n\n")

	if len(appData.Plans) > 0 {
		sb.WriteString("// Plans contains all plan IDs in your Portcall app\n")
		sb.WriteString("var Plans = struct {\n")
		for _, plan := range appData.Plans {
			identifier := toGoIdentifier(plan.Name)
			sb.WriteString(fmt.Sprintf("\t%s string\n", identifier))
		}
		sb.WriteString("}{\n")
		for _, plan := range appData.Plans {
			identifier := toGoIdentifier(plan.Name)
			sb.WriteString(fmt.Sprintf("\t%s: %q,\n", identifier, plan.ID))
		}
		sb.WriteString("}\n\n")

		sb.WriteString("// Plan ID constants\n")
		sb.WriteString("const (\n")
		for _, plan := range appData.Plans {
			constName := "Plan" + toGoIdentifier(plan.Name)
			sb.WriteString(fmt.Sprintf("\t%s = %q\n", constName, plan.ID))
		}
		sb.WriteString(")\n\n")
	}

	// Feature IDs section
	sb.WriteString("// =============================================================================\n")
	sb.WriteString("// Feature IDs\n")
	sb.WriteString("// =============================================================================\n\n")

	if len(allFeatures) > 0 {
		sb.WriteString("// Features contains all feature IDs in your Portcall app\n")
		sb.WriteString("var Features = struct {\n")
		for _, feature := range allFeatures {
			identifier := toGoIdentifier(feature.ID)
			sb.WriteString(fmt.Sprintf("\t%s string\n", identifier))
		}
		sb.WriteString("}{\n")
		for _, feature := range allFeatures {
			identifier := toGoIdentifier(feature.ID)
			sb.WriteString(fmt.Sprintf("\t%s: %q,\n", identifier, feature.ID))
		}
		sb.WriteString("}\n\n")

		sb.WriteString("// Feature ID constants\n")
		sb.WriteString("const (\n")
		for _, feature := range allFeatures {
			constName := "Feature" + toGoIdentifier(feature.ID)
			sb.WriteString(fmt.Sprintf("\t%s = %q\n", constName, feature.ID))
		}
		sb.WriteString(")\n\n")
	}

	// MeteredFeatures struct
	if len(appData.MeteredFeatures) > 0 {
		sb.WriteString("// MeteredFeatures contains only metered feature IDs\n")
		sb.WriteString("var MeteredFeatures = struct {\n")
		for _, feature := range appData.MeteredFeatures {
			identifier := toGoIdentifier(feature.ID)
			sb.WriteString(fmt.Sprintf("\t%s string\n", identifier))
		}
		sb.WriteString("}{\n")
		for _, feature := range appData.MeteredFeatures {
			identifier := toGoIdentifier(feature.ID)
			sb.WriteString(fmt.Sprintf("\t%s: %q,\n", identifier, feature.ID))
		}
		sb.WriteString("}\n\n")
	}

	// Plan Metadata section
	if len(appData.Plans) > 0 {
		sb.WriteString("// =============================================================================\n")
		sb.WriteString("// Plan Metadata\n")
		sb.WriteString("// =============================================================================\n\n")

		sb.WriteString("// PlanMeta contains metadata about a plan\n")
		sb.WriteString("type PlanMeta struct {\n")
		sb.WriteString("\tID              string\n")
		sb.WriteString("\tName            string\n")
		sb.WriteString("\tStatus          string\n")
		sb.WriteString("\tCurrency        string\n")
		sb.WriteString("\tInterval        string\n")
		sb.WriteString("\tIntervalCount   int\n")
		sb.WriteString("\tTrialPeriodDays int\n")
		sb.WriteString("\tIsFree          bool\n")
		sb.WriteString("}\n\n")

		sb.WriteString("// AllPlans contains metadata for all plans\n")
		sb.WriteString("var AllPlans = []PlanMeta{\n")
		for _, plan := range appData.Plans {
			sb.WriteString(fmt.Sprintf("\t{ID: %q, Name: %q, Status: %q, Currency: %q, Interval: %q, IntervalCount: %d, TrialPeriodDays: %d, IsFree: %v},\n",
				plan.ID, plan.Name, plan.Status, plan.Currency, plan.Interval, plan.IntervalCount, plan.TrialPeriodDays, plan.IsFree))
		}
		sb.WriteString("}\n\n")

		sb.WriteString("// GetPlanMeta returns the metadata for a plan by ID\n")
		sb.WriteString("func GetPlanMeta(planID string) *PlanMeta {\n")
		sb.WriteString("\tfor _, p := range AllPlans {\n")
		sb.WriteString("\t\tif p.ID == planID {\n")
		sb.WriteString("\t\t\treturn &p\n")
		sb.WriteString("\t\t}\n")
		sb.WriteString("\t}\n")
		sb.WriteString("\treturn nil\n")
		sb.WriteString("}\n\n")

		sb.WriteString("// GetFreePlanMeta returns the first free plan metadata, if any\n")
		sb.WriteString("func GetFreePlanMeta() *PlanMeta {\n")
		sb.WriteString("\tfor _, p := range AllPlans {\n")
		sb.WriteString("\t\tif p.IsFree {\n")
		sb.WriteString("\t\t\treturn &p\n")
		sb.WriteString("\t\t}\n")
		sb.WriteString("\t}\n")
		sb.WriteString("\treturn nil\n")
		sb.WriteString("}\n\n")
	}

	// Feature Metadata section
	if len(allFeatures) > 0 {
		sb.WriteString("// =============================================================================\n")
		sb.WriteString("// Feature Metadata\n")
		sb.WriteString("// =============================================================================\n\n")

		sb.WriteString("// FeatureMeta contains metadata about a feature\n")
		sb.WriteString("type FeatureMeta struct {\n")
		sb.WriteString("\tID        string\n")
		sb.WriteString("\tIsMetered bool\n")
		sb.WriteString("}\n\n")

		sb.WriteString("// AllFeatures contains metadata for all features\n")
		sb.WriteString("var AllFeatures = []FeatureMeta{\n")
		for _, feature := range allFeatures {
			sb.WriteString(fmt.Sprintf("\t{ID: %q, IsMetered: %v},\n", feature.ID, feature.IsMetered))
		}
		sb.WriteString("}\n\n")

		sb.WriteString("// GetFeatureMeta returns the metadata for a feature by ID\n")
		sb.WriteString("func GetFeatureMeta(featureID string) *FeatureMeta {\n")
		sb.WriteString("\tfor _, f := range AllFeatures {\n")
		sb.WriteString("\t\tif f.ID == featureID {\n")
		sb.WriteString("\t\t\treturn &f\n")
		sb.WriteString("\t\t}\n")
		sb.WriteString("\t}\n")
		sb.WriteString("\treturn nil\n")
		sb.WriteString("}\n\n")

		sb.WriteString("// IsMeteredFeature returns true if the feature is metered\n")
		sb.WriteString("func IsMeteredFeature(featureID string) bool {\n")
		sb.WriteString("\tmeta := GetFeatureMeta(featureID)\n")
		sb.WriteString("\treturn meta != nil && meta.IsMetered\n")
		sb.WriteString("}\n\n")
	}

	// ==========================================================================
	// Client Wrapper
	// ==========================================================================
	sb.WriteString("// =============================================================================\n")
	sb.WriteString("// Client\n")
	sb.WriteString("// =============================================================================\n\n")

	sb.WriteString("// Client is the typed Portcall client with methods for your plans and features\n")
	sb.WriteString("type Client struct {\n")
	sb.WriteString("\tSDK      *portcallsdk.Client\n")
	sb.WriteString("\tCheck    *CheckHelpers\n")
	sb.WriteString("\tCheckout *CheckoutHelpers\n")
	sb.WriteString("\tRecord   *RecordHelpers\n")
	sb.WriteString("}\n\n")

	sb.WriteString("// New creates a new typed Portcall client.\n")
	sb.WriteString("// It reads PC_API_SECRET from environment and uses the BaseURL configured at generation time.\n")
	sb.WriteString("func New() *Client {\n")
	sb.WriteString("\treturn NewWithConfig(portcallsdk.Config{\n")
	sb.WriteString("\t\tAPIKey:  os.Getenv(\"PC_API_SECRET\"),\n")
	sb.WriteString("\t\tBaseURL: BaseURL,\n")
	sb.WriteString("\t})\n")
	sb.WriteString("}\n\n")

	sb.WriteString("// NewWithConfig creates a new typed Portcall client with custom configuration.\n")
	sb.WriteString("func NewWithConfig(config portcallsdk.Config) *Client {\n")
	sb.WriteString("\tif config.BaseURL == \"\" {\n")
	sb.WriteString("\t\tconfig.BaseURL = BaseURL\n")
	sb.WriteString("\t}\n")
	sb.WriteString("\tsdk := portcallsdk.New(config)\n")
	sb.WriteString("\treturn &Client{\n")
	sb.WriteString("\t\tSDK:      sdk,\n")
	sb.WriteString("\t\tCheck:    newCheckHelpers(sdk),\n")
	sb.WriteString("\t\tCheckout: newCheckoutHelpers(sdk),\n")
	sb.WriteString("\t\tRecord:   newRecordHelpers(sdk),\n")
	sb.WriteString("\t}\n")
	sb.WriteString("}\n\n")

	// ==========================================================================
	// Check Helpers (Entitlements)
	// ==========================================================================
	sb.WriteString("// =============================================================================\n")
	sb.WriteString("// Check Helpers (Entitlements)\n")
	sb.WriteString("// =============================================================================\n\n")

	sb.WriteString("// EntitlementChecker provides a fluent API for checking entitlements\n")
	sb.WriteString("type EntitlementChecker struct {\n")
	sb.WriteString("\tsdk       *portcallsdk.Client\n")
	sb.WriteString("\tfeatureID string\n")
	sb.WriteString("}\n\n")

	sb.WriteString("// For checks the entitlement for a specific user\n")
	sb.WriteString("func (e *EntitlementChecker) For(ctx context.Context, userID string) (*portcallsdk.Entitlement, error) {\n")
	sb.WriteString("\treturn e.sdk.Entitlements.Get(ctx, userID, e.featureID)\n")
	sb.WriteString("}\n\n")

	sb.WriteString("// IsEnabled checks if the entitlement is enabled for a user\n")
	sb.WriteString("func (e *EntitlementChecker) IsEnabled(ctx context.Context, userID string) (bool, error) {\n")
	sb.WriteString("\treturn e.sdk.Entitlements.IsEnabled(ctx, userID, e.featureID)\n")
	sb.WriteString("}\n\n")

	sb.WriteString("// GetRemaining returns the remaining quota for a user\n")
	sb.WriteString("func (e *EntitlementChecker) GetRemaining(ctx context.Context, userID string) (int64, error) {\n")
	sb.WriteString("\treturn e.sdk.Entitlements.GetRemaining(ctx, userID, e.featureID)\n")
	sb.WriteString("}\n\n")

	sb.WriteString("// GetUsage returns the full usage status for a user\n")
	sb.WriteString("func (e *EntitlementChecker) GetUsage(ctx context.Context, userID string) (*portcallsdk.QuotaStatus, error) {\n")
	sb.WriteString("\treturn e.sdk.Entitlements.GetUsage(ctx, userID, e.featureID)\n")
	sb.WriteString("}\n\n")

	// CheckHelpers struct
	sb.WriteString("// CheckHelpers provides typed entitlement checking for each feature\n")
	sb.WriteString("type CheckHelpers struct {\n")
	for _, feature := range allFeatures {
		identifier := toGoIdentifier(feature.ID)
		sb.WriteString(fmt.Sprintf("\t%s *EntitlementChecker\n", identifier))
	}
	sb.WriteString("}\n\n")

	sb.WriteString("func newCheckHelpers(sdk *portcallsdk.Client) *CheckHelpers {\n")
	sb.WriteString("\treturn &CheckHelpers{\n")
	for _, feature := range allFeatures {
		identifier := toGoIdentifier(feature.ID)
		sb.WriteString(fmt.Sprintf("\t\t%s: &EntitlementChecker{sdk: sdk, featureID: %q},\n", identifier, feature.ID))
	}
	sb.WriteString("\t}\n")
	sb.WriteString("}\n\n")

	// ==========================================================================
	// Checkout Helpers
	// ==========================================================================
	sb.WriteString("// =============================================================================\n")
	sb.WriteString("// Checkout Helpers\n")
	sb.WriteString("// =============================================================================\n\n")

	sb.WriteString("// CheckoutCreator provides a fluent API for creating checkout sessions\n")
	sb.WriteString("type CheckoutCreator struct {\n")
	sb.WriteString("\tsdk    *portcallsdk.Client\n")
	sb.WriteString("\tplanID string\n")
	sb.WriteString("}\n\n")

	sb.WriteString("// For creates a checkout session for a specific user\n")
	sb.WriteString("func (c *CheckoutCreator) For(ctx context.Context, userID, successURL, cancelURL string) (*portcallsdk.CheckoutSession, error) {\n")
	sb.WriteString("\treturn c.sdk.CheckoutSessions.Create(ctx, portcallsdk.CreateCheckoutSessionInput{\n")
	sb.WriteString("\t\tUserID:      userID,\n")
	sb.WriteString("\t\tPlanID:      c.planID,\n")
	sb.WriteString("\t\tRedirectURL: successURL,\n")
	sb.WriteString("\t\tCancelURL:   cancelURL,\n")
	sb.WriteString("\t})\n")
	sb.WriteString("}\n\n")

	// CheckoutHelpers struct
	sb.WriteString("// CheckoutHelpers provides typed checkout session creation for each plan\n")
	sb.WriteString("type CheckoutHelpers struct {\n")
	for _, plan := range appData.Plans {
		identifier := toGoIdentifier(plan.Name)
		sb.WriteString(fmt.Sprintf("\t%s *CheckoutCreator\n", identifier))
	}
	sb.WriteString("}\n\n")

	sb.WriteString("func newCheckoutHelpers(sdk *portcallsdk.Client) *CheckoutHelpers {\n")
	sb.WriteString("\treturn &CheckoutHelpers{\n")
	for _, plan := range appData.Plans {
		identifier := toGoIdentifier(plan.Name)
		sb.WriteString(fmt.Sprintf("\t\t%s: &CheckoutCreator{sdk: sdk, planID: %q},\n", identifier, plan.ID))
	}
	sb.WriteString("\t}\n")
	sb.WriteString("}\n\n")

	// ==========================================================================
	// Record Helpers (Meter Events)
	// ==========================================================================
	sb.WriteString("// =============================================================================\n")
	sb.WriteString("// Record Helpers (Meter Events)\n")
	sb.WriteString("// =============================================================================\n\n")

	sb.WriteString("// UsageRecorder provides a fluent API for recording usage\n")
	sb.WriteString("type UsageRecorder struct {\n")
	sb.WriteString("\tsdk       *portcallsdk.Client\n")
	sb.WriteString("\tfeatureID string\n")
	sb.WriteString("}\n\n")

	sb.WriteString("// For records usage for a specific user\n")
	sb.WriteString("func (r *UsageRecorder) For(ctx context.Context, userID string, quantity int64) error {\n")
	sb.WriteString("\treturn r.sdk.MeterEvents.Record(ctx, userID, r.featureID, quantity)\n")
	sb.WriteString("}\n\n")

	sb.WriteString("// Increment increments usage by 1 for a specific user\n")
	sb.WriteString("func (r *UsageRecorder) Increment(ctx context.Context, userID string) error {\n")
	sb.WriteString("\treturn r.sdk.MeterEvents.Increment(ctx, userID, r.featureID)\n")
	sb.WriteString("}\n\n")

	// RecordHelpers struct (only metered features)
	sb.WriteString("// RecordHelpers provides typed usage recording for each metered feature\n")
	sb.WriteString("type RecordHelpers struct {\n")
	for _, feature := range appData.MeteredFeatures {
		identifier := toGoIdentifier(feature.ID)
		sb.WriteString(fmt.Sprintf("\t%s *UsageRecorder\n", identifier))
	}
	sb.WriteString("}\n\n")

	sb.WriteString("func newRecordHelpers(sdk *portcallsdk.Client) *RecordHelpers {\n")
	sb.WriteString("\treturn &RecordHelpers{\n")
	for _, feature := range appData.MeteredFeatures {
		identifier := toGoIdentifier(feature.ID)
		sb.WriteString(fmt.Sprintf("\t\t%s: &UsageRecorder{sdk: sdk, featureID: %q},\n", identifier, feature.ID))
	}
	sb.WriteString("\t}\n")
	sb.WriteString("}\n")

	return sb.String()
}

func generateYAMLConfig(appData *AppData) string {
	var sb strings.Builder

	sb.WriteString("# Portcall Configuration\n")
	sb.WriteString("# Generated by portcall-gen\n")
	sb.WriteString("# \n")
	sb.WriteString("# This file provides a reference of all plans and features in your Portcall app.\n")
	sb.WriteString("# DO NOT EDIT MANUALLY - run 'portcall-gen generate' to regenerate.\n")
	sb.WriteString("#\n")
	sb.WriteString("# For more information, visit: https://useportcall.com/docs\n\n")

	sb.WriteString("version: \"1.0\"\n")
	sb.WriteString(fmt.Sprintf("generated_at: %q\n\n", time.Now().Format(time.RFC3339)))

	// Plans
	sb.WriteString("plans:\n")
	for _, plan := range appData.Plans {
		sb.WriteString(fmt.Sprintf("  - id: %q\n", plan.ID))
		sb.WriteString(fmt.Sprintf("    name: %q\n", plan.Name))
		sb.WriteString(fmt.Sprintf("    status: %q\n", plan.Status))
		sb.WriteString(fmt.Sprintf("    currency: %q\n", plan.Currency))
		sb.WriteString(fmt.Sprintf("    interval: %q\n", plan.Interval))
		sb.WriteString(fmt.Sprintf("    interval_count: %d\n", plan.IntervalCount))
		sb.WriteString(fmt.Sprintf("    trial_period_days: %d\n", plan.TrialPeriodDays))
		sb.WriteString(fmt.Sprintf("    is_free: %v\n", plan.IsFree))
		if len(plan.Features) > 0 {
			sb.WriteString("    features:\n")
			for _, f := range plan.Features {
				sb.WriteString(fmt.Sprintf("      - feature_id: %q\n", f.FeatureID))
				sb.WriteString(fmt.Sprintf("        quota: %d\n", f.Quota))
				sb.WriteString(fmt.Sprintf("        interval: %q\n", f.Interval))
			}
		}
		sb.WriteString("\n")
	}

	// Features
	sb.WriteString("features:\n")
	sb.WriteString("  standard:\n")
	for _, f := range appData.Features {
		sb.WriteString(fmt.Sprintf("    - id: %q\n", f.ID))
		sb.WriteString(fmt.Sprintf("      is_metered: %v\n", f.IsMetered))
	}
	sb.WriteString("  metered:\n")
	for _, f := range appData.MeteredFeatures {
		sb.WriteString(fmt.Sprintf("    - id: %q\n", f.ID))
		sb.WriteString(fmt.Sprintf("      is_metered: %v\n", f.IsMetered))
	}

	return sb.String()
}

func writeFiles(outputDir, goContent, yamlContent string) ([]string, error) {
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	var files []string

	// Write Go file
	goPath := filepath.Join(outputDir, "generated.go")
	if err := os.WriteFile(goPath, []byte(goContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write Go file: %w", err)
	}
	files = append(files, goPath)

	// Write YAML file if provided
	if yamlContent != "" {
		yamlPath := filepath.Join(outputDir, "portcall.yaml")
		if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
			return nil, fmt.Errorf("failed to write YAML file: %w", err)
		}
		files = append(files, yamlPath)
	}

	return files, nil
}
