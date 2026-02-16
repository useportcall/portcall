package inputs

import (
	"reflect"
	"testing"
)

func TestResolveProvidersDefault(t *testing.T) {
	got, err := ResolveProviders(nil, "")
	if err != nil {
		t.Fatalf("ResolveProviders error: %v", err)
	}
	if !reflect.DeepEqual(got, []string{"digitalocean"}) {
		t.Fatalf("providers mismatch: %v", got)
	}
}

func TestResolveProvidersFallback(t *testing.T) {
	got, err := ResolveProviders(nil, "digitalocean")
	if err != nil {
		t.Fatalf("ResolveProviders error: %v", err)
	}
	if !reflect.DeepEqual(got, []string{"digitalocean"}) {
		t.Fatalf("providers mismatch: %v", got)
	}
}

func TestResolveProvidersRejectUnsupported(t *testing.T) {
	if _, err := ResolveProviders([]string{"aws"}, ""); err == nil {
		t.Fatal("expected unsupported provider error")
	}
}
