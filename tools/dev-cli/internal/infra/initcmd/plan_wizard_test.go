package initcmd

import "testing"

func TestApplyInteractivePlanDefaultsPromptsForDomain(t *testing.T) {
	plan := applyInteractivePlanDefaults(Plan{Domain: "example.com"}, Options{}, Deps{
		IsInteractive: func() bool { return true },
		Warn:          func(string, ...any) {},
		AskText:       func(string, string) string { return "portcall.com" },
	})
	if plan.Domain != "portcall.com" {
		t.Fatalf("expected prompted domain update, got %q", plan.Domain)
	}
}

func TestApplyInteractivePlanDefaultsSkipsInNonInteractive(t *testing.T) {
	plan := applyInteractivePlanDefaults(Plan{Domain: "example.com"}, Options{}, Deps{
		IsInteractive: func() bool { return false },
	})
	if plan.Domain != "example.com" {
		t.Fatalf("unexpected domain change: %q", plan.Domain)
	}
}
