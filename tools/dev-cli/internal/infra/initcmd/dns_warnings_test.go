package initcmd

import (
	"fmt"
	"strings"
	"testing"
)

func TestPrintDNSChecklistWithDomainAndTarget(t *testing.T) {
	infoLines := []string{}
	plainLines := []string{}
	warnLines := []string{}
	deps := Deps{
		Section: func(string) {},
		Info:    func(msg string, args ...any) { infoLines = append(infoLines, fmt.Sprintf(msg, args...)) },
		Plain:   func(msg string, args ...any) { plainLines = append(plainLines, fmt.Sprintf(msg, args...)) },
		Warn:    func(msg string, args ...any) { warnLines = append(warnLines, fmt.Sprintf(msg, args...)) },
		RunCmdOut: func(string, ...string) (string, error) {
			return "203.0.113.10", nil
		},
	}
	printDNSChecklist(Plan{Domain: "portcall.test"}, deps)
	if len(infoLines) == 0 || !strings.Contains(infoLines[0], "203.0.113.10") {
		t.Fatalf("expected ingress target info line, got %v", infoLines)
	}
	if !containsLineWith(plainLines, "api.portcall.test") {
		t.Fatalf("expected api host in checklist, got %v", plainLines)
	}
	if len(warnLines) == 0 {
		t.Fatalf("expected at least one warning")
	}
}

func TestPrintDNSChecklistPlaceholderWarnsEarly(t *testing.T) {
	calls := 0
	warnLines := []string{}
	deps := Deps{
		Section: func(string) {},
		Warn:    func(msg string, args ...any) { warnLines = append(warnLines, fmt.Sprintf(msg, args...)) },
		Plain:   func(string, ...any) {},
		RunCmdOut: func(string, ...string) (string, error) {
			calls++
			return "", nil
		},
	}
	printDNSChecklist(Plan{Domain: "example.com"}, deps)
	if calls != 0 {
		t.Fatalf("did not expect kubectl lookup for placeholder domain")
	}
	if len(warnLines) < 2 {
		t.Fatalf("expected placeholder warnings, got %v", warnLines)
	}
}

func containsLineWith(lines []string, needle string) bool {
	for _, line := range lines {
		if strings.Contains(line, needle) {
			return true
		}
	}
	return false
}
