package initcmd

import "fmt"

var ingressHostPrefixes = []string{
	"admin",
	"dashboard",
	"api",
	"auth",
	"quote",
	"checkout",
	"webhook",
	"file",
}

func ingressHostsForDomain(domain string) []string {
	hosts := make([]string, 0, len(ingressHostPrefixes))
	for _, prefix := range ingressHostPrefixes {
		hosts = append(hosts, fmt.Sprintf("%s.%s", prefix, domain))
	}
	return hosts
}
