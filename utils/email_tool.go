package utils

import (
	"invar/logs"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	mail "github.com/xhit/go-simple-mail"
)

func SendEmail(to string, subject string, text string) {
	from := os.Getenv("EMAIL_ACCOUNT")
	pass := os.Getenv("EMAIL_PASSWORD")
	server := os.Getenv("EMAIL_SERVER")
	port, err := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	if err != nil {
		logs.Logger().WithFields(logrus.Fields{
			"name": "Smtp",
		}).Error("Expected nil, got", err, "for port convert to int.")
		return
	}

	client := mail.NewSMTPClient()
	client.Host = server
	client.Port = port
	client.Username = from
	client.Password = pass
	client.Encryption = mail.EncryptionSSL
	client.ConnectTimeout = 10 * time.Second
	client.SendTimeout = 10 * time.Second
	client.KeepAlive = false

	smtpClient, err := client.Connect()

	if err != nil {
		logs.Logger().WithFields(logrus.Fields{
			"name": "Smtp",
		}).Error("Expected nil, got", err, "generating email")
		return
	}

	err = smtpClient.Noop()

	if err != nil {
		logs.Logger().WithFields(logrus.Fields{
			"name": "Smtp",
		}).Error("Expected nil, got", err, "noop to client")
		return
	}

	email := mail.NewMSG()
	email.SetFrom("InVar <noreply@invar.finance>")
	email.AddTo(to)
	email.SetSubject(subject)
	email.SetBody(mail.TextPlain, text)

	if email.Error != nil {
		logs.Logger().WithFields(logrus.Fields{
			"name": "Smtp",
		}).Error("Expected nil, got", err, "generating email")
		return
	}

	err = email.Send(smtpClient)

	email.GetError()

	if err != nil {
		logs.Logger().WithFields(logrus.Fields{
			"name": "Smtp",
		}).Error("Expected nil, got", err, "sending email")
		return
	}
}
