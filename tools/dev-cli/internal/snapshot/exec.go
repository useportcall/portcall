package snapshot

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func collectTargets(args []string, mode string) ([]Target, error) {
	var targets []Target
	if mode != "" {
		m := Mode(mode)
		modeTargets := filterByMode(m)
		if len(modeTargets) == 0 {
			return nil, fmt.Errorf("no targets for mode %q; available: %v", mode, availableModes())
		}
		targets = append(targets, modeTargets...)
	}
	if len(args) > 0 {
		explicit, err := lookupTargets(args)
		if err != nil {
			return nil, err
		}
		targets = append(targets, explicit...)
	}
	return dedup(targets), nil
}

func dedup(targets []Target) []Target {
	seen := map[string]bool{}
	var out []Target
	for _, t := range targets {
		if !seen[t.Name] {
			seen[t.Name] = true
			out = append(out, t)
		}
	}
	return out
}

func execSnapshot(findRoot rootFinder, targets []Target, live, headed bool) error {
	root, err := findRoot()
	if err != nil {
		return err
	}
	grep := buildGrepPattern(targets)
	args := []string{"playwright", "test", "--grep", grep}
	if headed {
		args = append(args, "--headed")
	}
	fmt.Printf("Running %d snapshot(s):\n", len(targets))
	for _, t := range targets {
		fmt.Printf("  - %-28s [%s]\n", t.Name, t.Mode)
	}
	env := os.Environ()
	if live {
		env = loadDiscordEnv(root, env)
		fmt.Println("Mode: live (uploading to Discord)")
	} else {
		fmt.Println("Mode: local (saving artifacts only)")
	}
	c := exec.Command("npx", args...)
	c.Dir = root + "/e2etest/browser"
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Env = env
	if err := c.Run(); err != nil {
		return fmt.Errorf("playwright snapshot run failed: %w", err)
	}
	fmt.Println("Snapshots completed.")
	return nil
}

func loadDiscordEnv(root string, base []string) []string {
	files := []string{"apps/dashboard/.envs", "apps/api/.envs", "apps/billing/.envs"}
	for _, f := range files {
		data, err := os.ReadFile(root + "/" + f)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			base = append(base, line)
		}
	}
	return base
}
