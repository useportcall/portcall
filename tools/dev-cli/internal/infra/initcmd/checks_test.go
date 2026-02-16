package initcmd

import "testing"

func TestConfirmGuardrails_CreateWithYesOnExistingFails(t *testing.T) {
	deps := Deps{
		GetClusterState: func(string) (ClusterState, bool) { return ClusterState{}, true },
		RunCmdOut:       func(string, ...string) (string, error) { return "", nil },
		Section:         func(string) {},
		Warn:            func(string, ...any) {},
		Plain:           func(string, ...any) {},
		Info:            func(string, ...any) {},
		AskYesNo:        func(string, bool) bool { return false },
		AskText:         func(string, string) string { return "" },
	}
	opts := Options{Yes: true, Action: "create"}
	err := confirmGuardrails(Plan{ClusterName: "portcall-micro"}, opts, deps)
	if err == nil {
		t.Fatal("expected create guardrail failure with --yes")
	}
}

func TestConfirmGuardrails_UpdateWithYesOnExistingPasses(t *testing.T) {
	deps := Deps{
		GetClusterState: func(string) (ClusterState, bool) { return ClusterState{}, true },
		RunCmdOut:       func(string, ...string) (string, error) { return "", nil },
		Section:         func(string) {},
		Warn:            func(string, ...any) {},
		Plain:           func(string, ...any) {},
		Info:            func(string, ...any) {},
		AskYesNo:        func(string, bool) bool { return false },
		AskText:         func(string, string) string { return "" },
	}
	opts := Options{Yes: true, Action: "update"}
	if err := confirmGuardrails(Plan{ClusterName: "portcall-micro"}, opts, deps); err != nil {
		t.Fatalf("expected update guardrail pass, got %v", err)
	}
}
