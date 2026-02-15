//go:build e2e_live_email

package smtprelay

import (
	"net"
	netsmtp "net/smtp"
	"os"
	"testing"
	"time"

	gsmtp "github.com/emersion/go-smtp"
)

func TestSMTPRelayLiveResend(t *testing.T) {
	apiKey := os.Getenv("RESEND_API_KEY")
	from := os.Getenv("E2E_EMAIL_FROM")
	to := os.Getenv("E2E_EMAIL_TO")
	if apiKey == "" || from == "" || to == "" {
		t.Fatal("set RESEND_API_KEY, E2E_EMAIL_FROM, and E2E_EMAIL_TO")
	}

	be := &Backend{cfg: Config{EmailProvider: "resend", ResendToken: apiKey, ResendAPI: "https://api.resend.com/emails"}}
	s := gsmtp.NewServer(be)
	s.Domain = "localhost"
	s.AllowInsecureAuth = true
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	go func() { _ = s.Serve(ln) }()

	subject := "Live E2E Password Reset " + time.Now().UTC().Format("20060102-150405")
	msg := []byte("From: " + from + "\r\nTo: " + to + "\r\nSubject: " + subject + "\r\n\r\nReset link test\r\n")
	if err := netsmtp.SendMail(ln.Addr().String(), nil, from, []string{to}, msg); err != nil {
		t.Fatal(err)
	}
	t.Logf("check inbox for subject: %s", subject)
}
