package keycloak

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type client struct {
	pfCmd *exec.Cmd
	token string
}

func newClient(namespace string) (*client, error) {
	fmt.Println("Getting Keycloak admin credentials...")
	user, err := secret(namespace, "KC_BOOTSTRAP_ADMIN_USERNAME")
	if err != nil {
		return nil, fmt.Errorf("read admin user: %w", err)
	}
	pass, err := secret(namespace, "KC_BOOTSTRAP_ADMIN_PASSWORD")
	if err != nil {
		return nil, fmt.Errorf("read admin pass: %w", err)
	}
	pf := exec.Command("kubectl", "port-forward", "-n", namespace, "deployment/keycloak", "18080:8080")
	if err := pf.Start(); err != nil {
		return nil, fmt.Errorf("port-forward: %w", err)
	}
	time.Sleep(3 * time.Second)
	token, err := loginToken(user, pass)
	if err != nil {
		_ = pf.Process.Kill()
		return nil, err
	}
	return &client{pfCmd: pf, token: token}, nil
}

func (kc *client) close() {
	if kc.pfCmd != nil && kc.pfCmd.Process != nil {
		_ = kc.pfCmd.Process.Kill()
	}
}

func (kc *client) updateSMTP(realm string) error {
	payload := `{"smtpServer":{"host":"email-worker.portcall.svc.cluster.local","port":"2525","from":"notifications@mail.useportcall.com","fromDisplayName":"Portcall","auth":"false","starttls":"false"}}`
	out, err := run("curl", "-s", "-w", "\n%{http_code}", "-X", "PUT",
		fmt.Sprintf("http://localhost:18080/admin/realms/%s", realm),
		"-H", "Authorization: Bearer "+kc.token, "-H", "Content-Type: application/json", "-d", payload)
	if err != nil {
		return fmt.Errorf("update realm: %w", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	code := lines[len(lines)-1]
	if code != "204" && code != "200" {
		return fmt.Errorf("update realm returned HTTP %s", code)
	}
	return kc.showSMTPStatus(realm)
}

func (kc *client) showSMTPStatus(realm string) error {
	out, err := run("curl", "-s", fmt.Sprintf("http://localhost:18080/admin/realms/%s", realm),
		"-H", "Authorization: Bearer "+kc.token)
	if err != nil {
		return fmt.Errorf("get realm: %w", err)
	}
	var data map[string]any
	if err := json.Unmarshal([]byte(out), &data); err != nil {
		return fmt.Errorf("parse realm json: %w", err)
	}
	smtp, _ := data["smtpServer"].(map[string]any)
	resetAllowed, _ := data["resetPasswordAllowed"].(bool)
	fmt.Println("\nSMTP Configuration")
	for _, k := range []string{"host", "port", "from", "fromDisplayName", "auth", "starttls"} {
		fmt.Printf("  %-20s %v\n", k+":", smtp[k])
	}
	fmt.Printf("  %-20s %v\n\n", "resetPasswordAllowed:", resetAllowed)
	return nil
}
