package mail

import (
	"fmt"
	"net/smtp"
	"strings"
)

type Mailer struct {
	host     string
	port     int
	username string
	password string
	auth     smtp.Auth
	from     string
}

func NewMailer(host string, port int, username, password string) *Mailer {
	auth := smtp.PlainAuth("", username, password, host)
	return &Mailer{
		host:     host,
		port:     port,
		username: username,
		password: password,
		auth:     auth,
		from:     username,
	}
}

// SendEmail sends a simple plain text email.
func (m *Mailer) SendEmail(to []string, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", m.host, m.port)

	msg := strings.Builder{}
	msg.WriteString(fmt.Sprintf("From: %s\r\n", m.from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(to, ",")))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)

	return smtp.SendMail(addr, m.auth, m.from, to, []byte(msg.String()))
}
