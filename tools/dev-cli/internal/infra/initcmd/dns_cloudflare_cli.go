package initcmd

import (
	"os/exec"
	"strings"
)

func reportCloudflareCLIStatus(deps Deps) {
	if _, err := exec.LookPath("wrangler"); err != nil {
		deps.Warn("wrangler CLI not found in PATH (Cloudflare CLI check skipped).")
		return
	}
	out, err := deps.RunCmdOut("wrangler", "whoami")
	if err != nil {
		deps.Warn("wrangler auth check failed: %s", compactText(out))
		deps.Plain("Run: wrangler login")
		return
	}
	deps.OK("wrangler is authenticated: %s", compactText(out))
}

func compactText(in string) string {
	trimmed := strings.Join(strings.Fields(strings.TrimSpace(in)), " ")
	if trimmed == "" {
		return "(no output)"
	}
	if len(trimmed) <= 180 {
		return trimmed
	}
	return trimmed[:177] + "..."
}
