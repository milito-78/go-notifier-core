package go_notifier_core

import (
	"encoding/json"
	"net/smtp"
)

type (
	Mailer interface {
		Send(fromName, fromMail, to, subject, message string) error
		SetConfig(config []byte)
	}

	SmtpConfig struct {
		Host       string
		Port       string
		Username   string
		Password   string
		Encryption string
	}

	SmtpMailer struct {
		config *SmtpConfig
	}
)

func (s *SmtpMailer) Send(fromName, fromMail, to, subject, message string) error {
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	err := smtp.SendMail(
		s.config.Host+":"+s.config.Port,
		auth,
		fromMail,
		[]string{to},
		[]byte("From: "+fromName+" <"+fromMail+">\r\n"+
			"To: "+to+"\r\n"+
			"Subject: "+subject+"\r\n"+
			"\r\n"+
			message+"\r\n"),
	)
	return err
}

func (s *SmtpMailer) SetConfig(config []byte) {
	err := json.Unmarshal(config, &s.config)
	if err != nil {
		return
	}
}
