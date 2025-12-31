package mailer

import (
	"fmt"
	"furniture-shop/internal/config"
	"log"
	"net/smtp"
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

func NewSender() Sender {
	return &smtpSender{
		host: config.Env.EmailSenderHost,
		port: config.Env.EmailSenderPort,
		user: config.Env.EmailSenderUser,
		pass: config.Env.EmailSenderPass,
		from: config.Env.EmailSenderFrom,
	}
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
