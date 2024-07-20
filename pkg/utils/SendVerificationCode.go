package utils

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
)

func GenerateRandomCode(length int) string {
	const charset = "0123456789"
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}
var SendVerificationEmailIMpl = SendVerificationEmail
func SendVerificationEmail(email, code string) error {
	from := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")
	to := []string{email}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	subject := "Subject: Email Verification\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
		<html>
			<body>
				<h2>Email Verification</h2>
				<p>Your verification code is: <strong>%s</strong></p>
				<p>This code will expire in 15 minutes.</p>
			</body>
		</html>
	`, code)

	message := []byte(subject + mime + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	return err
}

