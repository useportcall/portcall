package snapshot

import "fmt"

func printList() {
	fmt.Println("Available snapshot targets:")
	fmt.Println()
	lastMode := Mode("")
	for _, t := range targets {
		if t.Mode != lastMode {
			fmt.Printf("  [%s]\n", t.Mode)
			lastMode = t.Mode
		}
		fmt.Printf("    %-28s %s\n", t.Name, t.Description)
	}
	fmt.Println()
	fmt.Println("Groups (expand to multiple targets):")
	for name, members := range groups {
		fmt.Printf("  %-28s %v\n", name, members)
	}
}
