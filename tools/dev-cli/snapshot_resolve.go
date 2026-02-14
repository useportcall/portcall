package main

import (
	"fmt"
	"strings"
)

// lookupTargets resolves CLI arguments into a flat list of SnapshotTarget.
// Each arg can be a target name or a group name.
func lookupTargets(names []string) ([]SnapshotTarget, error) {
	var resolved []SnapshotTarget
	seen := map[string]bool{}
	for _, name := range names {
		targets, err := resolveArg(name)
		if err != nil {
			return nil, err
		}
		for _, t := range targets {
			if !seen[t.Name] {
				seen[t.Name] = true
				resolved = append(resolved, t)
			}
		}
	}
	return resolved, nil
}

func resolveArg(name string) ([]SnapshotTarget, error) {
	// Check group first.
	if members, ok := snapshotGroups[name]; ok {
		return lookupTargets(members)
	}
	// Then individual target.
	for _, t := range snapshotTargets {
		if t.Name == name {
			return []SnapshotTarget{t}, nil
		}
	}
	return nil, fmt.Errorf("unknown snapshot target or group: %q (use --list)", name)
}

// filterByMode returns targets matching the given mode.
func filterByMode(mode SnapshotMode) []SnapshotTarget {
	var out []SnapshotTarget
	for _, t := range snapshotTargets {
		if t.Mode == mode {
			out = append(out, t)
		}
	}
	return out
}

// buildGrepPattern combines multiple targets into a single Playwright
// --grep regex using alternation.
func buildGrepPattern(targets []SnapshotTarget) string {
	var parts []string
	for _, t := range targets {
		parts = append(parts, t.Grep)
	}
	return strings.Join(parts, "|")
}

// availableModes returns the distinct modes in the registry.
func availableModes() []SnapshotMode {
	seen := map[SnapshotMode]bool{}
	var modes []SnapshotMode
	for _, t := range snapshotTargets {
		if !seen[t.Mode] {
			seen[t.Mode] = true
			modes = append(modes, t.Mode)
		}
	}
	return modes
}
