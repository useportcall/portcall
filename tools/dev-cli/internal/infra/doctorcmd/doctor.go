package doctorcmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	infraState "github.com/useportcall/portcall/tools/dev-cli/internal/infra/state"
)

type doctorCheck struct {
	Name     string
	Passed   bool
	Severity string
	Message  string
}

type Deps struct {
	EnsureRootDir         func() error
	ResolveClusterContext func(cluster string) string
	GetInfraClusterState  func(cluster string) (infraState.ClusterState, bool)
	RunCmdOut             func(name string, args ...string) (string, error)
	OK                    func(msg string, args ...any)
	Warn                  func(msg string, args ...any)
	Fail                  func(msg string, args ...any)
}

func New(deps Deps) *cobra.Command {
	var opts struct {
		cluster   string
		namespace string
		localOnly bool
		dryRun    bool
	}
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Validate infra state and cluster wiring",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := deps.EnsureRootDir(); err != nil {
				return err
			}
			checks := []doctorCheck{}
			checks = append(checks, checkDoctorTools()...)
			ctx := deps.ResolveClusterContext(opts.cluster)
			checks = append(checks, checkStateForCluster(deps.GetInfraClusterState, opts.cluster, ctx))
			if opts.localOnly || opts.dryRun {
				deps.Warn("doctor local mode enabled: skipping cluster connectivity, runtime config, secrets, and legacy resource checks")
			} else {
				checks = append(checks, checkClusterConnectivity(deps.RunCmdOut, ctx, opts.namespace)...)
				checks = append(checks, checkRuntimeConfig(deps.RunCmdOut, ctx, opts.namespace)...)
				checks = append(checks, checkRequiredSecrets(deps.RunCmdOut, ctx, opts.namespace)...)
				checks = append(checks, checkLegacyResources(deps.RunCmdOut, ctx, opts.namespace)...)
			}
			fails := 0
			for _, check := range checks {
				if check.Passed {
					deps.OK("[OK] %s: %s", check.Name, check.Message)
					continue
				}
				if check.Severity == "warn" {
					deps.Warn("[WARN] %s: %s", check.Name, check.Message)
					continue
				}
				deps.Fail("[FAIL] %s: %s", check.Name, check.Message)
				fails++
			}
			if fails > 0 {
				return fmt.Errorf("infra doctor found %d blocking issue(s)", fails)
			}
			deps.OK("infra doctor passed")
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.cluster, "cluster", "digitalocean", "Cluster alias")
	cmd.Flags().StringVarP(&opts.namespace, "namespace", "n", "portcall", "Kubernetes namespace")
	cmd.Flags().BoolVar(&opts.localOnly, "local-only", false, "Run only local checks (tools + saved state), skip cluster calls")
	cmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "Preview doctor checks without cluster calls")
	return cmd
}

func checkDoctorTools() []doctorCheck {
	names := []string{"terraform", "doctl", "kubectl", "helm", "docker"}
	out := []doctorCheck{}
	for _, name := range names {
		_, err := exec.LookPath(name)
		out = append(out, doctorCheck{Name: "tool:" + name, Passed: err == nil, Severity: "fail", Message: "available in PATH"})
		if err != nil {
			out[len(out)-1].Message = "missing from PATH"
		}
	}
	return out
}

func checkStateForCluster(getInfraClusterState func(string) (infraState.ClusterState, bool), cluster, ctx string) doctorCheck {
	cfg, ok := getInfraClusterState(cluster)
	if !ok {
		return doctorCheck{Name: "state", Passed: false, Severity: "warn", Message: "no .dev-cli.infra.json entry for alias"}
	}
	if cfg.Values != "" {
		if _, err := os.Stat(cfg.Values); err != nil {
			return doctorCheck{Name: "state", Passed: false, Severity: "warn", Message: "saved values file not found"}
		}
	}
	if cfg.Context != "" && cfg.Context != ctx {
		return doctorCheck{Name: "state", Passed: false, Severity: "warn", Message: "saved context differs from resolved context"}
	}
	return doctorCheck{Name: "state", Passed: true, Severity: "warn", Message: "infra state entry found"}
}

