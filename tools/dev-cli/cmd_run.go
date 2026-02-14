package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var dockerApps, terminalApps, preset string
var listApps bool

func newRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Start the development environment",
		Long:  "Start apps in Docker or local Terminal windows.\n\nExamples:\n  dev-cli run --preset=dashboard\n  dev-cli run --docker=api --terminal=dashboard",
		RunE:  runCommand,
	}
	cmd.Flags().StringVar(&dockerApps, "docker", "", "Comma-separated apps to run in Docker")
	cmd.Flags().StringVar(&terminalApps, "terminal", "", "Comma-separated apps to run in Terminal windows")
	cmd.Flags().StringVar(&preset, "preset", "", "Use a predefined configuration")
	cmd.Flags().IntVar(&parallel, "parallel", 2, "Docker compose parallelism")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be executed")
	cmd.Flags().BoolVar(&listApps, "list", false, "List available apps and presets")
	return cmd
}

func runCommand(cmd *cobra.Command, args []string) error {
	var err error
	rootDir, err = findRootDir()
	if err != nil {
		return fmt.Errorf("could not find project root: %w", err)
	}
	if listApps {
		printAppsAndPresets()
		return nil
	}

	var docker, terminal []string
	if preset != "" {
		p, ok := presets[preset]
		if !ok {
			return fmt.Errorf("unknown preset: %s", preset)
		}
		docker, terminal = p.Docker, p.Terminal
		fmt.Printf("Using preset: %s\n  %s\n\n", preset, p.Description)
	} else if dockerApps != "" || terminalApps != "" {
		if dockerApps != "" {
			docker = parseAppList(dockerApps)
		}
		if terminalApps != "" {
			terminal = parseAppList(terminalApps)
		}
	} else {
		return runInteractive()
	}
	if err := validateApps(docker); err != nil {
		return err
	}
	if err := validateApps(terminal); err != nil {
		return err
	}
	return runEnvironment(docker, terminal)
}

func printAppsAndPresets() {
	fmt.Println("Available Apps:")
	for i, app := range apps {
		extra := ""
		if app.FrontendCmd != "" {
			extra += " (frontend)"
		}
		if app.IsWorker {
			extra += " [worker]"
		}
		fmt.Printf("  %2d) %-12s%s\n", i+1, app.Name, extra)
	}
	fmt.Println("\nPresets:")
	for name, p := range presets {
		fmt.Printf("  %-12s %s\n", name+":", p.Description)
	}
}

func runInteractive() error {
	printAppsAndPresets()
	fmt.Print("\nPreset> ")
	var input string
	fmt.Scanln(&input)
	if input = strings.TrimSpace(input); input == "" {
		input = "dashboard"
	}
	if p, ok := presets[input]; ok {
		return runEnvironment(p.Docker, p.Terminal)
	}
	return fmt.Errorf("invalid preset: %s", input)
}
