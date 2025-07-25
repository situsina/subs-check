package channel

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/bestruirui/bestsub/internal/modules/register"
)

type Email struct {
	Server   string `json:"server" type:"string" required:"true" description:"SMTP服务器地址"`
	Port     int    `json:"port" type:"int" required:"true" description:"SMTP服务器端口"`
	Username string `json:"username" type:"string" required:"true" description:"SMTP服务器用户名"`
	Password string `json:"password" type:"string" required:"true" description:"SMTP服务器密码"`
	From     string `json:"from" type:"string" required:"true" description:"发件人邮箱地址"`
	To       string `json:"to" type:"string" required:"true" description:"发送给"`
	TLS      bool   `json:"tls" type:"bool" required:"true" description:"是否启用TLS"`

	addr       string
	auth       smtp.Auth
	recipients []string
}

func (e *Email) Init() error {
	e.addr = fmt.Sprintf("%s:%d", e.Server, e.Port)
	e.auth = smtp.PlainAuth("", e.Username, e.Password, e.Server)

	recipients := strings.Split(e.To, ",")
	e.recipients = make([]string, len(recipients))
	for i, recipient := range recipients {
		e.recipients[i] = strings.TrimSpace(recipient)
	}

	return nil
}

func (e *Email) Send(title string, body *bytes.Buffer) error {
	if body == nil {
		return fmt.Errorf("email body is nil")
	}

	message := e.buildMessage(title, body)

	if err := e.sendMail(message); err != nil {
		return fmt.Errorf("send email failed: %w", err)
	}

	return nil
}

func (e *Email) buildMessage(subject string, body *bytes.Buffer) *bytes.Buffer {
	var message bytes.Buffer

	message.WriteString(fmt.Sprintf("From: %s\r\n", e.From))
	message.WriteString(fmt.Sprintf("To: %s\r\n", e.To))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	message.WriteString("\r\n")

	body.WriteTo(&message)

	return &message
}
func (e *Email) sendMail(message *bytes.Buffer) error {
	if e.TLS {
		return e.sendMailWithTLS(message)
	} else {
		return smtp.SendMail(e.addr, e.auth, e.From, e.recipients, message.Bytes())
	}
}

func (e *Email) sendMailWithTLS(message *bytes.Buffer) error {
	tlsConfig := &tls.Config{
		ServerName: e.Server,
	}
	conn, err := tls.Dial("tcp", e.addr, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, e.Server)
	if err != nil {
		return err
	}
	defer client.Quit()

	if err := client.Auth(e.auth); err != nil {
		return err
	}

	if err := client.Mail(e.From); err != nil {
		return err
	}

	for _, recipient := range e.recipients {
		if err := client.Rcpt(recipient); err != nil {
			return err
		}
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}
	defer writer.Close()

	if _, err := writer.Write(message.Bytes()); err != nil {
		return err
	}

	return nil
}

func init() {
	register.Notify(&Email{})
}
