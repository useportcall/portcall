package main

import "fmt"

// printSnapshotList renders a human-friendly list of targets and groups.
func printSnapshotList() {
	fmt.Println("Available snapshot targets:")
	fmt.Println()
	lastMode := SnapshotMode("")
	for _, t := range snapshotTargets {
		if t.Mode != lastMode {
			fmt.Printf("  [%s]\n", t.Mode)
			lastMode = t.Mode
		}
		fmt.Printf("    %-28s %s\n", t.Name, t.Description)
	}
	fmt.Println()
	fmt.Println("Groups (expand to multiple targets):")
	for name, members := range snapshotGroups {
		fmt.Printf("  %-28s %v\n", name, members)
	}
}
