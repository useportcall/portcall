package harness

import (
	"os"
	"path/filepath"
)

// FindRootDir walks up from cwd looking for go.work to find the monorepo root.
func FindRootDir() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.work")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return dir
		}
		dir = parent
	}
}
