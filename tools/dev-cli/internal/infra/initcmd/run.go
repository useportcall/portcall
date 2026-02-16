package initcmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

func NewCommands(deps Deps) (init, create, update *cobra.Command) {
	var opts Options
	run := func(_ *cobra.Command, _ []string) error {
		return RunInit(opts, deps)
	}

	init = &cobra.Command{Use: "init", Short: "Compatibility alias for infra create/update", RunE: run}
	applyFlags(init, &opts)
	_ = init.Flags().MarkDeprecated("provider", "use --providers")

	create = &cobra.Command{Use: "create", Short: "Create a new infrastructure stack",
		RunE: func(c *cobra.Command, a []string) error {
			opts.Action = "create"
			return RunInit(opts, deps)
		}}
	applyFlags(create, &opts)
	_ = create.Flags().Set("action", "create")
	_ = create.Flags().MarkHidden("action")
	_ = create.Flags().MarkDeprecated("provider", "use --providers")

	update = &cobra.Command{Use: "update", Short: "Update an existing infrastructure stack",
		RunE: func(c *cobra.Command, a []string) error {
			opts.Action = "update"
			return RunInit(opts, deps)
		}}
	applyFlags(update, &opts)
	_ = update.Flags().Set("action", "update")
	_ = update.Flags().MarkHidden("action")
	_ = update.Flags().MarkDeprecated("provider", "use --providers")

	return init, create, update
}

func RunInit(opts Options, deps Deps) error {
	if err := deps.EnsureRootDir(); err != nil {
		return err
	}
	var err error
	opts, err = normalizeDNSOptions(opts)
	if err != nil {
		return err
	}
	deps.Section("Infra Init")
	plan, err := buildPlan(opts, deps)
	if err != nil {
		return err
	}
	plan = applyInteractivePlanDefaults(plan, opts, deps)
	steps, err := deps.ResolveInitSteps(plan, opts)
	if err != nil {
		return err
	}
	if err := confirmGuardrails(plan, opts, deps); err != nil {
		return err
	}
	if err := confirmDomainGuardrail(plan, opts, deps); err != nil {
		return err
	}
	if err := checkPrereqs(plan, opts, deps); err != nil {
		return err
	}
	if err := confirmInit(plan, opts, deps); err != nil {
		return err
	}
	if opts.DryRun {
		previewDryRun(plan, opts, steps, deps)
		deps.OK("Infra dry-run complete (no resources were created)")
		return nil
	}
	deps.Info("Preparing terraform workflow for alias %s", plan.Alias)
	if err := deps.RunTerraformSteps(plan, opts, steps); err != nil {
		return err
	}
	if err := finalizeInfra(plan, opts, deps); err != nil {
		return err
	}
	runDNSWorkflow(plan, opts, deps)
	deps.OK("Infra init complete")
	root := deps.RootDir()
	deps.Plain("Values file: %s", filepath.Join(root, ".infra", plan.Alias, "values.micro.yaml"))
	deps.Plain("Deploy next: go run ./tools/dev-cli deploy --cluster %s --apps all --version patch", plan.Alias)
	return nil
}
