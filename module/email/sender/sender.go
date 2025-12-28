package sender

import (
	"crypto/tls"
	"fmt"
	"hub-service/module/email/model"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

type EmailSender interface {
	Send(email *model.EmailMessage) error
}

// SMTPSender implements EmailSender using SMTP
type SMTPSender struct {
	host     string
	port     int
	username string
	password string
	from     string
	useTLS   bool
}

func NewSMTPSender() *SMTPSender {
	port, _ := strconv.Atoi(os.Getenv("SYSTEM_EMAIL_PORT"))
	if port == 0 {
		port = 587
	}

	useTLS := true
	if os.Getenv("SMTP_USE_TLS") == "false" {
		useTLS = false
	}

	return &SMTPSender{
		host:     os.Getenv("SYSTEM_EMAIL_HOST"),
		port:     port,
		username: os.Getenv("SYSTEM_EMAIL"),
		password: os.Getenv("SYSTEM_EMAIL_SERVER"),
		from:     os.Getenv("SYSTEM_EMAIL"),
		useTLS:   useTLS,
	}
}

func (s *SMTPSender) Send(email *model.EmailMessage) error {
	if s.host == "" {
		return fmt.Errorf("SMTP host not configured")
	}

	from := email.From
	if from == "" {
		from = s.from
	}

	// Build email message
	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(email.To, ", ")))

	if len(email.Cc) > 0 {
		msg.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(email.Cc, ", ")))
	}

	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))
	msg.WriteString("MIME-Version: 1.0\r\n")

	// Determine content type
	if email.HTMLBody != "" {
		msg.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
		msg.WriteString("\r\n")
		msg.WriteString(email.HTMLBody)
	} else {
		msg.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
		msg.WriteString("\r\n")
		msg.WriteString(email.Body)
	}

	// Collect all recipients
	recipients := append([]string{}, email.To...)
	recipients = append(recipients, email.Cc...)
	recipients = append(recipients, email.Bcc...)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	if s.useTLS && s.port == 465 {
		// SSL/TLS connection
		return s.sendWithTLS(addr, auth, from, recipients, msg.String())
	}

	// STARTTLS or plain connection
	return smtp.SendMail(addr, auth, from, recipients, []byte(msg.String()))
}

func (s *SMTPSender) sendWithTLS(addr string, auth smtp.Auth, from string, to []string, msg string) error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         s.host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS dial error: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("SMTP client error: %w", err)
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth error: %w", err)
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("SMTP mail error: %w", err)
	}

	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("SMTP rcpt error: %w", err)
		}
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP data error: %w", err)
	}

	_, err = writer.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("SMTP write error: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("SMTP close error: %w", err)
	}

	return client.Quit()
}
