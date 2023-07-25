package go_notifier_core

import (
	"testing"
)

// Mock SMTP server configuration for testing
var mockSmtpConfig = []byte(`
	{
		"Host": "sandbox.smtp.mailtrap.io",
		"Port": "2525",
		"Username": "username",
		"Password": "password",
		"Encryption": "tls"
	}`)

func TestSmtpMailerSend(t *testing.T) {
	// Create an instance of SmtpMailer
	mailer := &SmtpMailer{}

	// Set the configuration
	mailer.SetConfig(mockSmtpConfig)

	// Test email details
	fromName := "Test User"
	fromMail := "testuser@example.com"
	to := "recipient@example.com"
	subject := "Test Subject"
	message := "Hello, this is a test email."

	// Send an email
	err := mailer.Send(fromName, fromMail, to, subject, message)
	if err != nil {
		t.Fatalf("Error sending email: %v", err)
	}
}

func TestSmtpMailerSetConfig(t *testing.T) {
	// Create an instance of SmtpMailer
	mailer := &SmtpMailer{}

	// Set the configuration
	mailer.SetConfig(mockSmtpConfig)

	// Verify that the configuration was set correctly
	// You may add additional assertions here to check individual fields of the SmtpConfig if needed.
	if mailer.config == nil {
		t.Fatalf("SmtpMailer configuration is nil after SetConfig")
	}
}

func TestSmtpMailerSetConfigInvalidJSON(t *testing.T) {
	// Create an instance of SmtpMailer
	mailer := &SmtpMailer{}

	// Invalid JSON configuration
	invalidJSON := []byte("this is not valid JSON")

	// Set the invalid configuration
	mailer.SetConfig(invalidJSON)

	// Verify that the configuration remains nil
	if mailer.config != nil {
		t.Fatalf("SmtpMailer configuration is not nil after setting invalid configuration")
	}
}
