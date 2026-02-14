package main

import "fmt"

func deployAppsToK8s(apps []string) error {
	info("Switching to cluster: %s", deployOpts.cluster)
	if !deployOpts.dryRun {
		if err := switchCluster(deployOpts.cluster); err != nil {
			return err
		}
		if err := runCmd("kubectl", "get", "nodes", "--request-timeout=5s"); err != nil {
			return fmt.Errorf("failed to connect to cluster: %w", err)
		}
	}
	if !deployOpts.skipMigration {
		if err := runLocalMigration(); err != nil {
			return err
		}
	} else {
		warn("Skipping database migration")
	}
	versions := collectVersions()
	for _, name := range apps {
		if err := processOneApp(name, versions); err != nil {
			return err
		}
	}
	if !deployOpts.skipSmokeTests && !deployOpts.dryRun {
		info("Running smoke tests...")
		smoke(deployOpts.cluster, apps)
	}
	ok("Deployment complete")
	return nil
}
