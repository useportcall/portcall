package initcmd

func applyInteractivePlanDefaults(plan Plan, opts Options, deps Deps) Plan {
	if opts.Yes || opts.DryRun || !deps.IsInteractive() {
		return plan
	}
	if isPlaceholderDomain(plan.Domain) {
		deps.Warn("Using placeholder domain %q.", plan.Domain)
		plan.Domain = deps.AskText("Enter base domain for ingress hosts (Enter to keep placeholder): ", plan.Domain)
	}
	return plan
}
