package services

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strings"

	"go.uber.org/zap"
)

type EmailService struct {
	smtpHost     string
	smtpPort     string
	smtpUser     string
	smtpPassword string
	fromEmail    string
	fromName     string
	logger       *zap.Logger
	templates    map[string]*template.Template
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

type Email struct {
	To      []string
	Subject string
	Body    string
	HTML    string
}

func NewEmailService(config EmailConfig, logger *zap.Logger) *EmailService {
	es := &EmailService{
		smtpHost:     config.SMTPHost,
		smtpPort:     config.SMTPPort,
		smtpUser:     config.SMTPUser,
		smtpPassword: config.SMTPPassword,
		fromEmail:    config.FromEmail,
		fromName:     config.FromName,
		logger:       logger,
		templates:    make(map[string]*template.Template),
	}

	// Load email templates
	es.loadTemplates()

	return es
}

func (es *EmailService) loadTemplates() {
	templates := map[string]string{
		"password_reset": `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Password Reset</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h2 style="color: #2c3e50;">Password Reset Request</h2>
		<p>Hello,</p>
		<p>You requested to reset your password. Click the button below to reset it:</p>
		<div style="text-align: center; margin: 30px 0;">
			<a href="{{.ResetURL}}" style="background-color: #3498db; color: white; padding: 12px 24px; text-decoration: none; border-radius: 5px; display: inline-block;">Reset Password</a>
		</div>
		<p>Or copy and paste this link into your browser:</p>
		<p style="word-break: break-all; color: #3498db;">{{.ResetURL}}</p>
		<p>This link will expire in 1 hour.</p>
		<p>If you didn't request this, please ignore this email.</p>
		<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
		<p style="color: #7f8c8d; font-size: 12px;">This is an automated message, please do not reply.</p>
	</div>
</body>
</html>`,
		"welcome": `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Welcome</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h2 style="color: #2c3e50;">Welcome to Base App!</h2>
		<p>Hello {{.Name}},</p>
		<p>Thank you for signing up! Your account has been created successfully.</p>
		<p>You can now log in and start using our services.</p>
		<div style="text-align: center; margin: 30px 0;">
			<a href="{{.LoginURL}}" style="background-color: #27ae60; color: white; padding: 12px 24px; text-decoration: none; border-radius: 5px; display: inline-block;">Log In</a>
		</div>
		<p>If you have any questions, please don't hesitate to contact our support team.</p>
		<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
		<p style="color: #7f8c8d; font-size: 12px;">This is an automated message, please do not reply.</p>
	</div>
</body>
</html>`,
		"notification": `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Notification</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h2 style="color: #2c3e50;">{{.Title}}</h2>
		<p>{{.Message}}</p>
		{{if .Link}}
		<div style="text-align: center; margin: 30px 0;">
			<a href="{{.Link}}" style="background-color: #3498db; color: white; padding: 12px 24px; text-decoration: none; border-radius: 5px; display: inline-block;">View Details</a>
		</div>
		{{end}}
		<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
		<p style="color: #7f8c8d; font-size: 12px;">This is an automated message, please do not reply.</p>
	</div>
</body>
</html>`,
	}

	for name, content := range templates {
		tmpl, err := template.New(name).Parse(content)
		if err == nil {
			es.templates[name] = tmpl
		}
	}
}

func (es *EmailService) SendEmail(ctx context.Context, email Email) error {
	// If SMTP is not configured, log and return (for development)
	if es.smtpHost == "" || es.smtpPort == "" {
		es.logger.Info("Email not sent (SMTP not configured)",
			zap.String("to", strings.Join(email.To, ",")),
			zap.String("subject", email.Subject),
		)
		return nil // Don't fail if email is not configured
	}

	auth := smtp.PlainAuth("", es.smtpUser, es.smtpPassword, es.smtpHost)

	// Build email message
	msg := []byte(fmt.Sprintf("From: %s <%s>\r\n", es.fromName, es.fromEmail) +
		fmt.Sprintf("To: %s\r\n", strings.Join(email.To, ",")) +
		fmt.Sprintf("Subject: %s\r\n", email.Subject) +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		email.HTML)

	addr := fmt.Sprintf("%s:%s", es.smtpHost, es.smtpPort)
	err := smtp.SendMail(addr, auth, es.fromEmail, email.To, msg)
	if err != nil {
		es.logger.Error("Failed to send email", zap.Error(err))
		return err
	}

	es.logger.Info("Email sent successfully", zap.String("to", strings.Join(email.To, ",")))
	return nil
}

func (es *EmailService) SendPasswordResetEmail(ctx context.Context, to, resetToken, resetURL string) error {
	tmpl, ok := es.templates["password_reset"]
	if !ok {
		return fmt.Errorf("password reset template not found")
	}

	var buf bytes.Buffer
	data := map[string]string{
		"ResetURL": resetURL + "?token=" + resetToken,
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	email := Email{
		To:      []string{to},
		Subject: "Password Reset Request",
		HTML:    buf.String(),
	}

	return es.SendEmail(ctx, email)
}

func (es *EmailService) SendWelcomeEmail(ctx context.Context, to, name, loginURL string) error {
	tmpl, ok := es.templates["welcome"]
	if !ok {
		return fmt.Errorf("welcome template not found")
	}

	var buf bytes.Buffer
	data := map[string]string{
		"Name":     name,
		"LoginURL": loginURL,
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	email := Email{
		To:      []string{to},
		Subject: "Welcome to Base App!",
		HTML:    buf.String(),
	}

	return es.SendEmail(ctx, email)
}

func (es *EmailService) SendNotificationEmail(ctx context.Context, to, title, message, link string) error {
	tmpl, ok := es.templates["notification"]
	if !ok {
		// Fallback to simple HTML
		email := Email{
			To:      []string{to},
			Subject: title,
			HTML:    fmt.Sprintf("<html><body><h2>%s</h2><p>%s</p></body></html>", title, message),
		}
		return es.SendEmail(ctx, email)
	}

	var buf bytes.Buffer
	data := map[string]string{
		"Title":   title,
		"Message": message,
		"Link":    link,
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	email := Email{
		To:      []string{to},
		Subject: title,
		HTML:    buf.String(),
	}

	return es.SendEmail(ctx, email)
}

// GetEmailConfigFromEnv loads email configuration from environment variables
func GetEmailConfigFromEnv() EmailConfig {
	return EmailConfig{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPUser:     os.Getenv("SMTP_USER"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromEmail:    getEnv("SMTP_FROM_EMAIL", "noreply@baseapp.com"),
		FromName:     getEnv("SMTP_FROM_NAME", "Base App"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

