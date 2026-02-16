package tfflow

import "fmt"

func RunTerraformSteps(plan Plan, opts Options, selected []Step, deps Deps) error {
	if err := deps.RunCmd("terraform", "-chdir="+plan.StackDir, "init"); err != nil {
		return fmt.Errorf("terraform init failed: %w", err)
	}
	steps := selected
	if len(steps) == 0 {
		steps = TerraformSteps(plan, opts)
	}
	if len(steps) == 0 {
		deps.Warn("No terraform steps selected; skipping")
		return nil
	}
	for idx, step := range steps {
		deps.Info("[%d/%d] Terraform step: %s", idx+1, len(steps), step.Name)
		if err := runWithTargets(plan, opts, step.Name, step.Targets, deps); err != nil {
			return err
		}
	}
	return nil
}

func runWithTargets(plan Plan, opts Options, step string, targets []string, deps Deps) error {
	action := "apply"
	if opts.DryRun {
		action = "plan"
	}
	providerEnv := map[string]string{}
	if token, err := deps.ResolveDOToken(); err == nil && token != "" {
		providerEnv["DIGITALOCEAN_TOKEN"] = token
	}
	args := []string{"-chdir=" + plan.StackDir, action, "-var", "cluster_name=" + plan.ClusterName, "-var", "region=" + plan.Region, "-var", "node_size=" + plan.NodeSize, "-var", fmt.Sprintf("node_count=%d", plan.NodeCount), "-var", "redis_size=" + plan.RedisSize, "-var", "registry_name=" + plan.RegistryName, "-var", "vpc_ip_range=" + plan.VPCCIDR, "-var", "spaces_region=" + plan.SpacesRegion, "-var", "spaces_prefix=" + plan.SpacesPrefix}
	if opts.Yes && action == "apply" {
		args = append(args, "-auto-approve")
	}
	for _, target := range targets {
		args = append(args, "-target="+target)
	}
	if err := deps.RunCmdWithEnv(providerEnv, "terraform", args...); err != nil {
		return fmt.Errorf("terraform %s failed at step %q: %w", action, step, err)
	}
	return nil
}
