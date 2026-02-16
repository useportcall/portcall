package doauth

import "fmt"

type RunCmdOutWithEnv func(env map[string]string, name string, args ...string) (string, error)

type doctlCheck struct {
	Name string
	Args []string
}

func VerifyDigitalOceanAccess(token string, runCmdOutWithEnv RunCmdOutWithEnv) error {
	env := map[string]string{"DIGITALOCEAN_TOKEN": token}
	if err := runReadChecks(env, runCmdOutWithEnv); err != nil {
		return err
	}
	return runWriteProbes(env, runCmdOutWithEnv)
}

func runReadChecks(env map[string]string, runCmdOutWithEnv RunCmdOutWithEnv) error {
	for _, check := range doctlReadChecks() {
		out, err := runCmdOutWithEnv(env, "doctl", check.Args...)
		if err != nil {
			return fmt.Errorf("%s check failed: %s", check.Name, formatProbeError(out))
		}
	}
	return nil
}

func runWriteProbes(env map[string]string, runCmdOutWithEnv RunCmdOutWithEnv) error {
	for _, probe := range doctlWriteProbes() {
		out, err := runCmdOutWithEnv(env, "doctl", probe.Args...)
		if err == nil {
			return fmt.Errorf("%s probe unexpectedly succeeded; refusing to continue", probe.Name)
		}
		switch {
		case hasPermissionDenied(out):
			return fmt.Errorf("%s probe shows missing write permissions: %s", probe.Name, formatProbeError(out))
		case hasExpectedProbeFailure(out):
			continue
		default:
			return fmt.Errorf("%s probe returned an unknown error: %s", probe.Name, formatProbeError(out))
		}
	}
	return nil
}
