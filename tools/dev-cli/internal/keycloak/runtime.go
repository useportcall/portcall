package keycloak

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

func run(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s %s failed: %w\n%s", name, strings.Join(args, " "), err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

func secret(ns, key string) (string, error) {
	var wrapped struct {
		Data map[string]string `json:"data"`
	}
	out, err := run("kubectl", "get", "secret", "portcall-secrets", "-n", ns, "-o", "json")
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal([]byte(out), &wrapped); err != nil {
		return "", err
	}
	val, ok := wrapped.Data[key]
	if !ok || strings.TrimSpace(val) == "" {
		return "", fmt.Errorf("missing secret key %s", key)
	}
	decoded, err := run("bash", "-c", fmt.Sprintf("echo '%s' | base64 -d", val))
	if err != nil {
		return "", err
	}
	return decoded, nil
}

func loginToken(user, pass string) (string, error) {
	out, err := run("curl", "-s", "-X", "POST",
		"http://localhost:18080/realms/master/protocol/openid-connect/token",
		"-H", "Content-Type: application/x-www-form-urlencoded",
		"-d", "username="+user, "-d", "password="+pass,
		"-d", "grant_type=password", "-d", "client_id=admin-cli")
	if err != nil {
		return "", err
	}
	var tok struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal([]byte(out), &tok); err != nil {
		return "", fmt.Errorf("parse token: %w", err)
	}
	if tok.AccessToken == "" {
		return "", fmt.Errorf("empty access token")
	}
	return tok.AccessToken, nil
}
