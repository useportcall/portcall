package infra

import (
	"testing"
)

func TestResolveClusterContext_Defaults(t *testing.T) {
	root := t.TempDir()
	if got := ResolveClusterContext(root, "digitalocean"); got != "do-k8s-portcall-prod" {
		t.Fatalf("context mismatch: %s", got)
	}
	if got := ResolveClusterContext(root, "staging"); got != "staging" {
		t.Fatalf("context mismatch: %s", got)
	}
}

func TestResolveClusterContext_FromState(t *testing.T) {
	root := t.TempDir()
	state := infraState{Clusters: map[string]ClusterState{"staging": {Context: "do-k8s-portcall-staging"}}}
	if err := saveInfraState(root, state); err != nil {
		t.Fatalf("save state: %v", err)
	}
	if got := ResolveClusterContext(root, "staging"); got != "do-k8s-portcall-staging" {
		t.Fatalf("context mismatch: %s", got)
	}
}

func TestResolveValuesAndImageFromState(t *testing.T) {
	root := t.TempDir()
	state := infraState{Clusters: map[string]ClusterState{"staging": {
		Values: "/tmp/staging-values.yaml", Registry: "registry.digitalocean.com/staging-reg",
	}}}
	if err := saveInfraState(root, state); err != nil {
		t.Fatalf("save state: %v", err)
	}
	if got := ResolveValuesFile(root, "staging", ""); got != "/tmp/staging-values.yaml" {
		t.Fatalf("values mismatch: %s", got)
	}
	if got := ResolveImageRepository(root, "staging", "registry.digitalocean.com/portcall-registry/portcall-api"); got != "registry.digitalocean.com/staging-reg/portcall-api" {
		t.Fatalf("image repo mismatch: %s", got)
	}
}
