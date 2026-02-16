package initcmd

import "testing"

func TestResolveCloudflareDNSIntentFromFlags(t *testing.T) {
	run, auto := resolveCloudflareDNSIntent(Options{DNSProvider: "cloudflare", DNSAuto: true}, Deps{})
	if !run || !auto {
		t.Fatalf("expected run=true auto=true, got run=%t auto=%t", run, auto)
	}
}

func TestResolveCloudflareDNSIntentInteractivePrompt(t *testing.T) {
	answers := []bool{true, false}
	i := 0
	run, auto := resolveCloudflareDNSIntent(Options{}, Deps{
		IsInteractive: func() bool { return true },
		AskYesNo: func(string, bool) bool {
			v := answers[i]
			i++
			return v
		},
	})
	if !run || auto {
		t.Fatalf("expected run=true auto=false, got run=%t auto=%t", run, auto)
	}
}
