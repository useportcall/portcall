package infra

import "testing"

func TestInfraCommandIncludesCreateAndUpdate(t *testing.T) {
	cmd := newInfraCmd()
	names := map[string]bool{}
	for _, sub := range cmd.Commands() {
		names[sub.Name()] = true
	}
	for _, required := range []string{"create", "update", "init", "downscale", "status", "doctor", "pull", "cleanup"} {
		if !names[required] {
			t.Fatalf("missing infra subcommand %q", required)
		}
	}
}
