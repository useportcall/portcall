package initcmd

import "strings"

func printDNSChecklist(plan Plan, deps Deps) {
	deps.Section("DNS Checklist")
	if isPlaceholderDomain(plan.Domain) {
		deps.Warn("Domain %q is a placeholder. Deploy can proceed, but ingress endpoints will not be reachable.", plan.Domain)
		deps.Warn("Set a real domain before sharing endpoints.")
		return
	}
	target := resolveIngressTarget(deps)
	if target == "" {
		deps.Warn("Ingress load balancer address is not ready yet.")
		deps.Plain("Run after a minute: kubectl get svc ingress-nginx-controller -n ingress-nginx")
	} else {
		deps.Info("Ingress target: %s", target)
	}
	deps.Plain("Create DNS records for:")
	for _, host := range ingressHostsForDomain(plan.Domain) {
		deps.Plain("- %s", host)
	}
	deps.Warn("Apps behind ingress will not be accessible until these records resolve to the ingress target.")
}

func resolveIngressTarget(deps Deps) string {
	ip, err := deps.RunCmdOut("kubectl", "get", "svc", "ingress-nginx-controller", "-n", "ingress-nginx", "-o", "jsonpath={.status.loadBalancer.ingress[0].ip}")
	if err == nil && strings.TrimSpace(ip) != "" {
		return strings.TrimSpace(ip)
	}
	host, err := deps.RunCmdOut("kubectl", "get", "svc", "ingress-nginx-controller", "-n", "ingress-nginx", "-o", "jsonpath={.status.loadBalancer.ingress[0].hostname}")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(host)
}
