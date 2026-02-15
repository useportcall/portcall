package smtprelay

import (
	"log"
	"time"

	"github.com/emersion/go-smtp"
)

func RunFromEnv() error {
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		return err
	}
	return Run(cfg)
}

func Run(cfg Config) error {
	be := &Backend{cfg: cfg}
	s := smtp.NewServer(be)
	s.Addr = cfg.SMTPAddr
	s.Domain = cfg.SMTPAddr
	s.AllowInsecureAuth = true
	s.MaxMessageBytes = 10 * 1024 * 1024
	s.ReadTimeout = 30 * time.Second
	s.WriteTimeout = 30 * time.Second
	s.ErrorLog = newSMTPErrorLogger()
	log.Printf("SMTP relay starting on %s via %s", s.Addr, cfg.EmailProvider)
	return s.ListenAndServe()
}
