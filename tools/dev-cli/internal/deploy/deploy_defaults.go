package deploy

import "github.com/spf13/cobra"

func applyDeployDefaults(cmd *cobra.Command) error {
	if deployOpts.cluster == "" {
		deployOpts.cluster = "digitalocean"
	}
	if deployOpts.version == "" {
		deployOpts.version = "patch"
	}
	settings, err := loadDeploySettings(rootDir)
	if err != nil {
		return err
	}
	if cmd.Flags().Changed("skip-migration") {
		settings.SkipMigration = deployOpts.skipMigration
		if err := saveDeploySettings(rootDir, settings); err != nil {
			return err
		}
	} else {
		deployOpts.skipMigration = settings.SkipMigration
	}
	return nil
}

func resolvePreflightSelection(cmd *cobra.Command) error {
	if deployOpts.skipPreflight || deployOpts.yes {
		return nil
	}
	if !cmd.Flags().Changed("unit-tests") {
		deployOpts.runUnitTests = askYesNo("Run unit tests before deploy? [Y/n]: ", true)
	}
	if !cmd.Flags().Changed("integration-tests") {
		deployOpts.runIntegTests = askYesNo("Run integration tests before deploy? [Y/n]: ", true)
	}
	if !cmd.Flags().Changed("e2e-tests") {
		deployOpts.runE2ETests = askYesNo("Run e2e tests before deploy? [y/N]: ", false)
	}
	if !deployOpts.runUnitTests && !deployOpts.runIntegTests && !deployOpts.runE2ETests {
		warn("No preflight suites selected")
	}
	return nil
}

func askYesNo(prompt string, defaultYes bool) bool {
	for {
		in := readInput(prompt)
		if in == "" {
			return defaultYes
		}
		switch in[0] {
		case 'y', 'Y':
			return true
		case 'n', 'N':
			return false
		default:
			warn("Please answer y or n")
		}
	}
}
