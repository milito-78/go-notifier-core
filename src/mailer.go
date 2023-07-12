package src

import (
	"encoding/json"
	"net/smtp"
)

type (
	Mailer interface {
		Send(to, subject, message string) error
		SetConfig(config []byte)
	}

	SmtpConfig struct {
		Host     string
		Port     string
		From     string
		Password string
	}

	SmtpMailer struct {
		config *SmtpConfig
	}
)

func (s *SmtpMailer) Send(to, subject, message string) error {
	auth := smtp.PlainAuth("", s.config.From, s.config.Password, s.config.Host)
	err := smtp.SendMail(
		s.config.Host+":"+s.config.Port,
		auth,
		s.config.From,
		[]string{to},
		[]byte("Subject: "+subject+"\r\n\r\n"+
			message+"\r\n"),
	)
	return err
}

func (s *SmtpMailer) SetConfig(config []byte) {
	err := json.Unmarshal(config, s.config)
	if err != nil {
		return
	}
}
