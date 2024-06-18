package User_email

import (
	"shopperia/src/common/models"
	"shopperia/src/email"
)

type Email interface {
	SendEmail(form models.SendEmailForm) error
}

type emailSender struct {
	Email
}

func NewEmailSender(accountEmail, accountName, password, host, serverName string) *emailSender {
	sender := email.NewEmailSender(accountEmail, accountName, password, host, serverName)
	return &emailSender{
		Email: sender,
	}
}

var prefix string = "./src/core/user/infrastructure/Email/templates"

func (e *emailSender) SendWelcomeEmail(emailContent models.WelcomeEmail, email models.EmailDto) error {

	templatePath := prefix + "/Welcome.html"

	form := models.SendEmailForm{
		Subject:          email.Subject,
		DestinationEmail: email.DestinationEmail,
		DestinationName:  email.DestinationName,
		TemplatePath:     templatePath,
		TemplateData:     emailContent,
	}

	err := e.Email.SendEmail(form)
	if err != nil {
		return err
	}

	return nil
}

func (e *emailSender) SendLoginConfirmationEmail(emailContent models.EmailDto) error {
	templatePath := prefix + "/LoginConfirm.html"

	form := models.SendEmailForm{
		Subject:          "Confirm login on your device",
		DestinationEmail: emailContent.DestinationEmail,
		DestinationName:  emailContent.DestinationName,
		TemplatePath:     templatePath,
		TemplateData:     emailContent,
	}

	if err := e.Email.SendEmail(form); err != nil {
		return err
	}

	return nil
}
