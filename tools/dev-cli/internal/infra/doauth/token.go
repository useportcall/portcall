package doauth

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type RunCmdOut func(name string, args ...string) (string, error)

func ResolveDigitalOceanToken(runCmdOut RunCmdOut) (string, error) {
	if token := strings.TrimSpace(os.Getenv("DIGITALOCEAN_TOKEN")); token != "" {
		return token, nil
	}
	current, err := currentDoctlContext(runCmdOut)
	if err != nil {
		return "", fmt.Errorf("digitalocean auth not configured. run:\n- doctl auth init --context portcall\n- doctl auth switch --context portcall\n- export DIGITALOCEAN_TOKEN=<token>")
	}
	token, err := TokenFromDoctlConfig(current)
	if err != nil || strings.TrimSpace(token) == "" {
		return "", fmt.Errorf("could not resolve digitalocean token for context %q. run:\n- doctl auth init --context %s\n- doctl auth switch --context %s\n- export DIGITALOCEAN_TOKEN=<token>", current, current, current)
	}
	return token, nil
}

func currentDoctlContext(runCmdOut RunCmdOut) (string, error) {
	out, err := runCmdOut("doctl", "auth", "list", "--output", "json")
	if err != nil {
		return "", err
	}
	var rows []struct {
		Name    string `json:"name"`
		Current bool   `json:"current"`
	}
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		return "", err
	}
	for _, row := range rows {
		if row.Current {
			return strings.TrimSpace(row.Name), nil
		}
	}
	return "", fmt.Errorf("no current doctl context")
}
