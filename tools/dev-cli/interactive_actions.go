package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func runSelfCommand(args ...string) error {
	fmt.Printf("\n$ %s %s\n\n", os.Args[0], joinArgs(args))
	cmd := exec.Command(os.Args[0], args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runInfraWizard(action string) error {
	cluster, err := promptWithDefault("Cluster alias", "digitalocean")
	if err != nil {
		if errors.Is(err, errBack) {
			return nil
		}
		return err
	}
	name, err := promptWithDefault("Cluster name", "portcall-micro")
	if err != nil {
		if errors.Is(err, errBack) {
			return nil
		}
		return err
	}
	domain, err := promptWithDefault("Domain", "example.com")
	if err != nil {
		if errors.Is(err, errBack) {
			return nil
		}
		return err
	}
	dnsProvider, err := promptWithDefault("DNS provider (manual|cloudflare)", "manual")
	if err != nil {
		if errors.Is(err, errBack) {
			return nil
		}
		return err
	}
	args := []string{"infra", action, "--cluster", cluster, "--name", name, "--domain", domain, "--dns-provider", dnsProvider}
	return runSelfCommand(args...)
}

func runDeployWizard() error {
	cluster, err := promptWithDefault("Cluster alias", "digitalocean")
	if err != nil {
		if errors.Is(err, errBack) {
			return nil
		}
		return err
	}
	apps, err := promptWithDefault("Apps (comma list or all)", "all")
	if err != nil {
		if errors.Is(err, errBack) {
			return nil
		}
		return err
	}
	version, err := promptWithDefault("Version bump (major|minor|patch|skip)", "patch")
	if err != nil {
		if errors.Is(err, errBack) {
			return nil
		}
		return err
	}
	return runSelfCommand("deploy", "--cluster", cluster, "--apps", apps, "--version", version)
}

func runEnvironmentWizard() error {
	preset, err := promptWithDefault("Run preset", "dashboard")
	if err != nil {
		if errors.Is(err, errBack) {
			return nil
		}
		return err
	}
	return runSelfCommand("run", "--preset", preset)
}

func joinArgs(args []string) string {
	out := ""
	for i, a := range args {
		if i > 0 {
			out += " "
		}
		out += a
	}
	return out
}
