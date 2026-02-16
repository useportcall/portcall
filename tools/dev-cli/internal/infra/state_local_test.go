package infra

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInfraStateRoundTrip(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	state := infraState{Clusters: map[string]ClusterState{"digitalocean": {
		Context: "do-k8s-portcall-micro", Registry: "registry.digitalocean.com/portcall-registry",
		Values: ".infra/digitalocean/values.micro.yaml", Mode: "micro", Cluster: "portcall-micro", Namespace: "portcall",
	}}}
	if err := saveInfraState(dir, state); err != nil {
		t.Fatalf("saveInfraState error: %v", err)
	}
	got, err := loadInfraState(dir)
	if err != nil {
		t.Fatalf("loadInfraState error: %v", err)
	}
	cfg := got.Clusters["digitalocean"]
	if cfg.Context == "" || cfg.Registry == "" || cfg.Values == "" {
		t.Fatalf("unexpected state: %+v", cfg)
	}
}

func TestLoadInfraStateMissingFile(t *testing.T) {
	t.Parallel()
	state, err := loadInfraState(t.TempDir())
	if err != nil {
		t.Fatalf("loadInfraState error: %v", err)
	}
	if len(state.Clusters) != 0 {
		t.Fatalf("expected empty cluster map, got %d", len(state.Clusters))
	}
}

func TestLoadInfraStateInvalidJSON(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, ".dev-cli.infra.json")
	if err := os.WriteFile(path, []byte("{"), 0o644); err != nil {
		t.Fatalf("write invalid state: %v", err)
	}
	if _, err := loadInfraState(dir); err == nil {
		t.Fatal("expected parse error")
	}
}
