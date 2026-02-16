package initcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestCheckPrereqsCallsDigitalOceanVerifier(t *testing.T) {
	restore := setFakePath(t, "terraform", "doctl", "kubectl", "helm")
	defer restore()
	called := false
	deps := Deps{
		ResolveDOToken: func() (string, error) { return "dop_v1_test", nil },
		VerifyDOAccess: func(token string) error {
			called = token == "dop_v1_test"
			return nil
		},
	}
	plan := Plan{Providers: []string{"digitalocean"}}
	if err := checkPrereqs(plan, Options{}, deps); err != nil {
		t.Fatalf("expected prereqs to pass, got %v", err)
	}
	if !called {
		t.Fatal("expected VerifyDOAccess to be called")
	}
}

func TestCheckPrereqsReturnsHelpfulDoError(t *testing.T) {
	restore := setFakePath(t, "terraform", "doctl", "kubectl", "helm")
	defer restore()
	deps := Deps{
		ResolveDOToken: func() (string, error) { return "dop_v1_test", nil },
		VerifyDOAccess: func(string) error { return fmt.Errorf("missing write permissions") },
	}
	err := checkPrereqs(Plan{Providers: []string{"digitalocean"}}, Options{}, deps)
	if err == nil || !containsLineWith([]string{err.Error()}, "doctl auth init") {
		t.Fatalf("expected remediation in error, got %v", err)
	}
}

func setFakePath(t *testing.T, names ...string) func() {
	t.Helper()
	dir := t.TempDir()
	for _, name := range names {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
			t.Fatalf("write fake binary %s: %v", name, err)
		}
	}
	prev := os.Getenv("PATH")
	if err := os.Setenv("PATH", dir); err != nil {
		t.Fatalf("set PATH: %v", err)
	}
	return func() { _ = os.Setenv("PATH", prev) }
}
