package email

import (
	"log"
	"net/smtp"
	"shopperia/src/common/models"
)

var Client *smtp.Client

type EmailSender struct {
	AccountEmail string
	AccountName  string
	password     string
	host         string
	serverName   string
}

func NewEmailSender(accountEmail, accountName, password, host, serverName string) EmailSender {
	return EmailSender{
		AccountEmail: accountEmail,
		AccountName:  accountName,
		password:     password,
		host:         host,
		serverName:   serverName,
	}
}

func (sender *EmailSender) GetFrom() (string, string) {
	fromEmail := sender.AccountEmail
	fromName := sender.AccountName

	return fromEmail, fromName
}

func (sender *EmailSender) configClient(client *smtp.Client, from, to string) error {
	err := client.Mail(from)
	if err != nil {
		return err
	}

	err = client.Rcpt(to)
	if err != nil {
		return err
	}

	return nil
}

func (Sender EmailSender) SendEmail(form models.SendEmailForm) error {
	Email, err := NewEmail(form)
	if err != nil {
		return err
	}

	server := NewEmailServer(Sender.AccountEmail, Sender.password, Sender.host, Sender.serverName)

	client, err := newClient(server.conn, Sender.host, server.auth)
	if err != nil {
		log.Printf("error while creating client: %v", err)
	}

	err = Sender.configClient(client, Email.From.Email, Email.Dest.Email)
	if err != nil {
		log.Fatalf("Error al configurar el cliente: ERR %v", err)
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte(Email.Message))
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	if err := client.Quit(); err != nil {
		return err
	}

	return nil
}
