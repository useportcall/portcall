package doauth

import "strings"

func doctlReadChecks() []doctlCheck {
	return []doctlCheck{
		{Name: "doctl account", Args: []string{"account", "get", "--output", "json"}},
		{Name: "kubernetes read", Args: []string{"kubernetes", "options", "regions", "--output", "json"}},
		{Name: "databases read", Args: []string{"databases", "options", "regions", "--output", "json"}},
		{Name: "registry read", Args: []string{"registry", "options", "available-regions", "--output", "json"}},
	}
}

func doctlWriteProbes() []doctlCheck {
	return []doctlCheck{
		{Name: "kubernetes write", Args: []string{"kubernetes", "cluster", "node-pool", "create", "00000000-0000-0000-0000-000000000000", "--name", "dev-cli-probe", "--size", "s-1vcpu-2gb", "--count", "1"}},
		{Name: "databases write", Args: []string{"databases", "db", "create", "00000000-0000-0000-0000-000000000000", "dev_cli_probe"}},
		{Name: "registry write", Args: []string{"registry", "create", "dev-cli-permission-probe", "--region", "zz-invalid-region", "--subscription-tier", "basic"}},
		{Name: "vpc write", Args: []string{"vpcs", "create", "--name", "dev-cli-permission-probe", "--region", "zz-invalid-region", "--ip-range", "10.250.0.0/24"}},
	}
}

func hasPermissionDenied(out string) bool {
	lower := strings.ToLower(out)
	for _, hint := range []string{"forbidden", "permission", "scope", "unauthorized", "not authorized"} {
		if strings.Contains(lower, hint) {
			return true
		}
	}
	return false
}

func hasExpectedProbeFailure(out string) bool {
	lower := strings.ToLower(out)
	for _, hint := range []string{"not found", "invalid", "must be one of", "unprocessable", "bad request", "already exists", "cannot"} {
		if strings.Contains(lower, hint) {
			return true
		}
	}
	return false
}

func formatProbeError(out string) string {
	trimmed := strings.Join(strings.Fields(strings.TrimSpace(out)), " ")
	if trimmed == "" {
		return "no diagnostic output"
	}
	if len(trimmed) <= 280 {
		return trimmed
	}
	return trimmed[:277] + "..."
}
