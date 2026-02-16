package infra

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNoDuplicatePackageDeclarations(t *testing.T) {
	root := filepath.Join(".")
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		count := 0
		for _, line := range strings.Split(string(data), "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "package ") {
				count++
			}
		}
		if count > 1 {
			t.Fatalf("duplicate package declarations in %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk infra package: %v", err)
	}
}
