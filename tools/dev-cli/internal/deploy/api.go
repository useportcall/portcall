package deploy

import "github.com/spf13/cobra"

type rootFinder func() (string, error)

var (
	findRoot rootFinder
	rootDir  string
)

func NewCommand(findRootDir rootFinder) *cobra.Command {
	findRoot = findRootDir
	return newDeployCmd()
}

func ensureRootDir() error {
	if rootDir != "" {
		return nil
	}
	dir, err := findRoot()
	if err != nil {
		return err
	}
	rootDir = dir
	return nil
}
