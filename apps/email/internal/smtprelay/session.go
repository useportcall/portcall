package smtprelay

import (
	"io"
	"log"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/useportcall/portcall/libs/go/emailx"
)

type Session struct {
	backend *Backend
	from    string
	to      []string
	data    []byte
}

func (s *Session) AuthMechanisms() []string { return []string{sasl.Plain} }
func (s *Session) Reset()                   {}
func (s *Session) Logout() error            { return nil }

func (s *Session) Auth(mech string) (sasl.Server, error) {
	return sasl.NewPlainServer(func(_, username, _ string) error {
		log.Printf("auth mechanism=%s user=%s", mech, username)
		return nil
	}), nil
}

func (s *Session) Mail(from string, _ *smtp.MailOptions) error { s.from = from; return nil }
func (s *Session) Rcpt(to string, _ *smtp.RcptOptions) error   { s.to = append(s.to, to); return nil }

func (s *Session) Data(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	s.data = data
	switch s.backend.cfg.EmailProvider {
	case "resend":
		return s.sendViaResend()
	case "postmark":
		return s.sendViaPostmark()
	default:
		client, err := emailx.NewFromEnv()
		if err != nil {
			return err
		}
		subject, htmlBody, textBody := parseEmailData(string(s.data))
		content := htmlBody
		if content == "" {
			content = textBody
		}
		return client.Send(content, subject, s.from, s.to)
	}
}
