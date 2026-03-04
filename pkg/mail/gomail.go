package mail

import (
	"crypto/tls"
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

type (
	IMailService interface {
		SendEmail(from, subject, text, html string, to ...string) error
	}

	GomailService struct {
		*GomailServiceConfig
	}

	GomailServiceConfig struct {
		Host     string `envconfig:"MAIL_HOST" default:"host" validate:"required"`
		Port     int    `envconfig:"MAIL_PORT" default:"465" validate:"required"`
		Username string `envconfig:"MAIL_USERNAME" default:"username" validate:"required"`
		Password string `envconfig:"MAIL_PASSWORD" default:"password" validate:"required"`
		From     []string
	}
)

func NewGomailService(cfg *GomailServiceConfig) *GomailService {
	return &GomailService{
		cfg,
	}
}

func (g *GomailService) SendEmail(from, subject, text, html string, to ...string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)

	// Set body content
	if html != "" {
		m.SetBody("text/html", html)
	} else if text != "" {
		m.SetBody("text/plain", text)
	}

	d := gomail.NewDialer(g.Host, g.Port, g.Username, g.Password)

	d.TLSConfig = &tls.Config{
		ServerName:         g.Host,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
	}

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email via gomail: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent successfully via gomail to: %v", to)

	return nil
}
