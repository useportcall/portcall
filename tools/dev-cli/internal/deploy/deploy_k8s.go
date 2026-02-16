package deploy

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/useportcall/portcall/tools/dev-cli/internal/infra"
)

var versionRe = regexp.MustCompile(`v\d+\.\d+\.\d+`)

func switchCluster(cluster string) error {
	_ = ensureRootDir()
	return runCmd("kubectl", "config", "use-context", resolveClusterContext(cluster))
}

func resolveClusterContext(cluster string) string {
	_ = ensureRootDir()
	return infra.ResolveClusterContext(rootDir, cluster)
}

func currentVersion(app string) string {
	name := deployName(app)
	out, err := runCmdOut("kubectl", "get", "deployment", name, "-n", "portcall", "-o", "jsonpath={.spec.template.spec.containers[0].image}")
	if err != nil {
		return "v0.0.0"
	}
	if v := versionRe.FindString(out); v != "" {
		return v
	}
	return "v0.0.0"
}

func collectVersions() map[string]string {
	out := map[string]string{}
	for _, app := range deployApps {
		if app.ChartOnly {
			continue
		}
		if v := currentVersion(app.Name); v != "v0.0.0" {
			out[helmValueName(app.Name)] = v
		}
	}
	return out
}

func versionSetArgs(versions map[string]string) []string {
	keys := make([]string, 0, len(versions))
	for k := range versions {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	args := []string{}
	for _, k := range keys {
		args = append(args, "--set", fmt.Sprintf("%s.image.tag=%s", k, versions[k]))
	}
	return args
}
