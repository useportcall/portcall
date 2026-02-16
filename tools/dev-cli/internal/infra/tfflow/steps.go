package tfflow

func TerraformStepCatalog() []Step {
	return []Step{{Name: "cluster", Targets: CoreTargets}, {Name: "postgres", Targets: PostgresTargets}, {Name: "redis", Targets: RedisTargets}, {Name: "object-storage", Targets: SpacesTargets}}
}

func TerraformSteps(plan Plan, opts Options) []Step {
	catalog := TerraformStepCatalog()
	steps := []Step{}
	if plan.Step == "all" || plan.Step == "core" {
		steps = append(steps, catalog[0])
	}
	if plan.Step == "all" || plan.Step == "services" {
		if !opts.SkipPostgres {
			steps = append(steps, catalog[1])
		}
		if !opts.SkipRedis {
			steps = append(steps, catalog[2])
		}
		if !opts.SkipSpaces {
			steps = append(steps, catalog[3])
		}
	}
	return steps
}

func ResolveInitSteps(plan Plan, opts Options, deps Deps) ([]Step, error) {
	defaults := TerraformSteps(plan, opts)
	deps.Section("Init Wizard")
	deps.Plain("Action      : %s", plan.Action)
	deps.Plain("Alias       : %s", plan.Alias)
	deps.Plain("Cluster     : %s", plan.ClusterName)
	deps.Plain("Provider    : %s", plan.Provider)
	deps.Plain("Default mode: micro/%s", plan.Step)
	if opts.Yes {
		deps.Info("Non-interactive mode enabled by --yes; using default stage selection")
		return defaults, nil
	}
	if !deps.IsInteractive() {
		deps.Warn("No interactive terminal detected; using default stage selection")
		return defaults, nil
	}
	deps.Info("Stages available (default selected based on --step and skip flags):")
	catalog := TerraformStepCatalog()
	for _, stage := range catalog {
		state := "off"
		if isStageSelected(defaults, stage.Name) {
			state = "on"
		}
		deps.Plain("- %s [default:%s]: %s", stage.Name, state, stepDescription(stage.Name))
	}
	if deps.AskYesNo("Use default stage selection? [Y/n]: ", true) {
		return defaults, nil
	}
	selected := []Step{}
	for _, stage := range catalog {
		if deps.AskYesNo("Include stage \""+stage.Name+"\"? [y/N]: ", isStageSelected(defaults, stage.Name)) {
			selected = append(selected, stage)
		}
	}
	return selected, nil
}

func isStageSelected(steps []Step, name string) bool {
	for _, step := range steps {
		if step.Name == name {
			return true
		}
	}
	return false
}

func stepDescription(name string) string {
	desc := map[string]string{"cluster": "create/update VPC, Kubernetes cluster, and container registry", "postgres": "create/update managed Postgres and DB firewall", "redis": "create/update managed Redis and DB firewall", "object-storage": "create/update Spaces buckets and access key"}
	return desc[name]
}
