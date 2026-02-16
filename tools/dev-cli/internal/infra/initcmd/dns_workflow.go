package initcmd

import (
	"strings"
)

func runDNSWorkflow(plan Plan, opts Options, deps Deps) {
	printDNSChecklist(plan, deps)
	if opts.DryRun || isPlaceholderDomain(plan.Domain) {
		return
	}
	runCloudflare, autoApply := resolveCloudflareDNSIntent(opts, deps)
	if !runCloudflare {
		return
	}
	deps.Section("Cloudflare DNS")
	reportCloudflareCLIStatus(deps)
	target := resolveIngressTarget(deps)
	if strings.TrimSpace(target) == "" {
		deps.Warn("Cannot update DNS yet: ingress load balancer target is empty.")
		return
	}
	token := resolveCloudflareToken(deps)
	if token == "" {
		deps.Warn("Cloudflare token unavailable; skipping Cloudflare DNS check/apply.")
		return
	}
	runCloudflareDNS(plan, opts, token, target, autoApply, deps)
}

func resolveCloudflareDNSIntent(opts Options, deps Deps) (run bool, autoApply bool) {
	if opts.DNSProvider == "cloudflare" {
		return true, opts.DNSAuto
	}
	if opts.Yes || !deps.IsInteractive() {
		return false, false
	}
	if !deps.AskYesNo("Check Cloudflare CLI and DNS records now? [y/N]: ", false) {
		return false, false
	}
	return true, deps.AskYesNo("Create/update missing records automatically? [y/N]: ", false)
}
