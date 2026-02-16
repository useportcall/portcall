package initcmd

import (
	"fmt"
	"os/exec"
	"strings"
)

func checkPrereqs(plan Plan, opts Options, deps Deps) error {
	required, hasDigitalOcean := requiredTools(plan, opts)
	missing := missingTools(required)
	if len(missing) > 0 {
		return fmt.Errorf("missing required infra CLI(s):\n- %s", strings.Join(missing, "\n- "))
	}
	if !hasDigitalOcean {
		return nil
	}
	token, err := resolveDOTokenForInit(deps)
	if err != nil {
		return err
	}
	if deps.VerifyDOAccess == nil {
		return nil
	}
	if err := deps.VerifyDOAccess(token); err != nil {
		return fmt.Errorf("digitalocean preflight failed: %w\nfix your doctl/token setup and retry:\n- doctl auth init --context portcall\n- doctl auth switch --context portcall\n- export DIGITALOCEAN_TOKEN=<token>\n- doctl account get", err)
	}
	return nil
}

func requiredTools(plan Plan, opts Options) ([]string, bool) {
	required := []string{"terraform"}
	hasDigitalOcean := false
	for _, p := range plan.Providers {
		if p == "digitalocean" {
			required = append(required, "doctl")
			hasDigitalOcean = true
		}
	}
	if !opts.DryRun {
		required = append(required, "kubectl", "helm")
	}
	return required, hasDigitalOcean
}

func missingTools(required []string) []string {
	missing := []string{}
	for _, name := range uniqueToolNames(required) {
		if _, err := exec.LookPath(name); err != nil {
			missing = append(missing, missingToolHint(name))
		}
	}
	return missing
}

func uniqueToolNames(names []string) []string {
	out, seen := []string{}, map[string]bool{}
	for _, n := range names {
		if !seen[n] {
			seen[n] = true
			out = append(out, n)
		}
	}
	return out
}

func missingToolHint(name string) string {
	hints := map[string]string{
		"terraform": "terraform (https://developer.hashicorp.com/terraform/install)",
		"doctl":     "doctl (https://docs.digitalocean.com/reference/doctl/how-to/install/)",
		"kubectl":   "kubectl (https://kubernetes.io/docs/tasks/tools/)",
		"helm":      "helm (https://helm.sh/docs/intro/install/)",
	}
	if h, ok := hints[name]; ok {
		return h
	}
	return name
}
