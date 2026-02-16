package keycloak

import "testing"

func TestNewCommand_HasSubcommands(t *testing.T) {
	cmd := NewCommand(func(string) error { return nil })
	if cmd == nil {
		t.Fatal("expected command")
	}
	names := map[string]bool{}
	for _, c := range cmd.Commands() {
		names[c.Name()] = true
	}
	if !names["smtp-update"] || !names["smtp-status"] {
		t.Fatal("expected smtp-update and smtp-status subcommands")
	}
}
