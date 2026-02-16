package initcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type microValues struct {
	Registry     string
	Domain       string
	PostgresHost string
	PostgresPort string
	PostgresUser string
	PostgresDB   string
	KeycloakDB   string
	RedisHost    string
	RedisPort    string
	SpacesRegion   string
	SpacesEndpoint string
	SpacesBucket   string
	AllowedIPs   []string
}

func writeMicroValuesFile(alias string, in microValues, deps Deps) (string, error) {
	root := deps.RootDir()
	cfgDir := filepath.Join(root, ".infra", alias)
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		return "", fmt.Errorf("create infra config dir: %w", err)
	}
	tplDir := filepath.Join(root, "tools", "dev-cli", "internal", "infra", "templates")
	tpl, err := os.ReadFile(filepath.Join(tplDir, "values.micro.yaml.tmpl"))
	if err != nil {
		return "", fmt.Errorf("read values template: %w", err)
	}
	ingressTpl, err := os.ReadFile(filepath.Join(tplDir, "values.micro.ingress.yaml.tmpl"))
	if err != nil {
		return "", fmt.Errorf("read ingress template: %w", err)
	}
	body := strings.NewReplacer(
		"__REGISTRY__", in.Registry, "__DOMAIN__", in.Domain,
		"__PG_HOST__", in.PostgresHost, "__PG_PORT__", in.PostgresPort,
		"__PG_USER__", in.PostgresUser, "__PG_DB__", in.PostgresDB,
		"__KC_DB__", in.KeycloakDB, "__REDIS_HOST__", in.RedisHost,
		"__REDIS_PORT__", in.RedisPort, "__SPACES_REGION__", in.SpacesRegion,
		"__SPACES_ENDPOINT__", in.SpacesEndpoint, "__SPACES_BUCKET__", in.SpacesBucket,
		"__ALLOWED_IPS__", renderAllowedIPs(in.AllowedIPs),
	).Replace(string(tpl) + "\n" + string(ingressTpl))
	path := filepath.Join(cfgDir, "values.micro.yaml")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		return "", fmt.Errorf("write values file: %w", err)
	}
	return path, nil
}

func renderAllowedIPs(ips []string) string {
	if len(ips) == 0 {
		return " []"
	}
	items := make([]string, 0, len(ips))
	for _, ip := range ips {
		items = append(items, fmt.Sprintf("\n      - %q", ip))
	}
	return strings.Join(items, "")
}
