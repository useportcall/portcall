package initcmd

import (
	"fmt"
	"strings"
)

func normalizeDNSOptions(opts Options) (Options, error) {
	opts.DNSProvider = strings.ToLower(strings.TrimSpace(opts.DNSProvider))
	if opts.DNSProvider == "" {
		opts.DNSProvider = "manual"
	}
	switch opts.DNSProvider {
	case "manual", "cloudflare":
		return opts, nil
	default:
		return opts, fmt.Errorf("unsupported --dns-provider %q (supported: manual|cloudflare)", opts.DNSProvider)
	}
}