func checkClusterConnectivity(runCmdOut func(name string, args ...string) (string, error), ctx, namespace string) []doctorCheck {
	checks := []doctorCheck{}
	if _, err := runCmdOut("kubectl", "--context", ctx, "get", "nodes", "--request-timeout=5s"); err != nil {
		return append(checks, doctorCheck{Name: "cluster", Passed: false, Severity: "fail", Message: "cannot reach cluster context " + ctx})
	}
	checks = append(checks, doctorCheck{Name: "cluster", Passed: true, Severity: "fail", Message: "connected to context " + ctx})
	if _, err := runCmdOut("kubectl", "--context", ctx, "get", "namespace", namespace); err != nil {
		checks = append(checks, doctorCheck{Name: "namespace", Passed: false, Severity: "fail", Message: "namespace missing: " + namespace})
	} else {
		checks = append(checks, doctorCheck{Name: "namespace", Passed: true, Severity: "fail", Message: "namespace exists: " + namespace})
	}
	return checks
}

func checkRequiredSecrets(runCmdOut func(name string, args ...string) (string, error), ctx, namespace string) []doctorCheck {
	checks := []doctorCheck{}
	checks = append(checks, checkSecretKeys(runCmdOut, ctx, namespace, "portcall-secrets", []string{"DATABASE_URL", "POSTGRES_USER", "POSTGRES_PASSWORD", "AES_ENCRYPTION_KEY", "KC_BOOTSTRAP_ADMIN_USERNAME", "KC_BOOTSTRAP_ADMIN_PASSWORD", "KC_DB_USERNAME", "KC_DB_PASSWORD", "REDIS_PASSWORD"}))
	checks = append(checks, checkSecretKeys(runCmdOut, ctx, namespace, "spaces-credentials", []string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"}))
	return checks
}

func checkRuntimeConfig(runCmdOut func(name string, args ...string) (string, error), ctx, namespace string) []doctorCheck {
	out := []doctorCheck{}
	for _, key := range []string{"S3_ENDPOINT", "S3_REGION", "REDIS_ADDR"} {
		val, err := runCmdOut("kubectl", "--context", ctx, "get", "configmap", "portcall-config", "-n", namespace, "-o", "jsonpath={.data."+key+"}")
		if err != nil || val == "" {
			out = append(out, doctorCheck{Name: "configmap:portcall-config", Passed: false, Severity: "fail", Message: "missing key " + key})
			continue
		}
		out = append(out, doctorCheck{Name: "configmap:portcall-config", Passed: true, Severity: "fail", Message: key + " present"})
	}
	return out
}

func checkSecretKeys(runCmdOut func(name string, args ...string) (string, error), ctx, namespace, name string, keys []string) doctorCheck {
	for _, key := range keys {
		out, err := runCmdOut("kubectl", "--context", ctx, "get", "secret", name, "-n", namespace, "-o", "jsonpath={.data."+key+"}")
		if err != nil || out == "" {
			return doctorCheck{Name: "secret:" + name, Passed: false, Severity: "fail", Message: "missing key " + key}
		}
	}
	return doctorCheck{Name: "secret:" + name, Passed: true, Severity: "fail", Message: "required keys present"}
}

func checkLegacyResources(runCmdOut func(name string, args ...string) (string, error), ctx, namespace string) []doctorCheck {
	checks := []doctorCheck{}
	if _, err := runCmdOut("kubectl", "--context", ctx, "get", "deployment", "smtp-relay", "-n", namespace); err == nil {
		checks = append(checks, doctorCheck{Name: "legacy:smtp-relay", Passed: false, Severity: "warn", Message: "deprecated deployment exists; run infra cleanup legacy"})
	} else {
		checks = append(checks, doctorCheck{Name: "legacy:smtp-relay", Passed: true, Severity: "warn", Message: "no deprecated smtp-relay deployment"})
	}
	return checks
}
