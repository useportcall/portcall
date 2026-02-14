package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// collectTargets merges explicit names and mode filter into a target list.
func collectTargets(args []string, mode string) ([]SnapshotTarget, error) {
	var targets []SnapshotTarget

	if mode != "" {
		m := SnapshotMode(mode)
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

func dedup(targets []SnapshotTarget) []SnapshotTarget {
	seen := map[string]bool{}
	var out []SnapshotTarget
	for _, t := range targets {
		if !seen[t.Name] {
			seen[t.Name] = true
			out = append(out, t)
		}
	}
	return out
}

// execSnapshot runs npx playwright test with the appropriate --grep.
func execSnapshot(targets []SnapshotTarget, live, headed bool) error {
	root, err := findRootDir()
	if err != nil {
		return err
	}

	grep := buildGrepPattern(targets)
	playwrightArgs := []string{"playwright", "test", "--grep", grep}
	if headed {
		playwrightArgs = append(playwrightArgs, "--headed")
	}

	// Print what will run.
	fmt.Printf("Running %d snapshot(s):\n", len(targets))
	for _, t := range targets {
		fmt.Printf("  • %-28s [%s]\n", t.Name, t.Mode)
	}
	fmt.Println()

	env := os.Environ()
	if live {
		env = loadDiscordEnv(root, env)
		fmt.Println("Mode: live (uploading to Discord)")
	} else {
		fmt.Println("Mode: local (saving artifacts only)")
	}

	c := exec.Command("npx", playwrightArgs...)
	c.Dir = root + "/e2etest/browser"
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Env = env
	if err := c.Run(); err != nil {
		return fmt.Errorf("playwright snapshot run failed: %w", err)
	}

	fmt.Println("✅ Snapshots completed!")
	return nil
}

// loadDiscordEnv sources .envs files to pick up webhook URLs.
func loadDiscordEnv(root string, base []string) []string {
	envFiles := []string{
		"apps/dashboard/.envs",
		"apps/api/.envs",
		"apps/billing/.envs",
	}
	for _, f := range envFiles {
		path := root + "/" + f
		data, err := os.ReadFile(path)
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
