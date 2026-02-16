package initcmd

import (
	"fmt"
	"sort"
	"strings"
)

func bootstrapCluster(plan Plan, opts Options, secrets bootstrapSecrets, deps Deps) error {
	if err := deps.RunShell(fmt.Sprintf(
		"kubectl create namespace %s --dry-run=client -o yaml | kubectl apply -f -", plan.Namespace)); err != nil {
		return fmt.Errorf("bootstrap namespace: %w", err)
	}
	if !opts.SkipRegistryLogin {
		if err := deps.RunCmd("doctl", "registry", "login", "--expiry-seconds", "1200"); err != nil {
			return fmt.Errorf("doctl registry login failed: %w", err)
		}
	}
	if err := createRegistrySecret(plan.Namespace, deps); err != nil {
		return err
	}
	if len(secrets.Portcall) > 0 {
		if err := applySecretLiterals(plan.Namespace, "portcall-secrets", secrets.Portcall, deps); err != nil {
			return err
		}
	}
	if len(secrets.Spaces) > 0 {
		if err := applySecretLiterals(plan.Namespace, "spaces-credentials", secrets.Spaces, deps); err != nil {
			return err
		}
	}
	if opts.InstallIngress {
		return installIngressNginx(deps)
	}
	return nil
}

func createRegistrySecret(namespace string, deps Deps) error {
	script := fmt.Sprintf(
		"kubectl create secret generic portcall-registry -n %s "+
			"--type=kubernetes.io/dockerconfigjson "+
			"--from-file=.dockerconfigjson=$HOME/.docker/config.json "+
			"--dry-run=client -o yaml | kubectl apply -f -", namespace)
	if err := deps.RunShell(script); err != nil {
		return fmt.Errorf("create image pull secret: %w", err)
	}
	return nil
}

func installIngressNginx(deps Deps) error {
	if err := deps.RunCmd("helm", "repo", "add", "ingress-nginx", "https://kubernetes.github.io/ingress-nginx"); err != nil {
		return fmt.Errorf("helm repo add ingress-nginx: %w", err)
	}
	if err := deps.RunCmd("helm", "repo", "update"); err != nil {
		return fmt.Errorf("helm repo update: %w", err)
	}
	args := []string{"upgrade", "--install", "ingress-nginx", "ingress-nginx/ingress-nginx",
		"--namespace", "ingress-nginx", "--create-namespace"}
	if err := deps.RunCmd("helm", args...); err != nil {
		return fmt.Errorf("install ingress-nginx: %w", err)
	}
	if err := deps.RunCmd("kubectl", "rollout", "status",
		"deployment/ingress-nginx-controller", "-n", "ingress-nginx", "--timeout=180s"); err != nil {
		return fmt.Errorf("wait ingress-nginx rollout: %w", err)
	}
	return nil
}

func getSecretValue(ns, name, key string, deps Deps) string {
	out, err := deps.RunCmdOut("kubectl", "get", "secret", name, "-n", ns,
		"-o", "jsonpath={.data."+key+"}")
	if err != nil || strings.TrimSpace(out) == "" {
		return ""
	}
	decoded := deps.DecodeB64(strings.TrimSpace(out))
	if decoded == "(decode error)" {
		return ""
	}
	return decoded
}

func ensureStableSecret(ns, name, key string, bytes int, deps Deps) (string, error) {
	if existing := getSecretValue(ns, name, key, deps); existing != "" {
		return existing, nil
	}
	v, err := deps.RandomSecret(bytes)
	if err != nil {
		return "", fmt.Errorf("generate %s: %w", key, err)
	}
	return v, nil
}

func applySecretLiterals(ns, name string, data map[string]string, deps Deps) error {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := []string{fmt.Sprintf("kubectl create secret generic %s -n %s", name, ns)}
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("--from-literal=%s=%s", k, shellSingleQuote(data[k])))
	}
	parts = append(parts, "--dry-run=client -o yaml | kubectl apply -f -")
	if err := deps.RunShell(strings.Join(parts, " ")); err != nil {
		return fmt.Errorf("apply secret %s: %w", name, err)
	}
	return nil
}

func shellSingleQuote(v string) string {
	if v == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(v, "'", "'\\''") + "'"
}
