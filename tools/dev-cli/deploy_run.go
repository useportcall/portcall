package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func runDeploy(cmd *cobra.Command, _ []string) error {
	var err error
	rootDir, err = findRootDir()
	if err != nil {
		return err
	}
	if deployOpts.listApps {
		listDeployApps()
		return nil
	}
	if err = applyDeployDefaults(cmd); err != nil {
		return err
	}
	deployOpts.valuesFile = resolveValuesFile(deployOpts.cluster, deployOpts.valuesFile)
	if _, err = os.Stat(deployOpts.valuesFile); err != nil {
		return fmt.Errorf("values file not found: %s", deployOpts.valuesFile)
	}
	apps, err := selectApps(deployOpts.appInput)
	if err != nil {
		return err
	}
	if err = validateVersionType(); err != nil {
		return err
	}
	if err = resolvePreflightSelection(cmd); err != nil {
		return err
	}
	printSummary(apps)
	if !deployOpts.yes && strings.ToLower(readInput("Proceed with deployment? [y/N]: ")) != "y" {
		warn("Deployment cancelled")
		return nil
	}
	if deployOpts.dryRun {
		warn("DRY RUN - No changes will be made")
	}
	if !deployOpts.skipPreflight {
		if err := runPreflightSuites(); err != nil {
			return err
		}
	}
	return deployAppsToK8s(apps)
}

func validateVersionType() error {
	switch deployOpts.version {
	case "major", "minor", "patch", "skip":
		return nil
	default:
		return fmt.Errorf("invalid version type: %s", deployOpts.version)
	}
}

func printSummary(apps []string) {
	warn("====== DEPLOYMENT SUMMARY ======")
	plain("Apps:     %s", strings.Join(apps, ", "))
	plain("Cluster:  %s", deployOpts.cluster)
	plain("Values:   %s", deployOpts.valuesFile)
	plain("Version:  %s", deployOpts.version)
	plain("Build:    %v", !deployOpts.skipBuild)
	plain("Migration: %v", !deployOpts.skipMigration)
	plain("Preflight unit:        %v", deployOpts.runUnitTests && !deployOpts.skipPreflight)
	plain("Preflight integration: %v", deployOpts.runIntegTests && !deployOpts.skipPreflight)
	plain("Preflight e2e:         %v", deployOpts.runE2ETests && !deployOpts.skipPreflight)
	plain("Smoke tests:           %v", !deployOpts.skipSmokeTests)
	plain("Dry-run:  %v", deployOpts.dryRun)
	warn("================================")
}
