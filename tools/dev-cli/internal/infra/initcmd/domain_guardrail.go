package initcmd

import (
	"fmt"
	"strings"
)

func confirmDomainGuardrail(plan Plan, opts Options, deps Deps) error {
	if opts.DryRun || !isPlaceholderDomain(plan.Domain) {
		return nil
	}
	deps.Section("Domain Guardrail")
	deps.Warn("You are using placeholder domain %q", plan.Domain)
	deps.Plain("Ingress endpoints will be unreachable until DNS is configured for a real domain.")
	if opts.Yes {
		return fmt.Errorf("refusing non-dry-run apply with --yes and placeholder domain %q; set --domain <your-domain>", plan.Domain)
	}
	token := "placeholder domain"
	deps.Info("Recommended: rerun with --domain <your-domain>")
	if deps.AskText("Type \""+token+"\" to continue anyway: ", "") != token {
		return fmt.Errorf("init canceled")
	}
	if !deps.AskYesNo("Proceed with placeholder domain? [y/N]: ", false) {
		return fmt.Errorf("init canceled")
	}
	return nil
}

func isPlaceholderDomain(domain string) bool {
	normalized := strings.ToLower(strings.TrimSpace(domain))
	return normalized == "" || normalized == "example.com" || strings.HasSuffix(normalized, ".example.com")
}
