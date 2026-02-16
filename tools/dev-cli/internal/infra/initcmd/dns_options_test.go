package initcmd

import "testing"

func TestNormalizeDNSOptionsDefaultsToManual(t *testing.T) {
	opts, err := normalizeDNSOptions(Options{})
	if err != nil {
		t.Fatalf("normalize options: %v", err)
	}
	if opts.DNSProvider != "manual" {
		t.Fatalf("expected manual default, got %q", opts.DNSProvider)
	}
}

func TestNormalizeDNSOptionsRejectsUnknownProvider(t *testing.T) {
	_, err := normalizeDNSOptions(Options{DNSProvider: "route53"})
	if err == nil {
		t.Fatal("expected provider validation error")
	}
}
