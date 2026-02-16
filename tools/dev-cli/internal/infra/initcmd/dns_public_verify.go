package initcmd

import (
	"context"
	"net"
	"time"
)

func verifyPublicDNS(hosts []string, deps Deps) {
	deps.Info("Checking public DNS resolution...")
	for _, host := range hosts {
		status := lookupDNS(host)
		if status == "" {
			deps.Warn("%s: not resolving yet (propagation pending)", host)
			continue
		}
		deps.OK("%s: resolves to %s", host, status)
	}
}

func lookupDNS(host string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	addrs, err := net.DefaultResolver.LookupHost(ctx, host)
	if err != nil || len(addrs) == 0 {
		return ""
	}
	return addrs[0]
}
