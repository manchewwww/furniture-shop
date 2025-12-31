package mailer

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

type Sender interface {
	Send(to, subject, body string) error
}

type smtpSender struct {
	host string
	port string
	user string
	pass string
	from string
}

func NewSenderFromEnv() Sender {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("FROM_EMAIL")
	if host == "" || port == "" || user == "" || pass == "" || from == "" {
		log.Printf("MAILER: SMTP not fully configured; emails will be logged only")
		return &logSender{}
	}
	return &smtpSender{host: host, port: port, user: user, pass: pass, from: from}
}

func (s *smtpSender) Send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := smtp.PlainAuth("", s.user, s.pass, s.host)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nFrom: %s\r\n\r\n%s\r\n", to, subject, s.from, body))
	return smtp.SendMail(addr, auth, s.from, []string{to}, msg)
}

type logSender struct{}

func (l *logSender) Send(to, subject, body string) error {
	log.Printf("MAIL LOG -> to=%s subject=%q body=%q", to, subject, body)
	return nil
}
