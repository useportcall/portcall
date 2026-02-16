package tfflow

import (
	"reflect"
	"testing"
)

func TestTerraformSteps_DefaultAll(t *testing.T) {
	opts := Options{SkipPostgres: false, SkipRedis: false, SkipSpaces: false}
	got := TerraformSteps(Plan{Step: "all"}, opts)
	want := []string{"cluster", "postgres", "redis", "object-storage"}
	if len(got) != len(want) {
		t.Fatalf("unexpected step count: got %d want %d", len(got), len(want))
	}
	for i := range want {
		if got[i].Name != want[i] {
			t.Fatalf("step mismatch at %d: got %s want %s", i, got[i].Name, want[i])
		}
	}
}

func TestTerraformSteps_ServiceSkips(t *testing.T) {
	opts := Options{SkipPostgres: true, SkipRedis: false, SkipSpaces: true}
	got := TerraformSteps(Plan{Step: "services"}, opts)
	want := []Step{{Name: "redis", Targets: RedisTargets}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("steps mismatch: got %v want %v", got, want)
	}
}
