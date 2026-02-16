package initcmd

import "testing"

func TestConfirmDomainGuardrailRejectsPlaceholderWithYes(t *testing.T) {
	deps := Deps{
		Section:  func(string) {},
		Warn:     func(string, ...any) {},
		Plain:    func(string, ...any) {},
		Info:     func(string, ...any) {},
		AskText:  func(string, string) string { return "" },
		AskYesNo: func(string, bool) bool { return false },
	}
	err := confirmDomainGuardrail(Plan{Domain: "example.com"}, Options{Yes: true}, deps)
	if err == nil {
		t.Fatal("expected placeholder domain to fail with --yes")
	}
}

func TestConfirmDomainGuardrailAcceptsInteractiveOverride(t *testing.T) {
	deps := Deps{
		Section:  func(string) {},
		Warn:     func(string, ...any) {},
		Plain:    func(string, ...any) {},
		Info:     func(string, ...any) {},
		AskText:  func(string, string) string { return "placeholder domain" },
		AskYesNo: func(string, bool) bool { return true },
	}
	if err := confirmDomainGuardrail(Plan{Domain: "example.com"}, Options{}, deps); err != nil {
		t.Fatalf("expected interactive override to pass, got %v", err)
	}
}

func TestConfirmDomainGuardrailSkipsForRealDomain(t *testing.T) {
	if err := confirmDomainGuardrail(Plan{Domain: "portcall.com"}, Options{Yes: true}, Deps{}); err != nil {
		t.Fatalf("real domain should pass: %v", err)
	}
}
