package initcmd

import (
	"fmt"

	"github.com/useportcall/portcall/tools/dev-cli/internal/infra/dnscloudflare"
)

func runCloudflareDNS(plan Plan, opts Options, token, target string, autoApply bool, deps Deps) {
	client := dnscloudflare.NewClient(token)
	zoneID, err := client.ResolveZoneID(plan.Domain, opts.CloudflareZoneID)
	if err != nil {
		deps.Warn("Cloudflare zone lookup failed: %v", err)
		return
	}
	deps.Info("Cloudflare zone ID: %s", zoneID)
	hosts := ingressHostsForDomain(plan.Domain)
	results, err := client.EnsureRecords(zoneID, hosts, target, autoApply)
	if err != nil {
		deps.Warn("Cloudflare DNS ensure failed: %v", err)
		return
	}
	printCloudflareResults(results, autoApply, deps)
	verifyPublicDNS(hosts, deps)
}

func printCloudflareResults(results []dnscloudflare.RecordOutcome, autoApply bool, deps Deps) {
	for _, item := range results {
		line := fmt.Sprintf("%s: %s (%s)", item.Host, item.Status, item.Detail)
		switch item.Status {
		case "ok":
			deps.OK(line)
		case "created", "updated":
			deps.OK(line)
		case "missing", "mismatch":
			if autoApply {
				deps.Warn(line)
			} else {
				deps.Warn(line + " | rerun with --dns-provider cloudflare --dns-auto")
			}
		default:
			deps.Warn(line)
		}
	}
}
