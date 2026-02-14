package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readInput(prompt string) string {
	fmt.Print(prompt)
	line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.TrimSpace(line)
}

func selectApps(appInput string) ([]string, error) {
	if strings.TrimSpace(appInput) == "" {
		listDeployApps()
		appInput = readInput("Enter comma-separated names/numbers or 'all': ")
	}
	if strings.TrimSpace(appInput) == "" {
		return nil, fmt.Errorf("no apps selected")
	}
	toks := strings.Split(appInput, ",")
	out := []string{}
	seen := map[string]bool{}
	for _, tok := range toks {
		t := strings.ToLower(strings.TrimSpace(tok))
		if t == "" {
			continue
		}
		if t == "all" {
			out = []string{}
			for _, app := range deployApps {
				out = append(out, app.Name)
			}
			return out, nil
		}
		if n, err := strconv.Atoi(t); err == nil {
			if n < 1 || n > len(deployApps) {
				return nil, fmt.Errorf("invalid app number: %d", n)
			}
			t = deployApps[n-1].Name
		}
		if _, ok := findDeployApp(t); !ok {
			return nil, fmt.Errorf("unknown app: %s", t)
		}
		if !seen[t] {
			out = append(out, t)
			seen[t] = true
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no apps selected")
	}
	return out, nil
}
