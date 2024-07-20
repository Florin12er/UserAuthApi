package utils

import(
	"net/smtp"
	"os"
	"fmt"
)
func SendResetEmail(email, code string) error {
    from := os.Getenv("EMAIL_ADDRESS")
    password := os.Getenv("EMAIL_PASSWORD")
    to := []string{email}
    smtpHost := "smtp.gmail.com"
    smtpPort := "587"

    // HTML email template
    htmlTemplate := `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Password Reset</title>
        <style>
            body {
                font-family: Arial, sans-serif;
                line-height: 1.6;
                color: #333;
                max-width: 600px;
                margin: 0 auto;
                padding: 20px;
            }
            .container {
                background-color: #f9f9f9;
                border-radius: 5px;
                padding: 20px;
                text-align: center;
            }
            h1 {
                color: #2c3e50;
            }
            .code {
                font-size: 36px;
                font-weight: bold;
                color: #3498db;
                letter-spacing: 5px;
                margin: 20px 0;
                padding: 10px;
                background-color: #ecf0f1;
                border-radius: 5px;
            }
            .footer {
                margin-top: 20px;
                font-size: 12px;
                color: #7f8c8d;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <h1>Password Reset</h1>
            <p>You have requested to reset your password. Use the following code to complete the process:</p>
            <div class="code">%s</div>
            <p>This code will expire in 15 minutes.</p>
            <p>If you did not request a password reset, please ignore this email or contact support if you have concerns.</p>
            <div class="footer">
                <p>This is an automated message, please do not reply to this email.</p>
            </div>
        </div>
    </body>
    </html>
    `

    // Format the HTML template with the reset code
    htmlBody := fmt.Sprintf(htmlTemplate, code)

    // Compose the email
    mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
    subject := "Subject: Password Reset\n"
    msg := []byte(subject + mime + htmlBody)

    // Authenticate and send the email
    auth := smtp.PlainAuth("", from, password, smtpHost)
    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
    if err != nil {
        return err
    }

    return nil
}


