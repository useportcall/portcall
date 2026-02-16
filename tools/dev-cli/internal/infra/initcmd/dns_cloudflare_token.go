package initcmd

import (
	"os"
	"strings"

	"github.com/useportcall/portcall/tools/dev-cli/internal/infra/authstore"
)

func resolveCloudflareToken(deps Deps) string {
	if token := strings.TrimSpace(os.Getenv("CLOUDFLARE_API_TOKEN")); token != "" {
		return token
	}
	if token := readCloudflareTokenFromStore(deps); token != "" {
		return token
	}
	if !deps.IsInteractive() {
		return ""
	}
	deps.Warn("Cloudflare API token not found in env or local auth store.")
	if deps.AskYesNo("Run wrangler login now (browser auth)? [Y/n]: ", true) {
		if err := deps.RunCmd("wrangler", "login"); err != nil {
			deps.Warn("wrangler login failed: %v", err)
		}
	}
	token := strings.TrimSpace(deps.AskText("Paste Cloudflare API token (or Enter to skip): ", ""))
	if token == "" {
		return ""
	}
	if deps.AskYesNo("Save token in .dev-cli.auth.json for this repo? [Y/n]: ", true) {
		if err := authstore.SaveCloudflareToken(deps.RootDir(), token); err != nil {
			deps.Warn("failed to save auth token: %v", err)
		}
	}
	return token
}

func readCloudflareTokenFromStore(deps Deps) string {
	state, err := authstore.Load(deps.RootDir())
	if err != nil {
		return ""
	}
	return strings.TrimSpace(state.CloudflareAPIToken)
}
