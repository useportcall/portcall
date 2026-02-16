package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/useportcall/portcall/tools/dev-cli/internal/infra"
)

func switchCluster(cluster string) error {
	if err := ensureRootDir(); err != nil {
		return err
	}
	ctx := infra.ResolveClusterContext(rootDir, cluster)
	cmd := exec.Command("kubectl", "config", "use-context", ctx)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = rootDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("switch context %s: %w", ctx, err)
	}
	return nil
}
