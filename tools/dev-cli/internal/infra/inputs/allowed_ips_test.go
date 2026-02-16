package inputs

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestReadAdminAllowedIPs(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "values.yaml")
	content := "admin:\n  ingress:\n    allowedIPs:\n      - \"203.0.113.5/32\"\n      - 198.51.100.0/24\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write values: %v", err)
	}
	got, err := ReadAdminAllowedIPs(path)
	if err != nil {
		t.Fatalf("ReadAdminAllowedIPs error: %v", err)
	}
	want := []string{"203.0.113.5/32", "198.51.100.0/24"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("allowed IP mismatch: got %v want %v", got, want)
	}
}

func TestUniqueIPs(t *testing.T) {
	t.Parallel()
	got := UniqueIPs([]string{"203.0.113.5/32", "", "203.0.113.5/32", "198.51.100.0/24"})
	want := []string{"203.0.113.5/32", "198.51.100.0/24"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("UniqueIPs mismatch: got %v want %v", got, want)
	}
}
