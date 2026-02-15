package smtprelay

import (
	"log"
	"net/http"
	"time"

	"github.com/emersion/go-smtp"
)

type Backend struct {
	cfg        Config
	httpClient *http.Client
}

func (b *Backend) client() *http.Client {
	if b.httpClient != nil {
		return b.httpClient
	}
	return &http.Client{Timeout: 10 * time.Second}
}

func (b *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	log.Printf("new SMTP session from: %s", c.Conn().RemoteAddr())
	return &Session{backend: b}, nil
}
