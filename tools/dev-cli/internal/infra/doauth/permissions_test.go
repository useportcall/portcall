package doauth

import (
	"errors"
	"strings"
	"testing"
)

func TestVerifyDigitalOceanAccessPassesWithExpectedProbeErrors(t *testing.T) {
	run := func(env map[string]string, name string, args ...string) (string, error) {
		if name != "doctl" {
			t.Fatalf("unexpected command: %s", name)
		}
		if env["DIGITALOCEAN_TOKEN"] != "dop_v1_test" {
			t.Fatalf("missing token env")
		}
		key := strings.Join(args, " ")
		if strings.Contains(key, "options") || strings.Contains(key, "account get") {
			return "ok", nil
		}
		return "cluster not found", errors.New("failed")
	}
	if err := VerifyDigitalOceanAccess("dop_v1_test", run); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestVerifyDigitalOceanAccessFailsOnPermissionDenied(t *testing.T) {
	run := func(env map[string]string, name string, args ...string) (string, error) {
		key := strings.Join(args, " ")
		if strings.Contains(key, "options") || strings.Contains(key, "account get") {
			return "ok", nil
		}
		return "forbidden: token is read-only", errors.New("failed")
	}
	err := VerifyDigitalOceanAccess("dop_v1_test", run)
	if err == nil || !strings.Contains(err.Error(), "missing write permissions") {
		t.Fatalf("expected write permission error, got %v", err)
	}
}

func TestVerifyDigitalOceanAccessFailsOnReadCheck(t *testing.T) {
	run := func(env map[string]string, name string, args ...string) (string, error) {
		if strings.Join(args, " ") == "account get --output json" {
			return "unauthorized", errors.New("failed")
		}
		return "ok", nil
	}
	err := VerifyDigitalOceanAccess("dop_v1_test", run)
	if err == nil || !strings.Contains(err.Error(), "doctl account check failed") {
		t.Fatalf("expected account read failure, got %v", err)
	}
}
