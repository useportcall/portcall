package initcmd

import (
	"fmt"
	"os"
	"strings"
)

func resolveDOTokenForInit(deps Deps) (string, error) {
	token, err := deps.ResolveDOToken()
	if err == nil && strings.TrimSpace(token) != "" {
		return token, nil
	}
	if !deps.IsInteractive() {
		return "", err
	}
	deps.Section("DigitalOcean Auth")
	deps.Warn("DigitalOcean auth not detected.")
	if deps.AskYesNo("Run doctl auth init now? [Y/n]: ", true) {
		if initErr := deps.RunCmd("doctl", "auth", "init", "--context", "portcall"); initErr != nil {
			deps.Warn("doctl auth init failed: %v", initErr)
		}
	}
	if deps.AskYesNo("Switch doctl context to portcall? [Y/n]: ", true) {
		if switchErr := deps.RunCmd("doctl", "auth", "switch", "--context", "portcall"); switchErr != nil {
			deps.Warn("doctl auth switch failed: %v", switchErr)
		}
	}
	token, err = deps.ResolveDOToken()
	if err == nil && strings.TrimSpace(token) != "" {
		return token, nil
	}
	entered := strings.TrimSpace(deps.AskText("Paste DigitalOcean token (or Enter to cancel): ", ""))
	if entered == "" {
		return "", fmt.Errorf("digitalocean auth required")
	}
	_ = os.Setenv("DIGITALOCEAN_TOKEN", entered)
	return entered, nil
}
