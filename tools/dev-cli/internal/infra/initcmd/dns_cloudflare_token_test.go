package initcmd

import (
	"os"
	"testing"

	"github.com/useportcall/portcall/tools/dev-cli/internal/infra/authstore"
)

func TestResolveCloudflareTokenPrefersEnv(t *testing.T) {
	prev := os.Getenv("CLOUDFLARE_API_TOKEN")
	_ = os.Setenv("CLOUDFLARE_API_TOKEN", "env-token")
	defer func() { _ = os.Setenv("CLOUDFLARE_API_TOKEN", prev) }()
	token := resolveCloudflareToken(Deps{RootDir: func() string { return t.TempDir() }})
	if token != "env-token" {
		t.Fatalf("expected env token, got %q", token)
	}
}

func TestReadCloudflareTokenFromStore(t *testing.T) {
	root := t.TempDir()
	if err := authstore.SaveCloudflareToken(root, "stored-token"); err != nil {
		t.Fatalf("save token: %v", err)
	}
	token := readCloudflareTokenFromStore(Deps{RootDir: func() string { return root }})
	if token != "stored-token" {
		t.Fatalf("unexpected token: %q", token)
	}
}
