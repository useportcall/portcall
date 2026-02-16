package main

func runPortcallInteractive() error {
	for {
		choice, err := selectMainAction()
		if err != nil {
			return err
		}
		switch choice {
		case "infra":
			if err := runInfraMenu(); err != nil {
				return err
			}
		case "deploy":
			if err := runDeployWizard(); err != nil {
				return err
			}
		case "run":
			if err := runEnvironmentWizard(); err != nil {
				return err
			}
		case "setup":
			if err := runSelfCommand("setup"); err != nil {
				return err
			}
		case "exit":
			return nil
		}
	}
}
