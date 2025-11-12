package notification

import (
	"fmt"
	"net/smtp"

	cfg "github.com/pobyzaarif/go-config"
)

type SMTPSender struct {
	Host     string `env:"SMTP_HOST"`
	Port     string `env:"SMTP_PORT"`
	Username string `env:"SMTP_USERNAME"`
	Password string `env:"SMTP_PASSWORD"`
}

func NewSMTPSenderFromEnv() (*SMTPSender, error) {
	s := &SMTPSender{}
	if err := cfg.LoadConfig(s); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *SMTPSender) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	msg := []byte(
		"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n\r\n" +
			body + "\r\n",
	)

	addr := fmt.Sprintf("%s:%v", s.Host, s.Port)
	return smtp.SendMail(addr, auth, s.Username, []string{to}, msg)
}
