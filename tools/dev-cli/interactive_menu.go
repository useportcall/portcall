package main

func selectMainAction() (string, error) {
	return promptChoice("Portcall", []choiceItem{
		{Key: "1", Label: "Infra workflow", Value: "infra"},
		{Key: "2", Label: "Deploy apps", Value: "deploy"},
		{Key: "3", Label: "Run local env", Value: "run"},
		{Key: "4", Label: "Setup workstation", Value: "setup"},
		{Key: "q", Label: "Exit", Value: "exit"},
	})
}

func runInfraMenu() error {
	for {
		choice, err := promptChoice("Infra", []choiceItem{
			{Key: "1", Label: "Create infrastructure", Value: "create"},
			{Key: "2", Label: "Update infrastructure", Value: "update"},
			{Key: "3", Label: "Doctor checks", Value: "doctor"},
			{Key: "4", Label: "Status", Value: "status"},
			{Key: "5", Label: "Pull state", Value: "pull"},
			{Key: "b", Label: "Back", Value: "back"},
			{Key: "q", Label: "Exit", Value: "exit"},
		})
		if err != nil {
			return err
		}
		switch choice {
		case "create", "update":
			if err := runInfraWizard(choice); err != nil {
				return err
			}
		case "doctor":
			if err := runSelfCommand("infra", "doctor", "--cluster", "digitalocean"); err != nil {
				return err
			}
		case "status":
			if err := runSelfCommand("infra", "status", "--cluster", "digitalocean"); err != nil {
				return err
			}
		case "pull":
			if err := runSelfCommand("infra", "pull", "--cluster", "digitalocean"); err != nil {
				return err
			}
		case "back":
			return nil
		case "exit":
			return nil
		}
	}
}
