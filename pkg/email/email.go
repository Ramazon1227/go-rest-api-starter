package email

import (
	"fmt"
	"net/smtp"
)

var (
	Host     string
	Port     string
	Username string
	Password string
	From     string
)

func SendWelcomeEmail(to, name, userPassword string) error {
	auth := smtp.PlainAuth("", Username, Password, Host)

	subject := "Welcome to Our System"
	body := fmt.Sprintf(`
Dear %s,

An account has been created for you in the  ...  System.

Your login credentials:
Email: %s
Password: %s

Please change your password after your first login.

Best regards,
... System Team
`, name, to, userPassword)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", From, to, subject, body)

	addr := fmt.Sprintf("%s:%s", Host, Port)
	return smtp.SendMail(addr, auth, From, []string{to}, []byte(msg))
}
