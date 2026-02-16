package snapshot

import (
	"fmt"
	"strings"
)

func lookupTargets(names []string) ([]Target, error) {
	var resolved []Target
	seen := map[string]bool{}
	for _, name := range names {
		found, err := resolveArg(name)
		if err != nil {
			return nil, err
		}
		for _, t := range found {
			if !seen[t.Name] {
				seen[t.Name] = true
				resolved = append(resolved, t)
			}
		}
	}
	return resolved, nil
}

func resolveArg(name string) ([]Target, error) {
	if members, ok := groups[name]; ok {
		return lookupTargets(members)
	}
	for _, t := range targets {
		if t.Name == name {
			return []Target{t}, nil
		}
	}
	return nil, fmt.Errorf("unknown snapshot target or group: %q (use --list)", name)
}

func filterByMode(mode Mode) []Target {
	var out []Target
	for _, t := range targets {
		if t.Mode == mode {
			out = append(out, t)
		}
	}
	return out
}

func buildGrepPattern(ts []Target) string {
	var parts []string
	for _, t := range ts {
		parts = append(parts, t.Grep)
	}
	return strings.Join(parts, "|")
}

func availableModes() []Mode {
	seen := map[Mode]bool{}
	var modes []Mode
	for _, t := range targets {
		if !seen[t.Mode] {
			seen[t.Mode] = true
			modes = append(modes, t.Mode)
		}
	}
	return modes
}
